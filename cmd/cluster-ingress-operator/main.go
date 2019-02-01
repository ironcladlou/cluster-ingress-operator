package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openshift/cluster-ingress-operator/pkg/dns"
	awsdns "github.com/openshift/cluster-ingress-operator/pkg/dns/aws"
	"github.com/openshift/cluster-ingress-operator/pkg/operator"
	"github.com/openshift/cluster-ingress-operator/pkg/operator/support"

	configv1 "github.com/openshift/api/config/v1"

	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

const (
	// cloudCredentialsSecretName is the name of the secret in the
	// operator's namespace that will hold the credentials that the operator
	// will use to authenticate with the cloud API.
	cloudCredentialsSecretName = "cloud-credentials"
)

func main() {
	logrus.Infof("creating API client")
	// Get a kube client.
	kubeConfig, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("failed to get kube config: %v", err)
	}
	kubeClient, err := support.NewClient(kubeConfig)
	if err != nil {
		logrus.Fatalf("failed to create kube client: %v", err)
	}
	logrus.Infof("created API client")

	// Collect operator configuration.
	operatorNamespace, ok := os.LookupEnv("WATCH_NAMESPACE")
	if !ok {
		operatorNamespace = "openshift-ingress-operator"
	}

	// Set up the DNS manager.
	dnsManager, err := createDNSManager(kubeClient, operatorNamespace)
	if err != nil {
		logrus.Fatalf("failed to create DNS manager: %v", err)
	}
	logrus.Infof("created DNS manager")

	operatorConfig := operator.Config{
		Namespace:        operatorNamespace,
		OperandNamespace: "openshift-ingress",
	}
	operator, err := operator.New(operatorConfig, kubeClient, kubeConfig, dnsManager)
	if err != nil {
		logrus.Fatalf("failed to create operator: %v", err)
	}

	logrus.Infof("running operator until error or signal received")
	if err := operator.Start(signals.SetupSignalHandler()); err != nil {
		logrus.Fatalf("operator stopped with error: %v", err)
	} else {
		logrus.Infof("operator stopped by signal handler")
	}
}

// createDNSManager creates a DNS manager compatible with the given cluster
// configuration.
func createDNSManager(client client.Client, operatorNamespace string) (dns.Manager, error) {
	dnsConfig := &configv1.DNS{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, dnsConfig)
	if err != nil {
		logrus.Fatalf("failed to get dns 'cluster': %v", err)
	}

	infraConfig := &configv1.Infrastructure{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, infraConfig)
	if err != nil {
		logrus.Fatalf("failed to get infrastructure 'cluster': %v", err)
	}

	// Retrieve the typed cluster version config.
	clusterVersionConfig := &configv1.ClusterVersion{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: "version"}, clusterVersionConfig)
	if err != nil {
		logrus.Fatalf("failed to get clusterversion 'version': %v", err)
	}

	var dnsManager dns.Manager
	switch infraConfig.Status.Platform {
	case configv1.AWSPlatform:
		awsCreds := &corev1.Secret{}
		err := client.Get(context.TODO(), types.NamespacedName{Namespace: operatorNamespace, Name: cloudCredentialsSecretName}, awsCreds)
		if err != nil {
			return nil, fmt.Errorf("failed to get aws creds from %s/%s: %v", awsCreds.Namespace, awsCreds.Name, err)
		}
		manager, err := awsdns.NewManager(awsdns.Config{
			AccessID:   string(awsCreds.Data["aws_access_key_id"]),
			AccessKey:  string(awsCreds.Data["aws_secret_access_key"]),
			BaseDomain: strings.TrimSuffix(dnsConfig.Spec.BaseDomain, ".") + ".",
			ClusterID:  string(clusterVersionConfig.Spec.ClusterID),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create AWS DNS manager: %v", err)
		}
		logrus.Infof("using AWS DNS manager")
		dnsManager = manager
	default:
		logrus.Infof("using noop DNS manager")
		dnsManager = &dns.NoopManager{}
	}
	return dnsManager, nil
}
