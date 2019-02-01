package operator

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/openshift/cluster-ingress-operator/pkg/apis/ingress/v1alpha1"
	"github.com/openshift/cluster-ingress-operator/pkg/manifests"
	"github.com/openshift/cluster-ingress-operator/pkg/operator/support"
	"github.com/openshift/library-go/pkg/crypto"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	configv1 "github.com/openshift/api/config/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Controller struct {
	client client.Client

	operatorNamespace string
	operandNamespace  string
}

type graph struct {
	wildcardCAPair   *corev1.Secret
	wildcardCACert   *corev1.ConfigMap // owned by wildcardCAKey
	routerCR         *rbacv1.ClusterRole
	prometheusCR     *rbacv1.ClusterRole
	prometheusCRB    *rbacv1.ClusterRoleBinding // owned by prometheusCR
	operandNamespace *corev1.Namespace
	prometheusRole   *rbacv1.Role             // owned by operandNamespace
	prometheusRB     *rbacv1.RoleBinding      // owned by operandNamespace
	defaultIngress   *v1alpha1.ClusterIngress // owned by operandNamespace
}

func NewController(client client.Client, operatorNamespace, operandNamespace string) *Controller {
	return &Controller{
		client:            client,
		operatorNamespace: operatorNamespace,
		operandNamespace:  operandNamespace,
	}
}

func (c *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logrus.Infof("starting operator reconcile for request: %s", request)
	defer func() { logrus.Infof("finished operator reconcile for request: %s", request) }()

	desired, err := c.getDesiredGraph()
	if err != nil {
		return reconcile.Result{}, err
	}
	current, err := c.getCurrentGraph()
	if err != nil {
		return reconcile.Result{}, err
	}

	// compute changes in a pre-order traversal of the operator graph to ensure
	// dependencies are satisfied.
	changes := compare(current, desired)
	for _, change := range changes {
		logrus.Infof("applying change: %s", change)
		err := change.Apply(c.client)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to apply change: %v", err)
		}
	}
	return reconcile.Result{}, nil
}

// TODO: These could be split apart; right now the dependencies to
// generate each piece are implicit (e.g. the certs) or hidden in
// the manifest factory
func (c *Controller) getDesiredGraph() (*graph, error) {
	wildcardCAPair, err := c.getDesiredWildcardCAPair()
	if err != nil {
		return nil, err
	}

	wildcardCACert, err := c.getDesiredWildcardCACert(wildcardCAPair)
	if err != nil {
		return nil, err
	}

	routerClusterRole, err := c.getDesiredRouterClusterRole()
	if err != nil {
		return nil, err
	}

	operandNamespace, err := c.getDesiredOperandNamespace()
	if err != nil {
		return nil, err
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
	defaultIngress, err := c.getDesiredDefaultClusterIngress(ingressConfig, infraConfig)
	if err != nil {
		return nil, err
	}

	return &graph{
		wildcardCAPair:   wildcardCAPair,
		wildcardCACert:   wildcardCACert,
		routerCR:         routerClusterRole,
		operandNamespace: operandNamespace,
		defaultIngress:   defaultIngress,
	}, nil
}

func (c *Controller) getDesiredWildcardCAPair() (*corev1.Secret, error) {
	// TODO: Note that the signer and cert contents aren't deterministic, so comparing
	// will always produce a signer name diff on operator restart.
	signerName := fmt.Sprintf("%s@%d", "ingress-operator", time.Now().Unix())
	caConfig, err := crypto.MakeCAConfig(signerName, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to make CA config: %v", err)
	}
	certBytes, keyBytes, err := caConfig.GetPEMBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate: %v", err)
	}

	name := wildcardCAPairName(c.operatorNamespace)
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.crt": certBytes,
			"tls.key": keyBytes,
		},
	}, nil
}

func (c *Controller) getDesiredWildcardCACert(ca *corev1.Secret) (*corev1.ConfigMap, error) {
	name := wildcardCACertName(c.operatorNamespace)
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
		},
		Data: map[string]string{"ca-bundle.crt": string(ca.Data["tls.crt"])},
	}, nil
}

func (c *Controller) getDesiredRouterClusterRole() (*rbacv1.ClusterRole, error) {
	cr := manifests.RouterClusterRole()
	name := routerClusterRoleName()
	cr.Name = name.Name
	return cr, nil
}

func (c *Controller) getDesiredOperandNamespace() (*corev1.Namespace, error) {
	ns := manifests.OperandNamespace()
	name := operandNamespaceName(c.operandNamespace)
	ns.Name = name.Name
	return ns, nil
}

func (c *Controller) getDesiredDefaultClusterIngress(ingressConfig *configv1.Ingress, infraConfig *configv1.Infrastructure) (*v1alpha1.ClusterIngress, error) {
	ci := manifests.DefaultClusterIngress()

	name := defaultClusterIngressName(c.operatorNamespace)
	ci.Namespace = name.Namespace
	ci.Name = name.Name

	if len(ingressConfig.Spec.Domain) != 0 {
		ci.Spec.IngressDomain = &ingressConfig.Spec.Domain
	}

	switch infraConfig.Status.Platform {
	case configv1.AWSPlatform:
		ci.Spec.HighAvailability = &v1alpha1.ClusterIngressHighAvailability{
			Type: v1alpha1.CloudClusterIngressHA,
		}
	default:
		ci.Spec.HighAvailability = &v1alpha1.ClusterIngressHighAvailability{
			Type: v1alpha1.UserDefinedClusterIngressHA,
		}
	}

	return ci, nil
}

func (c *Controller) getCurrentGraph() (*graph, error) {
	wildcardCAPair, err := c.getCurrentWildcardCAPair()
	if err != nil {
		return nil, err
	}

	wildcardCACert, err := c.getCurrentWildcardCACert()
	if err != nil {
		return nil, err
	}

	routerCR, err := c.getCurrentRouterClusterRole()
	if err != nil {
		return nil, err
	}

	operandNamespace, err := c.getCurrentOperandNamespace()
	if err != nil {
		return nil, err
	}

	defaultIngress, err := c.getCurrentDefaultClusterIngress()
	if err != nil {
		return nil, err
	}

	return &graph{
		wildcardCAPair:   wildcardCAPair,
		wildcardCACert:   wildcardCACert,
		routerCR:         routerCR,
		operandNamespace: operandNamespace,
		defaultIngress:   defaultIngress,
	}, nil
}

func (c *Controller) getCurrentWildcardCAPair() (*corev1.Secret, error) {
	obj := &corev1.Secret{}
	err := c.client.Get(context.TODO(), wildcardCAPairName(c.operatorNamespace), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *Controller) getCurrentWildcardCACert() (*corev1.ConfigMap, error) {
	obj := &corev1.ConfigMap{}
	err := c.client.Get(context.TODO(), wildcardCACertName(c.operatorNamespace), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *Controller) getCurrentRouterClusterRole() (*rbacv1.ClusterRole, error) {
	obj := &rbacv1.ClusterRole{}
	err := c.client.Get(context.TODO(), routerClusterRoleName(), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *Controller) getCurrentOperandNamespace() (*corev1.Namespace, error) {
	obj := &corev1.Namespace{}
	err := c.client.Get(context.TODO(), operandNamespaceName(c.operandNamespace), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func (c *Controller) getCurrentDefaultClusterIngress() (*v1alpha1.ClusterIngress, error) {
	obj := &v1alpha1.ClusterIngress{}
	err := c.client.Get(context.TODO(), defaultClusterIngressName(c.operatorNamespace), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return obj, nil
}

func compare(current, desired *graph) (changes []support.Change) {
	changes = append(changes, getWildcardCAPairChanges(current, desired)...)
	changes = append(changes, getWildcardCACertChanges(current, desired)...)
	changes = append(changes, getRouterClusterRoleChanges(current, desired)...)
	changes = append(changes, getOperandNamespaceChanges(current, desired)...)
	changes = append(changes, getDefaultClusterIngressChanges(current, desired)...)
	return
}

func getWildcardCAPairChanges(current, desired *graph) (changes []support.Change) {
	if current.wildcardCAPair == nil {
		changes = append(changes, &support.CreateResource{Object: desired.wildcardCAPair})
	}
	return
}

func getWildcardCACertChanges(current, desired *graph) (changes []support.Change) {
	if current.wildcardCACert == nil {
		changes = append(changes, &support.CreateResource{Object: desired.wildcardCACert, Owner: desired.wildcardCAPair})
	}
	return
}

func getRouterClusterRoleChanges(current, desired *graph) (changes []support.Change) {
	if current.routerCR == nil {
		changes = append(changes, &support.CreateResource{Object: desired.routerCR})
	}
	return
}

func getOperandNamespaceChanges(current, desired *graph) (changes []support.Change) {
	if current.operandNamespace == nil {
		changes = append(changes, &support.CreateResource{Object: desired.operandNamespace})
	}
	return
}

func getDefaultClusterIngressChanges(current, desired *graph) (changes []support.Change) {
	if current.defaultIngress == nil {
		changes = append(changes, &support.CreateResource{Object: desired.defaultIngress})
	}
	return
}

func wildcardCAPairName(operatorNamespace string) types.NamespacedName {
	return types.NamespacedName{Namespace: operatorNamespace, Name: "router-ca"}
}

func wildcardCACertName(operatorNamespace string) types.NamespacedName {
	return types.NamespacedName{Namespace: operatorNamespace, Name: "router-ca"}
}

func routerClusterRoleName() types.NamespacedName {
	return types.NamespacedName{Name: "openshift-ingress-router"}
}

func operandNamespaceName(operandNamespace string) types.NamespacedName {
	return types.NamespacedName{Name: operandNamespace}
}

func defaultClusterIngressName(operatorNamespace string) types.NamespacedName {
	return types.NamespacedName{Namespace: operatorNamespace, Name: "default"}
}
