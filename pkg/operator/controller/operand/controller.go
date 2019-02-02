package operand

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/openshift/cluster-ingress-operator/pkg/apis/ingress/v1alpha1"
	"github.com/openshift/cluster-ingress-operator/pkg/dns"
	"github.com/openshift/cluster-ingress-operator/pkg/manifests"
	operatorconfig "github.com/openshift/cluster-ingress-operator/pkg/operator/config"
	support "github.com/openshift/cluster-ingress-operator/pkg/operator/support"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	configv1 "github.com/openshift/api/config/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/openshift/library-go/pkg/crypto"
)

type Controller struct {
	client            client.Client
	dns               dns.Manager
	configGetter      operatorConfigGetter
	operatorNamespace string
	operandNamespace  string
}

type graph struct {
	sa *corev1.ServiceAccount
	// owned by sa
	binding      *rbacv1.ClusterRoleBinding
	wildcardCert *corev1.Secret
	deployment   *appsv1.Deployment
	// owned by deployment
	clusterSvc      *corev1.Service
	loadBalancerSvc *corev1.Service
	// owned by loadBalancerSvc
	publicDNS  *dnsRecord
	privateDNS *dnsRecord
}

type dnsRecord struct{}

func NewController(client client.Client, dns dns.Manager, operatorNamespace, operandNamespace string) *Controller {
	return &Controller{
		client:            client,
		dns:               dns,
		operatorNamespace: operatorNamespace,
		operandNamespace:  operandNamespace,
		configGetter:      &envOperatorConfigGetter{},
	}
}

func (c *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logrus.Infof("starting operand reconcile of %s", request)
	defer func() { logrus.Infof("finished operand reconcile of %s", request) }()

	ci := &v1alpha1.ClusterIngress{}
	err := c.client.Get(context.TODO(), request.NamespacedName, ci)
	if err != nil {
		if errors.IsNotFound(err) {
			logrus.Infof("ignoring update to nonexistent key %s", request)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	if ci.Status.ObservedGeneration == ci.Generation {
		logrus.Infof("ignoring non-spec update to %s", request)
		return reconcile.Result{}, nil
	}
	if ci.Status.ObservedGeneration > ci.Generation {
		logrus.Infof("ignoring stale generation of %s", request)
		return reconcile.Result{}, nil
	}

	current, err := c.getCurrentGraph(ci)
	if err != nil {
		return reconcile.Result{}, err
	}
	desired, err := c.getDesiredGraph(ci)
	if err != nil {
		return reconcile.Result{}, err
	}

	changes := compare(current, desired, ci.DeletionTimestamp != nil)
	for _, change := range changes {
		logrus.Infof("applying change: %s", change)
		err := change.Apply(c.client)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to apply change to %s: %v", request, err)
		}
	}

	observed := ci.DeepCopyObject().(*v1alpha1.ClusterIngress)
	observed.Status.ObservedGeneration = observed.Generation
	err = c.client.Status().Update(context.TODO(), observed)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to record generation %d: %v", observed.Generation, err)
	} else {
		logrus.Infof("recorded new generation %d for %s", observed.Generation, request)
	}

	return reconcile.Result{}, nil
}

func (c *Controller) getDesiredGraph(ci *v1alpha1.ClusterIngress) (*graph, error) {
	sa, err := c.getDesiredServiceAccount(ci)
	if err != nil {
		return nil, err
	}

	binding, err := c.getDesiredClusterRoleBinding(ci)
	if err != nil {
		return nil, err
	}

	ca := &corev1.Secret{}
	err = c.client.Get(context.TODO(), wildcardCAPairName(c.operatorNamespace), ca)
	if err != nil {
		return nil, fmt.Errorf("failed to get wildcard CA pair: %v", err)
	}
	wildcardCert, err := c.getDesiredWildcardCert(ci, ca)
	if err != nil {
		return nil, err
	}

	operatorConfig, err := c.configGetter.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get operator confuig %v", err)
	}
	infraConfig := &configv1.Infrastructure{}
	err = c.client.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, infraConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get infrastructure config: %v", err)
	}
	ingressConfig := &configv1.Ingress{}
	err = c.client.Get(context.TODO(), types.NamespacedName{Name: "cluster"}, ingressConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingress config: %v", err)
	}
	deployment, err := c.getDesiredDeployment(ci, operatorConfig, infraConfig, ingressConfig)
	if err != nil {
		return nil, err
	}

	return &graph{
		sa:           sa,
		binding:      binding,
		wildcardCert: wildcardCert,
		deployment:   deployment,
	}, nil
}

func (c *Controller) getDesiredServiceAccount(ci *v1alpha1.ClusterIngress) (*corev1.ServiceAccount, error) {
	sa := manifests.RouterServiceAccount()
	name := serviceAccountName(c.operandNamespace, ci.Name)
	sa.Namespace = name.Namespace
	sa.Name = name.Name
	sa.Finalizers = append(sa.Finalizers, "ingress.openshift.io/ingress-controller")
	return sa, nil
}

func (c *Controller) getDesiredClusterRoleBinding(ci *v1alpha1.ClusterIngress) (*rbacv1.ClusterRoleBinding, error) {
	crb := manifests.RouterClusterRoleBinding()
	crbName := clusterRoleBindingName(ci.Name)
	saName := serviceAccountName(c.operandNamespace, ci.Name)
	crb.Name = crbName.Name
	crb.Subjects[0].Namespace = saName.Namespace
	crb.Subjects[0].Name = saName.Name
	return crb, nil
}

func (c *Controller) getDesiredWildcardCert(ci *v1alpha1.ClusterIngress, signingCA *corev1.Secret) (*corev1.Secret, error) {
	ca, err := crypto.GetCAFromBytes(signingCA.Data["tls.crt"], signingCA.Data["tls.key"])
	if err != nil {
		return nil, err
	}

	hostnames := sets.NewString(fmt.Sprintf("*.%s", *ci.Spec.IngressDomain))
	cert, err := ca.MakeServerCert(hostnames, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to make CA: %v", err)
	}

	certBytes, keyBytes, err := cert.GetPEMBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate: %v", err)
	}

	name := wildcardCertName(c.operandNamespace, ci.Name)
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.crt": certBytes,
			"tls.key": keyBytes,
		},
	}, nil
}

func (c *Controller) getCurrentGraph(ci *v1alpha1.ClusterIngress) (*graph, error) {
	sa, err := c.getCurrentServiceAccount(ci)
	if err != nil {
		return nil, err
	}

	binding, err := c.getCurrentClusterRoleBinding(ci)
	if err != nil {
		return nil, err
	}

	wildcardCert, err := c.getCurrentWildcardCert(ci)
	if err != nil {
		return nil, err
	}

	deployment, err := c.getCurrentDeployment(ci)
	if err != nil {
		return nil, err
	}

	return &graph{
		sa:           sa,
		binding:      binding,
		wildcardCert: wildcardCert,
		deployment:   deployment,
	}, nil
}

func (c *Controller) getCurrentServiceAccount(ci *v1alpha1.ClusterIngress) (*corev1.ServiceAccount, error) {
	obj := &corev1.ServiceAccount{}
	err := c.client.Get(context.TODO(), serviceAccountName(c.operandNamespace, ci.Name), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *Controller) getCurrentClusterRoleBinding(ci *v1alpha1.ClusterIngress) (*rbacv1.ClusterRoleBinding, error) {
	obj := &rbacv1.ClusterRoleBinding{}
	err := c.client.Get(context.TODO(), clusterRoleBindingName(ci.Name), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *Controller) getCurrentWildcardCert(ci *v1alpha1.ClusterIngress) (*corev1.Secret, error) {
	obj := &corev1.Secret{}
	err := c.client.Get(context.TODO(), wildcardCertName(c.operandNamespace, ci.Name), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *Controller) getCurrentDeployment(ci *v1alpha1.ClusterIngress) (*appsv1.Deployment, error) {
	obj := &appsv1.Deployment{}
	err := c.client.Get(context.TODO(), deploymentName(c.operandNamespace, ci.Name), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func compare(current, desired *graph, deleted bool) (changes []support.Change) {
	if deleted && current.sa != nil {
		changes = append(changes, &support.DeleteResource{Object: desired.sa})
	}
	changes = append(changes, getRouterServiceAccountChanges(current, desired, deleted)...)
	changes = append(changes, getRouterClusterRoleBindingChanges(current, desired, deleted)...)
	changes = append(changes, getRouterWildcardCertChanges(current, desired, deleted)...)
	changes = append(changes, getRouterDeploymentChanges(current, desired, deleted)...)
	// changes = append(changes, getRouterClusterServiceChanges(current, desired)...)
	// changes = append(changes, getRouterLBServiceChanges(current, desired)...)
	return
}

func getRouterServiceAccountChanges(current, desired *graph, deleted bool) (changes []support.Change) {
	if deleted {
		if current.sa != nil {
			changes = append(changes, &support.FinalizeResource{
				Object:    desired.sa,
				Finalizer: "ingress.openshift.io/ingress-controller",
				Actions: []interface{}{
					&support.DeleteResource{Object: desired.binding},
				},
			})
		}
		return
	}
	if current.sa == nil {
		changes = append(changes, &support.CreateResource{Object: desired.sa})
	}
	return
}

func getRouterClusterRoleBindingChanges(current, desired *graph, deleted bool) (changes []support.Change) {
	if deleted {
		// owner is responsible for cleaning up
		return
	}
	if current.binding == nil {
		changes = append(changes, &support.CreateResource{Object: desired.binding})
	}
	return
}

func getRouterWildcardCertChanges(current, desired *graph, deleted bool) (changes []support.Change) {
	if deleted {
		// owner is responsible for cleaning up
		return
	}
	if current.wildcardCert == nil {
		changes = append(changes, &support.CreateResource{Object: desired.wildcardCert, Owner: desired.sa})
	}
	return
}

func getRouterDeploymentChanges(current, desired *graph, deleted bool) (changes []support.Change) {
	if deleted {
		// owner is responsible for cleaning up
		return
	}
	if current.deployment == nil {
		changes = append(changes, &support.CreateResource{Object: desired.deployment, Owner: desired.sa})
		return
	}
	currentReplicas, desiredReplicas := current.deployment.Spec.Replicas, desired.deployment.Spec.Replicas
	if currentReplicas != nil && desiredReplicas != nil && *currentReplicas != *desiredReplicas {
		changes = append(changes, &support.ScaleResource{Object: desired.deployment, Replicas: *desiredReplicas})
	}
	return
}

func serviceAccountName(operandNamespace, clusterIngressName string) types.NamespacedName {
	return types.NamespacedName{Namespace: operandNamespace, Name: "router-" + clusterIngressName}
}

func clusterRoleBindingName(clusterIngressName string) types.NamespacedName {
	return types.NamespacedName{Name: "openshift-ingress-router:" + clusterIngressName}
}

func wildcardCAPairName(operatorNamespace string) types.NamespacedName {
	return types.NamespacedName{Namespace: operatorNamespace, Name: "router-ca"}
}

func wildcardCertName(operandNamespace, clusterIngressName string) types.NamespacedName {
	return types.NamespacedName{Namespace: operandNamespace, Name: "router-certs-" + clusterIngressName}
}

func deploymentName(operandNamespace, clusterIngressName string) types.NamespacedName {
	return types.NamespacedName{Namespace: operandNamespace, Name: "router-" + clusterIngressName}
}

type operatorConfigGetter interface {
	Get() (*operatorconfig.Config, error)
}

type envOperatorConfigGetter struct{}

func (g *envOperatorConfigGetter) Get() (*operatorconfig.Config, error) {
	routerImage := os.Getenv("IMAGE")
	if len(routerImage) == 0 {
		return nil, fmt.Errorf("couldn't find IMAGE environment variable")
	}
	return &operatorconfig.Config{
		RouterImage: routerImage,
	}, nil
}
