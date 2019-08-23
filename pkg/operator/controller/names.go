package controller

import (
	"fmt"

	operatorv1 "github.com/openshift/api/operator/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// GlobalMachineSpecifiedConfigNamespace is the location for global
	// config.  In particular, the operator will put the configmap with the
	// CA certificate in this namespace.
	GlobalMachineSpecifiedConfigNamespace = "openshift-config-managed"

	// ControllerDeploymentLabel identifies a deployment as an ingress controller
	// deployment, and the value is the name of the owning ingress controller.
	ControllerDeploymentLabel = "ingresscontroller.operator.openshift.io/deployment-ingresscontroller"
)

// IngressClusterOperatorName returns the namespaced name of the ClusterOperator
// resource for the operator.
func IngressClusterOperatorName() types.NamespacedName {
	return types.NamespacedName{
		Name: "ingress",
	}
}

// RouterDeploymentName returns the namespaced name for the router deployment.
func RouterDeploymentName(ci *operatorv1.IngressController) types.NamespacedName {
	return types.NamespacedName{
		Namespace: "openshift-ingress",
		Name:      "router-" + ci.Name,
	}
}

// RouterCASecretName returns the namespaced name for the router CA secret.
// This secret holds the CA certificate that the operator will use to create
// default certificates for ingresscontrollers.
func RouterCASecretName(operatorNamespace string) types.NamespacedName {
	return types.NamespacedName{
		Namespace: operatorNamespace,
		Name:      "router-ca",
	}
}

// RouterCAConfigMapName returns the namespaced name for the router CA configmap.
// The operator uses this configmap to publish the public key for the CA
// certificate, so that other operators can include it into their trust bundles.
func RouterCAConfigMapName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: GlobalMachineSpecifiedConfigNamespace,
		Name:      "router-ca",
	}
}

// RouterCertsGlobalSecretName returns the namespaced name for the router certs
// secret.  The operator uses this secret to publish the default certificates and
// their keys, so that the authentication operator can configure the OAuth server
// to use the same certificates.
func RouterCertsGlobalSecretName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: GlobalMachineSpecifiedConfigNamespace,
		Name:      "router-certs",
	}
}

// RouterOperatorGeneratedDefaultCertificateSecretName returns the namespaced name for
// the operator-generated router default certificate secret.
func RouterOperatorGeneratedDefaultCertificateSecretName(ci *operatorv1.IngressController, namespace string) types.NamespacedName {
	return types.NamespacedName{
		Namespace: namespace,
		Name:      fmt.Sprintf("router-certs-%s", ci.Name),
	}
}

// RouterPodDisruptionBudgetName returns the namespaced name for the router
// deployment's pod disruption budget.
func RouterPodDisruptionBudgetName(ic *operatorv1.IngressController) types.NamespacedName {
	return types.NamespacedName{
		Namespace: "openshift-ingress",
		Name:      "router-" + ic.Name,
	}
}

// RsyslogConfigMapName returns the namespaced name for the rsyslog configmap.
func RsyslogConfigMapName(ic *operatorv1.IngressController) types.NamespacedName {
	return types.NamespacedName{
		Namespace: "openshift-ingress",
		Name:      "rsyslog-conf-" + ic.Name,
	}
}

// RouterEffectiveDefaultCertificateSecretName returns the namespaced name for
// the in-use router default certificate secret.
func RouterEffectiveDefaultCertificateSecretName(ic *operatorv1.IngressController) types.NamespacedName {
	if cert := ic.Spec.DefaultCertificate; cert != nil {
		return types.NamespacedName{Namespace: ic.Namespace, Name: cert.Name}
	}
	return RouterOperatorGeneratedDefaultCertificateSecretName(ic, ic.Namespace)
}

func IngressControllerDeploymentLabel(ic *operatorv1.IngressController) string {
	return ic.Name
}

func IngressControllerDeploymentPodSelector(ic *operatorv1.IngressController) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: map[string]string{
			ControllerDeploymentLabel: IngressControllerDeploymentLabel(ic),
		},
	}
}

func InternalIngressControllerServiceName(ic *operatorv1.IngressController) types.NamespacedName {
	// TODO: remove hard-coded namespace
	return types.NamespacedName{Namespace: "openshift-ingress", Name: "router-internal-" + ic.Name}
}

func IngressControllerServiceMonitorName(ic *operatorv1.IngressController) types.NamespacedName {
	return types.NamespacedName{
		Namespace: "openshift-ingress",
		Name:      "router-" + ic.Name,
	}
}

func LoadBalancerServiceName(ic *operatorv1.IngressController) types.NamespacedName {
	return types.NamespacedName{Namespace: "openshift-ingress", Name: "router-" + ic.Name}
}

func WildcardDNSRecordName(ic *operatorv1.IngressController) types.NamespacedName {
	return types.NamespacedName{
		Namespace: ic.Namespace,
		Name:      fmt.Sprintf("%s-wildcard", ic.Name),
	}
}
