package iaaseki

import (
	"context"
	clientset "ekiOperator/generated/iaaseki/clientset/versioned"
	ekischeme "ekiOperator/generated/iaaseki/clientset/versioned/scheme"
	ekiinformers "ekiOperator/generated/iaaseki/informers/externalversions/iaaseki/v1"
	ekilisters "ekiOperator/generated/iaaseki/listers/iaaseki/v1"
	"fmt"
	"time"

	iaasekiv1 "ekiOperator/apis/iaaseki/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netinress "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	netinformers "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	netlisters "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const controllerAgentName = "EkiMonitor"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a EkiMonitor fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by EkiMonitor"
	// MessageResourceSynced is the message used for an Event fired when a EkiMonitor
	// is synced successfully
	MessageResourceSynced = "EkiMonitor synced successfully"
)

type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// ekiclientset is a clientset for our own API group
	ekiclientset clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	servicesLister    corelisters.ServiceLister
	serviceSynced     cache.InformerSynced
	ingressLister     netlisters.IngressLister
	ingressSynced     cache.InformerSynced

	ekiMonitorLister ekilisters.EkiMonitorLister
	ekiMonitorSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	ekiclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	serviceInformer coreinformers.ServiceInformer,
	IngressInformer netinformers.IngressInformer,
	ekiMonitorInformer ekiinformers.EkiMonitorInformer) *Controller {

	// Create event broadcaster
	// Add app-controller types to the default Kubernetes Scheme so Events can be
	// logged for app-controller types.
	utilruntime.Must(ekischeme.AddToScheme(ekischeme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	ekiController := &Controller{
		kubeclientset:     kubeclientset,
		ekiclientset:      ekiclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		servicesLister:    serviceInformer.Lister(),
		serviceSynced:     serviceInformer.Informer().HasSynced,
		ingressLister:     IngressInformer.Lister(),
		ingressSynced:     IngressInformer.Informer().HasSynced,
		ekiMonitorLister:  ekiMonitorInformer.Lister(),
		ekiMonitorSynced:  ekiMonitorInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Apps"),
		recorder:          recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when App resources change
	ekiMonitorInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: ekiController.enqueueApp,
		UpdateFunc: func(old, new interface{}) {
			ekiController.enqueueApp(new)
		},
	})

	return ekiController

}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting App controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")

	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.serviceSynced, c.ingressSynced, c.ekiMonitorSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	klog.Info("Starting workers")
	// Launch two workers to process App resources
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// App resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the App resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the App resource with this namespace/name
	ekiApp, err := c.ekiMonitorLister.EkiMonitors(namespace).Get(name)
	if err != nil {
		// The App resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("ekiApp '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	// Get the deployment with the name specified in ekiApp.spec
	deployment, err := c.deploymentsLister.Deployments(ekiApp.Namespace).Get(ekiApp.Spec.Deployment.Name)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(ekiApp.Namespace).Create(context.TODO(), newDeployment(ekiApp), metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	service, err := c.servicesLister.Services(ekiApp.Namespace).Get(ekiApp.Spec.Service.Name)
	if errors.IsNotFound(err) {
		service, err = c.kubeclientset.CoreV1().Services(ekiApp.Namespace).Create(context.TODO(), newService(ekiApp), metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	ingress, err := c.ingressLister.Ingresses(ekiApp.Namespace).Get(ekiApp.Spec.Ingress.Name)
	if errors.IsNotFound(err) {
		ingress, err = c.kubeclientset.NetworkingV1().Ingresses(ekiApp.Namespace).Create(context.TODO(), newIngress(ekiApp), metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	// If the Deployment is not controlled by this App resource, we should log
	// a warning to the event recorder and return error msg.
	if !metav1.IsControlledBy(deployment, ekiApp) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(ekiApp, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}
	if !metav1.IsControlledBy(service, ekiApp) {
		msg := fmt.Sprintf(MessageResourceExists, service.Name)
		c.recorder.Event(ekiApp, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}
	if !metav1.IsControlledBy(ingress, ekiApp) {
		msg := fmt.Sprintf(MessageResourceExists, ingress.Name)
		c.recorder.Event(ekiApp, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}
	c.recorder.Event(ekiApp, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// enqueueApp takes a App resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than App.
func (c *Controller) enqueueApp(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the App resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that App resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a App, we should not do anything more
		// with it.
		if ownerRef.Kind != "EkiMonitorApp" {
			return
		}

		App, err := c.ekiMonitorLister.EkiMonitors(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of App '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueApp(App)
		return
	}
}

// newDeployment creates a new Deployment for a App resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the App resource that 'owns' it.
func newDeployment(ekimonitor *iaasekiv1.EkiMonitor) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "ekimonitor-deployment",
		"controller": ekimonitor.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ekimonitor.Spec.Deployment.Name,
			Namespace: ekimonitor.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(ekimonitor, iaasekiv1.SchemeGroupVersion.WithKind("EkiMonitor")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &ekimonitor.Spec.Deployment.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  ekimonitor.Spec.Deployment.Name,
							Image: ekimonitor.Spec.Deployment.Image,
						},
					},
				},
			},
		},
	}
}
func newService(ekimonitor *iaasekiv1.EkiMonitor) *corev1.Service {
	labels := map[string]string{
		"app":        "ekimonitor-service",
		"controller": ekimonitor.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ekimonitor.Spec.Service.Name,
			Namespace: ekimonitor.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(ekimonitor, iaasekiv1.SchemeGroupVersion.WithKind("EkiMonitor")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.IntOrString{IntVal: 80},
				},
			},
		},
	}
}
func newIngress(ekimonitor *iaasekiv1.EkiMonitor) *netinress.Ingress {
	pathType := netinress.PathTypePrefix
	return &netinress.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ekimonitor.Spec.Ingress.Name,
			Namespace: ekimonitor.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(ekimonitor, iaasekiv1.SchemeGroupVersion.WithKind("EkiMonitor")),
			},
		},
		Spec: netinress.IngressSpec{
			Rules: []netinress.IngressRule{
				{
					IngressRuleValue: netinress.IngressRuleValue{
						HTTP: &netinress.HTTPIngressRuleValue{
							Paths: []netinress.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: netinress.IngressBackend{
										Service: &netinress.IngressServiceBackend{
											Name: ekimonitor.Spec.Ingress.Name,
											Port: netinress.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
