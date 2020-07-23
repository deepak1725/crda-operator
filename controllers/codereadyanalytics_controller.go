/*


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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	f8av1alpha1 "github.com/deepak1725/crda-operator/api/v1alpha1"
)

// CodeReadyAnalyticsReconciler reconciles a CodeReadyAnalytics object
type CodeReadyAnalyticsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=f8a.example.com,resources=codereadyanalytics,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=f8a.example.com,resources=codereadyanalytics/status,verbs=get;update;patch

func (r *CodeReadyAnalyticsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("codereadyanalytics", req.NamespacedName)

	// your logic here

	r.Log.Info("********************YESSSS********************************")

	// Fetch the CodeReadyAnalytics instance
	instance := &f8av1alpha1.CodeReadyAnalytics{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	var result *ctrl.Result

	// Secrets and Config Map Checks

	//POSTGRES Secret
	result, err = r.ensureSecret(req, instance, r.postgresSecret(instance))
	if result != nil {
		return *result, err
	}
	// AWS Secret
	result, err = r.ensureSecret(req, instance, r.awsSecret(instance))
	if result != nil {
		return *result, err
	}
	// 3Scale Secret
	result, err = r.ensureSecret(req, instance, r.threeScaleSecret(instance))
	if result != nil {
		return *result, err
	}
	// Worker Secret
	result, err = r.ensureSecret(req, instance, r.workerSecret(instance))
	if result != nil {
		return *result, err
	}
	// Gemini Secret
	result, err = r.ensureSecret(req, instance, r.geminiSecret(instance))
	if result != nil {
		return *result, err
	}
	// Config Maps
	result, err = r.ensureConfigMap(req, instance, r.bayesianConfigMap(instance))
	if result != nil {
		return *result, err
	}

	// == GREMLIN  ==========
	result, err = r.ensureDeployment(req, instance, r.gremlinDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.gremlinService(instance))
	if result != nil {
		return *result, err
	}

	// Persistent Vol
	r.Log.Info("Trying to create PV")
	result, err = r.ensurePV(req, instance, r.pvDeployment(instance))
	if err != nil {
		return *result, err
	}
	// Persistent Vol Claim
	result, err = r.ensurePVC(req, instance, r.pvcDeployment(instance))
	if err != nil {
		return *result, err
	}

	// PgBouncer
	result, err = r.ensureBouncerDeployment(req, instance, r.bouncerDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.bouncerService(instance))
	if result != nil {
		return *result, err
	}

	//Worker
	result, err = r.ensureDeployment(req, instance, r.workerDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == API Server Service  ==========
	result, err = r.ensureDeployment(req, instance, r.apiDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.apiService(instance))
	if result != nil {
		return *result, err
	}

	err = r.updateBackendStatus(instance)
	if err != nil {
		// Requeue the request if the status could not be updated
		return *result, err
	}

	return ctrl.Result{}, nil
}

func (r *CodeReadyAnalyticsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&f8av1alpha1.CodeReadyAnalytics{}).
		Complete(r)
}
