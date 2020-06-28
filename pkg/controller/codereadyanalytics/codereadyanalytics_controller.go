package codereadyanalytics

import (
	"context"

	openshiftv1alpha1 "operator/crda-operator/pkg/apis/openshift/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_codereadyanalytics")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CodeReadyAnalytics Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCodeReadyAnalytics{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("codereadyanalytics-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CodeReadyAnalytics
	err = c.Watch(&source.Kind{Type: &openshiftv1alpha1.CodeReadyAnalytics{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner CodeReadyAnalytics
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &openshiftv1alpha1.CodeReadyAnalytics{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCodeReadyAnalytics implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCodeReadyAnalytics{}

// ReconcileCodeReadyAnalytics reconciles a CodeReadyAnalytics object
type ReconcileCodeReadyAnalytics struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CodeReadyAnalytics object and makes changes based on the state read
// and what is in the CodeReadyAnalytics.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCodeReadyAnalytics) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CodeReadyAnalytics")

	// Fetch the CodeReadyAnalytics instance
	instance := &openshiftv1alpha1.CodeReadyAnalytics{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var result *reconcile.Result

	// Secrets and Config Map Checks
	result, err = r.ensureSecret(request, instance, r.postgresSecret(instance))
	if result != nil {
		return *result, err
	}
	// AWS
	result, err = r.ensureSecret(request, instance, r.awsSecret(instance))
	if result != nil {
		return *result, err
	}
	// 3Scale
	result, err = r.ensureSecret(request, instance, r.threeScaleSecret(instance))
	if result != nil {
		return *result, err
	}
	// Config Maps
	result, err = r.ensureConfigMap(request, instance, r.bayesianConfigMap(instance))
	if result != nil {
		return *result, err
	}

	// == API Server Service  ==========
	result, err = r.ensureDeployment(request, instance, r.apiDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(request, instance, r.apiService(instance))
	if result != nil {
		return *result, err
	}

	err = r.updateBackendStatus(instance)
	if err != nil {
		// Requeue the request if the status could not be updated
		return reconcile.Result{}, err
	}

	// Persistent Vol
	log.Info("Trying to create PV")
	result, err = r.ensurePV(request, instance, r.pvDeployment(instance))
	if err != nil {
		return *result, err
	}
	// Persistent Vol Claim
	result, err = r.ensurePVC(request, instance, r.pvcDeployment(instance))
	if err != nil {
		return *result, err
	}

	result, err = r.ensureBouncerDeployment(request, instance, r.bouncerDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(request, instance, r.bouncerService(instance))
	if result != nil {
		return *result, err
	}

	// result, err = r.handleApiServerChanges(instance)
	// if result != nil {
	// 	return *result, err
	// }

	return reconcile.Result{}, nil
}
