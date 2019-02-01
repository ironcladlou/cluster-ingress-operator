package manifests

import (
	"bytes"
	"io"

	ingressv1alpha1 "github.com/openshift/cluster-ingress-operator/pkg/apis/ingress/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	DefaultClusterIngressAsset    = "assets/defaults/cluster-ingress.yaml"
	OperandNamespaceAsset         = "assets/router/namespace.yaml"
	RouterServiceAccountAsset     = "assets/router/service-account.yaml"
	RouterClusterRoleAsset        = "assets/router/cluster-role.yaml"
	RouterClusterRoleBindingAsset = "assets/router/cluster-role-binding.yaml"
	RouterDeploymentAsset         = "assets/router/deployment.yaml"
	RouterClusterServiceAsset     = "assets/router/service-internal.yaml"
	RouterCloudServiceAsset       = "assets/router/service-cloud.yaml"

	// Annotation used to inform the certificate generation service to
	// generate a cluster-signed certificate and populate the secret.
	ServingCertSecretAnnotation = "service.alpha.openshift.io/serving-cert-secret-name"

	// Annotation used to enable the proxy protocol on the AWS load balancer.
	AWSLBProxyProtocolAnnotation = "service.beta.kubernetes.io/aws-load-balancer-proxy-protocol"
)

func MustAssetReader(asset string) io.Reader {
	return bytes.NewReader(MustAsset(asset))
}

func DefaultClusterIngress() *ingressv1alpha1.ClusterIngress {
	ci, err := NewClusterIngress(MustAssetReader(DefaultClusterIngressAsset))
	if err != nil {
		panic(err)
	}
	return ci
}

func OperandNamespace() *corev1.Namespace {
	ns, err := NewNamespace(MustAssetReader(OperandNamespaceAsset))
	if err != nil {
		panic(err)
	}
	return ns
}

func RouterServiceAccount() *corev1.ServiceAccount {
	sa, err := NewServiceAccount(MustAssetReader(RouterServiceAccountAsset))
	if err != nil {
		panic(err)
	}
	return sa
}

func RouterClusterRole() *rbacv1.ClusterRole {
	cr, err := NewClusterRole(MustAssetReader(RouterClusterRoleAsset))
	if err != nil {
		panic(err)
	}
	return cr
}

func RouterClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	crb, err := NewClusterRoleBinding(MustAssetReader(RouterClusterRoleBindingAsset))
	if err != nil {
		panic(crb)
	}
	return crb
}

func RouterDeployment() *appsv1.Deployment {
	deployment, err := NewDeployment(MustAssetReader(RouterDeploymentAsset))
	if err != nil {
		panic(err)
	}
	return deployment
}

func RouterClusterService(cr *ingressv1alpha1.ClusterIngress) *corev1.Service {
	s, err := NewService(MustAssetReader(RouterClusterServiceAsset))
	if err != nil {
		panic(s)
	}
	return s

	/*
		name := "router-internal-" + cr.Name

		s.Name = name

		if s.Labels == nil {
			s.Labels = map[string]string{}
		}
		s.Labels["router"] = name

		if s.Annotations == nil {
			s.Annotations = map[string]string{}
		}
		s.Annotations[ServingCertSecretAnnotation] = fmt.Sprintf("router-metrics-certs-%s", cr.Name)

		if s.Spec.Selector == nil {
			s.Spec.Selector = map[string]string{}
		}
		s.Spec.Selector["router"] = "router-" + cr.Name

		return s, nil
	*/
}

func RouterCloudService(cr *ingressv1alpha1.ClusterIngress) *corev1.Service {
	s, err := NewService(MustAssetReader(RouterCloudServiceAsset))
	if err != nil {
		panic(err)
	}
	return s
	/*

		name := "router-" + cr.Name

		s.Name = name

		if s.Labels == nil {
			s.Labels = map[string]string{}
		}
		s.Labels["router"] = name

		if s.Spec.Selector == nil {
			s.Spec.Selector = map[string]string{}
		}
		s.Spec.Selector["router"] = name

		if f.config.Platform == configv1.AWSPlatform {
			if s.Annotations == nil {
				s.Annotations = map[string]string{}
			}
			s.Annotations[AWSLBProxyProtocolAnnotation] = "*"
		}

		return s, nil
	*/
}

func NewServiceAccount(manifest io.Reader) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&sa); err != nil {
		return nil, err
	}

	return &sa, nil
}

func NewRole(manifest io.Reader) (*rbacv1.Role, error) {
	r := rbacv1.Role{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

func NewRoleBinding(manifest io.Reader) (*rbacv1.RoleBinding, error) {
	rb := rbacv1.RoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&rb); err != nil {
		return nil, err
	}

	return &rb, nil
}

func NewClusterRole(manifest io.Reader) (*rbacv1.ClusterRole, error) {
	cr := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&cr); err != nil {
		return nil, err
	}

	return &cr, nil
}

func NewClusterRoleBinding(manifest io.Reader) (*rbacv1.ClusterRoleBinding, error) {
	crb := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&crb); err != nil {
		return nil, err
	}

	return &crb, nil
}

func NewService(manifest io.Reader) (*corev1.Service, error) {
	s := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

func NewNamespace(manifest io.Reader) (*corev1.Namespace, error) {
	ns := corev1.Namespace{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&ns); err != nil {
		return nil, err
	}

	return &ns, nil
}

func NewDeployment(manifest io.Reader) (*appsv1.Deployment, error) {
	o := appsv1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&o); err != nil {
		return nil, err
	}

	return &o, nil
}

func NewClusterIngress(manifest io.Reader) (*ingressv1alpha1.ClusterIngress, error) {
	o := ingressv1alpha1.ClusterIngress{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&o); err != nil {
		return nil, err
	}

	return &o, nil
}
