package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/http/httpproxy"

	"github.com/ghodss/yaml"

	"github.com/openshift/cluster-ingress-operator/pkg/dns"
	awsdns "github.com/openshift/cluster-ingress-operator/pkg/dns/aws"
	azuredns "github.com/openshift/cluster-ingress-operator/pkg/dns/azure"
	logf "github.com/openshift/cluster-ingress-operator/pkg/log"
	"github.com/openshift/cluster-ingress-operator/pkg/manifests"
	"github.com/openshift/cluster-ingress-operator/pkg/operator"
	operatorclient "github.com/openshift/cluster-ingress-operator/pkg/operator/client"
	operatorconfig "github.com/openshift/cluster-ingress-operator/pkg/operator/config"
	statuscontroller "github.com/openshift/cluster-ingress-operator/pkg/operator/controller/status"

	configv1 "github.com/openshift/api/config/v1"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

const (
	// cloudCredentialsSecretName is the name of the secret in the
	// operator's namespace that will hold the credentials that the operator
	// will use to authenticate with the cloud API.
	cloudCredentialsSecretName = "cloud-credentials"
)

var log = logf.Logger.WithName("entrypoint")

func main() {
	metrics.DefaultBindAddress = ":60000"

	// Get a kube client.
	kubeConfig, err := config.GetConfig()
	if err != nil {
		log.Error(err, "failed to get kube config")
		os.Exit(1)
	}
	kubeClient, err := operatorclient.NewClient(kubeConfig)
	if err != nil {
		log.Error(err, "failed to create kube client")
		os.Exit(1)
	}

	// Collect operator configuration.
	operatorNamespace := os.Getenv("WATCH_NAMESPACE")
	if len(operatorNamespace) == 0 {
		operatorNamespace = manifests.DefaultOperatorNamespace
	}
	log.Info("using operator namespace", "namespace", operatorNamespace)

	ingressControllerImage := os.Getenv("IMAGE")
	if len(ingressControllerImage) == 0 {
		log.Error(fmt.Errorf("missing environment variable"), "'IMAGE' environment variable must be set")
		os.Exit(1)
	}
	releaseVersion := os.Getenv("RELEASE_VERSION")
	if len(releaseVersion) == 0 {
		releaseVersion = statuscontroller.UnknownVersionValue
		log.Info("RELEASE_VERSION environment variable missing", "release version", statuscontroller.UnknownVersionValue)
	}

	// Retrieve the cluster infrastructure config.
	infraConfig := &configv1.Infrastructure{}
	err = kubeClient.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, infraConfig)
	if err != nil {
		log.Error(err, "failed to get infrastructure 'config'")
		os.Exit(1)
	}

	dnsConfig := &configv1.DNS{}
	err = kubeClient.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, dnsConfig)
	if err != nil {
		log.Error(err, "failed to get dns 'cluster'")
		os.Exit(1)
	}

	platformStatus, err := getPlatformStatus(kubeClient, infraConfig)
	if err != nil {
		log.Error(err, "failed to get platform status")
		os.Exit(1)
	}

	// If the cluster proxy config specifies proxy info, create an HTTP client using the proxy
	// which can be used to configure the HTTP client for all non-apiserver connections.
	proxyConfig := &configv1.Proxy{}
	if err := kubeClient.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, proxyConfig); err != nil {
		if errors.IsNotFound(err) {
			log.Info("cluster proxy config was not found")
		} else {
			log.Error(err, "failed to get proxy 'cluster'")
			os.Exit(1)
		}
	}
	var proxyClient *http.Client
	if len(proxyConfig.Status.HTTPProxy) > 0 || len(proxyConfig.Status.HTTPSProxy) > 0 {
		log.Info("using HTTP client with cluster proxy config: %#v", proxyConfig)
		proxyClient = newProxyClient(proxyConfig)
	}

	operatorConfig := operatorconfig.Config{
		OperatorReleaseVersion: releaseVersion,
		Namespace:              operatorNamespace,
		IngressControllerImage: ingressControllerImage,
	}

	// Set up the DNS manager.
	dnsProvider, err := createDNSProvider(kubeClient, operatorConfig, dnsConfig, platformStatus, proxyClient)
	if err != nil {
		log.Error(err, "failed to create DNS manager")
		os.Exit(1)
	}

	// Set up and start the operator.
	op, err := operator.New(operatorConfig, dnsProvider, kubeConfig)
	if err != nil {
		log.Error(err, "failed to create operator")
		os.Exit(1)
	}
	if err := op.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "failed to start operator")
		os.Exit(1)
	}
}

// createDNSManager creates a DNS manager compatible with the given cluster
// configuration.
func createDNSProvider(cl client.Client, operatorConfig operatorconfig.Config, dnsConfig *configv1.DNS, platformStatus *configv1.PlatformStatus, httpClient *http.Client) (dns.Provider, error) {
	var dnsProvider dns.Provider
	switch platformStatus.Type {
	case configv1.AWSPlatformType:
		awsCreds := &corev1.Secret{}
		err := cl.Get(context.TODO(), types.NamespacedName{Namespace: operatorConfig.Namespace, Name: cloudCredentialsSecretName}, awsCreds)
		if err != nil {
			return nil, fmt.Errorf("failed to get aws creds from secret %s/%s: %v", awsCreds.Namespace, awsCreds.Name, err)
		}
		log.Info("using aws creds from secret", "namespace", awsCreds.Namespace, "name", awsCreds.Name)
		provider, err := awsdns.NewProvider(awsdns.Config{
			AccessID:   string(awsCreds.Data["aws_access_key_id"]),
			AccessKey:  string(awsCreds.Data["aws_secret_access_key"]),
			DNS:        dnsConfig,
			Region:     platformStatus.AWS.Region,
			HTTPClient: httpClient,
		}, operatorConfig.OperatorReleaseVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to create AWS DNS manager: %v", err)
		}
		dnsProvider = provider
	case configv1.AzurePlatformType:
		azureCreds := &corev1.Secret{}
		err := cl.Get(context.TODO(), types.NamespacedName{Namespace: operatorConfig.Namespace, Name: cloudCredentialsSecretName}, azureCreds)
		if err != nil {
			return nil, fmt.Errorf("failed to get azure creds from secret %s/%s: %v", azureCreds.Namespace, azureCreds.Name, err)
		}
		log.Info("using azure creds from secret", "namespace", azureCreds.Namespace, "name", azureCreds.Name)
		provider, err := azuredns.NewProvider(azuredns.Config{
			Environment:    "AzurePublicCloud",
			ClientID:       string(azureCreds.Data["azure_client_id"]),
			ClientSecret:   string(azureCreds.Data["azure_client_secret"]),
			TenantID:       string(azureCreds.Data["azure_tenant_id"]),
			SubscriptionID: string(azureCreds.Data["azure_subscription_id"]),
			DNS:            dnsConfig,
			HTTPClient:     httpClient,
		}, operatorConfig.OperatorReleaseVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure DNS manager: %v", err)
		}
		dnsProvider = provider
	default:
		dnsProvider = &dns.FakeProvider{}
	}
	return dnsProvider, nil
}

// getPlatformStatus provides a backwards-compatible way to look up platform status. AWS is the
// special case. 4.1 clusters on AWS expose the region config only through install-config. New AWS clusters
// and all other 4.2+ platforms are configured via platform status.
func getPlatformStatus(client client.Client, infra *configv1.Infrastructure) (*configv1.PlatformStatus, error) {
	status := infra.Status.PlatformStatus

	// Only AWS needs backwards compatibility with install-config
	if status.Type != configv1.AWSPlatformType {
		return status, nil
	}

	// Check whether the cluster config is already migrated
	if status.AWS != nil && len(status.AWS.Region) > 0 {
		return status, nil
	}

	// Otherwise build a platform status from the deprecated install-config
	type installConfig struct {
		Platform struct {
			AWS struct {
				Region string `json:"region"`
			} `json:"aws"`
		} `json:"platform"`
	}
	clusterConfigName := types.NamespacedName{Namespace: "kube-system", Name: "cluster-config-v1"}
	clusterConfig := &corev1.ConfigMap{}
	if err := client.Get(context.TODO(), clusterConfigName, clusterConfig); err != nil {
		return nil, fmt.Errorf("failed to get configmap %s: %v", clusterConfigName, err)
	}
	data, ok := clusterConfig.Data["install-config"]
	if !ok {
		return nil, fmt.Errorf("missing install-config in configmap")
	}
	var ic installConfig
	if err := yaml.Unmarshal([]byte(data), &ic); err != nil {
		return nil, fmt.Errorf("invalid install-config: %v\njson:\n%s", err, data)
	}
	return &configv1.PlatformStatus{
		Type: configv1.AWSPlatformType,
		AWS: &configv1.AWSPlatformStatus{
			Region: ic.Platform.AWS.Region,
		},
	}, nil
}

// newProxyClient returns an http.Client based on http.DefaultClient with Proxy
// set to a proxy function based on proxyConfig.
func newProxyClient(proxyConfig *configv1.Proxy) *http.Client {
	cfg := &httpproxy.Config{
		HTTPProxy:  proxyConfig.Status.HTTPProxy,
		HTTPSProxy: proxyConfig.Status.HTTPSProxy,
		NoProxy:    proxyConfig.Status.NoProxy,
	}
	proxyFn := cfg.ProxyFunc()
	proxy := func(req *http.Request) (*url.URL, error) {
		return proxyFn(req.URL)
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
