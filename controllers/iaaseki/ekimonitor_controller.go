/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package iaaseki

import (
	"context"
	iaasekiv1 "ekiOperator/apis/iaaseki/v1"
	ekiclientset "ekiOperator/generated/iaaseki/clientset/versioned"
	ekiinformer "ekiOperator/generated/iaaseki/informers/externalversions"
	ekisignal "ekiOperator/pkg/signals"
	"k8s.io/apimachinery/pkg/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// EkiMonitorReconciler reconciles a EkiMonitor object
type EkiMonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=iaaseki.cmss,resources=deployment,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=iaaseki.cmss,resources=service,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=iaaseki.cmss,resources=ingress,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=iaaseki.cmss,resources=ekimonitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=iaaseki.cmss,resources=ekimonitors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=iaaseki.cmss,resources=ekimonitors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EkiMonitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *EkiMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// TODO(user): your logic here
	klog.InitFlags(nil)
	// get config
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	kubeClient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubernetes client: %s", err.Error())
	}
	ekiClient, err := ekiclientset.NewForConfig(kubeconfig)
	if err != nil {
		klog.Fatalf("Error building eki client: %s", err.Error())
	}
	// generate two informers
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	ekiInformerFactory := ekiinformer.NewSharedInformerFactory(ekiClient, time.Second*30)

	ekiController := NewController(kubeClient, ekiClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		kubeInformerFactory.Core().V1().Services(),
		kubeInformerFactory.Networking().V1().Ingresses(),
		ekiInformerFactory.Iaaseki().V1().EkiMonitors())

	stopCh := ekisignal.SetupSignalHandler()
	kubeInformerFactory.Start(stopCh)
	ekiInformerFactory.Start(stopCh)

	if err = ekiController.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EkiMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&iaasekiv1.EkiMonitor{}).
		Complete(r)
}
