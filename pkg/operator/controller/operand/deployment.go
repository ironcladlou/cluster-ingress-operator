package operand

import (
	"fmt"

	"github.com/openshift/cluster-ingress-operator/pkg/apis/ingress/v1alpha1"
	"github.com/openshift/cluster-ingress-operator/pkg/manifests"
	operatorconfig "github.com/openshift/cluster-ingress-operator/pkg/operator/config"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	configv1 "github.com/openshift/api/config/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Controller) getDesiredDeployment(
	ci *v1alpha1.ClusterIngress,
	operatorConfig *operatorconfig.Config,
	infraConfig *configv1.Infrastructure,
	ingressConfig *configv1.Ingress,
) (*appsv1.Deployment, error) {
	deployment := manifests.RouterDeployment()
	name := deploymentName(c.operandNamespace, ci.Name)
	saName := serviceAccountName(c.operandNamespace, ci.Name)

	deployment.Namespace = name.Namespace
	deployment.Name = name.Name
	deployment.Spec.Template.Labels["router"] = name.Name
	deployment.Spec.Selector.MatchLabels["router"] = name.Name
	deployment.Spec.Template.Spec.ServiceAccountName = saName.Name

	env := []corev1.EnvVar{
		{Name: "ROUTER_SERVICE_NAME", Value: ci.Name},
	}
	if ci.Spec.IngressDomain != nil {
		env = append(env, corev1.EnvVar{Name: "ROUTER_CANONICAL_HOSTNAME", Value: *ci.Spec.IngressDomain})
	}
	if ci.Spec.HighAvailability != nil && ci.Spec.HighAvailability.Type == v1alpha1.CloudClusterIngressHA {
		// For now, check if we are on AWS. This can really be done for
		// for any external [cloud] LBs that support the proxy protocol.
		if infraConfig.Status.Platform == configv1.AWSPlatform {
			env = append(env, corev1.EnvVar{Name: "ROUTER_USE_PROXY_PROTOCOL", Value: "true"})
		}
	}

	if ci.Spec.NodePlacement != nil {
		if ci.Spec.NodePlacement.NodeSelector != nil {
			nodeSelector, err := metav1.LabelSelectorAsMap(ci.Spec.NodePlacement.NodeSelector)
			if err != nil {
				// TODO: this should be in validation and shouldn't be requeued
				return nil, fmt.Errorf("clusteringress has invalid spec.nodePlacement.nodeSelector: %v", err)
			}
			deployment.Spec.Template.Spec.NodeSelector = nodeSelector
		}
	}

	if ci.Spec.NamespaceSelector != nil {
		namespaceSelector, err := metav1.LabelSelectorAsSelector(ci.Spec.NamespaceSelector)
		if err != nil {
			// TODO: this should be in validation and shouldn't be requeued
			return nil, fmt.Errorf("clusteringress has invalid spec.namespaceSelector: %v", err)
		}
		env = append(env, corev1.EnvVar{
			Name:  "NAMESPACE_LABELS",
			Value: namespaceSelector.String(),
		})
	}

	replicas := ci.Spec.Replicas
	deployment.Spec.Replicas = &replicas

	if ci.Spec.RouteSelector != nil {
		routeSelector, err := metav1.LabelSelectorAsSelector(ci.Spec.RouteSelector)
		if err != nil {
			return nil, fmt.Errorf("clusteringress has invalid spec.routeSelector: %v", err)
		}
		env = append(env, corev1.EnvVar{Name: "ROUTE_LABELS", Value: routeSelector.String()})
	}

	deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, env...)

	deployment.Spec.Template.Spec.Containers[0].Image = operatorConfig.RouterImage

	if ci.Spec.HighAvailability != nil && ci.Spec.HighAvailability.Type == v1alpha1.UserDefinedClusterIngressHA {
		// Expose ports 80 and 443 on the host to provide endpoints for
		// the user's HA solution.
		deployment.Spec.Template.Spec.HostNetwork = true

		// With container networking, probes default to using the pod IP
		// address.  With host networking, probes default to using the
		// node IP address.  Using localhost avoids potential routing
		// problems or firewall restrictions.
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe.Handler.HTTPGet.Host = "localhost"
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler.HTTPGet.Host = "localhost"
	}

	// Fill in the default certificate secret name.
	secretName := wildcardCertName(c.operandNamespace, ci.Name)
	if ci.Spec.DefaultCertificateSecret != nil && len(*ci.Spec.DefaultCertificateSecret) > 0 {
		secretName = wildcardCertName(c.operandNamespace, *ci.Spec.DefaultCertificateSecret)
	}
	deployment.Spec.Template.Spec.Volumes[0].Secret.SecretName = secretName.Name

	return deployment, nil
}
