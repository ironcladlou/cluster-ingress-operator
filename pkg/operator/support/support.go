package support

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/davecgh/go-spew/spew"

	"github.com/openshift/cluster-ingress-operator/pkg/apis"
	"github.com/openshift/cluster-ingress-operator/pkg/util/slice"

	appsv1 "k8s.io/api/apps/v1"

	configv1 "github.com/openshift/api/config/v1"

	kscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Scheme contains all the API types necessary for the operator's dynamic
// clients to work. Any new non-core types must be added here.
var Scheme *runtime.Scheme

func init() {
	Scheme = kscheme.Scheme
	if err := apis.AddToScheme(Scheme); err != nil {
		panic(err)
	}
	if err := configv1.Install(Scheme); err != nil {
		panic(err)
	}
}

// NewClient builds an operator-compatible kube client from the given REST config.
func NewClient(kubeConfig *rest.Config) (client.Client, error) {
	managerOptions := manager.Options{
		Scheme:         Scheme,
		MapperProvider: apiutil.NewDiscoveryRESTMapper,
	}
	mapper, err := managerOptions.MapperProvider(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get API Group-Resources")
	}
	kubeClient, err := client.New(kubeConfig, client.Options{
		Scheme: Scheme,
		Mapper: mapper,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create kube client: %v", err)
	}
	return kubeClient, nil
}

type Change interface {
	Apply(cl client.Client) error
}

type CreateResource struct {
	Object runtime.Object
	Owner  runtime.Object
}

func (c *CreateResource) Apply(cl client.Client) error {
	obj := c.Object.DeepCopyObject()

	if c.Owner != nil && !reflect.ValueOf(c.Owner).IsNil() {
		owner := c.Owner.DeepCopyObject()
		ownerKey, err := client.ObjectKeyFromObject(owner)
		if err != nil {
			return err
		}
		err = cl.Get(context.TODO(), ownerKey, owner)
		if err != nil {
			return err
		}
		ownerAccessor, err := meta.Accessor(owner)
		if err != nil {
			return err
		}
		objAccessor, err := meta.Accessor(obj)
		if err != nil {
			return err
		}
		// TODO: inject scheme?
		ownerGVK, err := apiutil.GVKForObject(owner, Scheme)
		if err != nil {
			return err
		}
		ownerRef := metav1.NewControllerRef(ownerAccessor, ownerGVK)
		objAccessor.SetOwnerReferences([]metav1.OwnerReference{*ownerRef})
	}

	logrus.Infof("CreateResource: %v", spew.Sdump(obj))
	err := cl.Create(context.TODO(), obj)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

func (c *CreateResource) String() string {
	objKey, err := client.ObjectKeyFromObject(c.Object)
	if err != nil {
		objKey = types.NamespacedName{Namespace: "<unknown>", Name: "<unknown>"}
	}
	ownerKey, err := client.ObjectKeyFromObject(c.Owner)
	if err != nil {
		ownerKey = types.NamespacedName{Namespace: "<unknown>", Name: "<unknown>"}
	}
	return fmt.Sprintf("CreateResource[object: %s, owner: %s]", objKey, ownerKey)
}

type DeleteResource struct {
	Object runtime.Object
}

func (c *DeleteResource) Apply(cl client.Client) error {
	obj := c.Object.DeepCopyObject()
	logrus.Infof("DeleteResource: %v", spew.Sdump(obj))
	return cl.Delete(context.TODO(), obj)
}

func (c *DeleteResource) String() string {
	objKey, err := client.ObjectKeyFromObject(c.Object)
	if err != nil {
		objKey = types.NamespacedName{Namespace: "<unknown>", Name: "<unknown>"}
	}
	return fmt.Sprintf("DeleteResource[object: %s]", objKey)
}

type FinalizeResource struct {
	Object    runtime.Object
	Finalizer string
	Actions   []interface{}
}

func (c *FinalizeResource) Apply(cl client.Client) error {
	errors := []error{}
	for _, action := range c.Actions {
		switch a := action.(type) {
		case Change:
			err := a.Apply(cl)
			if err != nil {
				errors = append(errors, err)
			}
		default:
			errors = append(errors, fmt.Errorf("unsupported action: %v", action))
		}
	}
	if len(errors) > 0 {
		return utilerrors.NewAggregate(errors)
	}
	obj := c.Object.DeepCopyObject()
	objAccessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	if slice.ContainsString(objAccessor.GetFinalizers(), c.Finalizer) {
		objAccessor.SetFinalizers(slice.RemoveString(objAccessor.GetFinalizers(), c.Finalizer))
		logrus.Infof("finalizeResource: removing finalizer %q from %v", c.Finalizer, spew.Sdump(obj))
		err = cl.Update(context.TODO(), obj)
		if err != nil {
			return err
		}
	}
	logrus.Infof("FinalizeResource: finalized %v", spew.Sdump(obj))
	return nil
}

func (c *FinalizeResource) String() string {
	objKey, err := client.ObjectKeyFromObject(c.Object)
	if err != nil {
		objKey = types.NamespacedName{Namespace: "<unknown>", Name: "<unknown>"}
	}
	return fmt.Sprintf("FinalizeResource[object: %s, finalizer: %s: actions: %d]", objKey, c.Finalizer, len(c.Actions))
}

type ScaleResource struct {
	// TODO: Make this generic
	Object   *appsv1.Deployment
	Replicas int32
}

func (c *ScaleResource) Apply(cl client.Client) error {
	obj := c.Object.DeepCopyObject().(*appsv1.Deployment)
	obj.Spec.Replicas = &c.Replicas
	logrus.Infof("ScaleResource: %v", spew.Sdump(obj))
	return cl.Update(context.TODO(), obj)
}

func (c *ScaleResource) String() string {
	objKey, err := client.ObjectKeyFromObject(c.Object)
	if err != nil {
		objKey = types.NamespacedName{Namespace: "<unknown>", Name: "<unknown>"}
	}
	return fmt.Sprintf("ScaleResource[object: %s]", objKey)
}

func ObjectKeyFromObject(obj runtime.Object) client.ObjectKey {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		panic(err)
	}
	return client.ObjectKey{Namespace: accessor.GetNamespace(), Name: accessor.GetName()}
}
