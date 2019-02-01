package operator

import (
	"fmt"
	"time"

	"github.com/openshift/cluster-ingress-operator/pkg/apis/ingress/v1alpha1"
	"github.com/openshift/cluster-ingress-operator/pkg/dns"
	operandcontroller "github.com/openshift/cluster-ingress-operator/pkg/operator/controller/operand"
	operatorcontroller "github.com/openshift/cluster-ingress-operator/pkg/operator/controller/operator"
	"github.com/openshift/cluster-ingress-operator/pkg/operator/support"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Config is configuration for the operator and should include things like
// operated images, scheduling configuration, etc.
type Config struct {
	// Namespace is the operator namespace.
	Namespace string
	// OperandNamespace is the namespace of the component managed by the operator.
	OperandNamespace string
}

// Operator is the scaffolding for the ingress operator. It sets up dependencies
// and defines the toplogy of the operator and its managed components, wiring
// them together. Operator knows what namespace the operator lives in, and what
// specific resoure types in other namespaces should produce operator events.
type Operator struct {
	operatorManager        manager.Manager
	syncOperatorController chan event.GenericEvent

	config Config
}

// New creates (but does not start) a new operator from configuration.
func New(config Config, client client.Client, kubeConfig *rest.Config, dnsManager dns.Manager) (*Operator, error) {
	operatorManager, err := manager.New(kubeConfig, manager.Options{
		Namespace: config.Namespace,
		Scheme:    support.Scheme,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create operator manager: %v", err)
	}

	operatorController, err := controller.New("operator-controller", operatorManager, controller.Options{
		Reconciler: operatorcontroller.NewController(client, config.Namespace, config.OperandNamespace),
	})
	if err != nil {
		return nil, err
	}
	syncOperatorController := make(chan event.GenericEvent)
	operatorController.Watch(
		&source.Channel{Source: syncOperatorController},
		&handler.EnqueueRequestForObject{},
	)

	operandController, err := controller.New("operand-controller", operatorManager, controller.Options{
		Reconciler: operandcontroller.NewController(client, dnsManager, config.Namespace, config.OperandNamespace),
	})
	if err != nil {
		return nil, err
	}
	err = operandController.Watch(&source.Kind{Type: &v1alpha1.ClusterIngress{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return nil, err
	}

	// TODO: Set up watches for other resources in the operand namespace to queue
	// events in the operator and operand controllers.

	return &Operator{
		operatorManager:        operatorManager,
		syncOperatorController: syncOperatorController,
		config:                 config,
	}, nil
}

// Start creates the default ClusterIngress and then starts the operator
// synchronously until a message is received on the stop channel.
// TODO: Move the default ClusterIngress logic elsewhere.
func (o *Operator) Start(stop <-chan struct{}) error {
	logrus.Infof("starting operator")
	defer func() { logrus.Info("stopping operator") }()

	errChan := make(chan error)
	go func() {
		logrus.Infof("starting operator manager")
		if err := o.operatorManager.Start(stop); err != nil {
			logrus.Errorf("operator manager returned with an error: %v", err)
			errChan <- err
		} else {
			logrus.Infof("operator manager returned without error")
		}
	}()

	// TODO: Is there some other object we could watch to avoid this?
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			logrus.Infof("forcing operator controller resync")
			o.syncOperatorController <- getGenericOperatorEvent(o.config.Namespace, "ingress-operator")
		}
	}()
	o.syncOperatorController <- getGenericOperatorEvent(o.config.Namespace, "ingress-operator")

	var err error
	select {
	case <-stop:
	case managerErr := <-errChan:
		err = managerErr
	}
	ticker.Stop()
	return err
}

func getGenericOperatorEvent(operatorNamespace, name string) event.GenericEvent {
	return event.GenericEvent{
		Meta: &metav1.ObjectMeta{
			Namespace: operatorNamespace,
			Name:      name,
		},
	}
}
