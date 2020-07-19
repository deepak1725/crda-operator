package controllers

import (
	"context"
	f8av1alpha1 "github.com/deepak1725/crda-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *CodeReadyAnalyticsReconciler) ensureDeployment(request reconcile.Request,
	instance *f8av1alpha1.CodeReadyAnalytics,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		r.Log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			r.Log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		r.Log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}
	// Deployment was successful
	return nil, nil
}

func (r *CodeReadyAnalyticsReconciler) ensureService(request reconcile.Request,
	instance *f8av1alpha1.CodeReadyAnalytics,
	s *corev1.Service,
) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		r.Log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.Client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		r.Log.Error(err, "Failed to get Service")
		return &reconcile.Result{}, err
	}
	// Creation was successful
	return nil, nil
}

func (r *CodeReadyAnalyticsReconciler) ensureSecret(request reconcile.Request,
	instance *f8av1alpha1.CodeReadyAnalytics,
	secret *corev1.Secret,
) (*reconcile.Result, error) {
	found := &corev1.Secret{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      secret.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the secret
		r.Log.Info("Creating a new secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.Client.Create(context.TODO(), secret)

		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the secret not existing
		r.Log.Error(err, "Failed to get Secret")
		return &reconcile.Result{}, err
	}
	// Creation was successful
	return nil, nil
}

func (r *CodeReadyAnalyticsReconciler) ensureConfigMap(request reconcile.Request,
	instance *f8av1alpha1.CodeReadyAnalytics,
	configMap *corev1.ConfigMap,
) (*reconcile.Result, error) {

	// Check if this ConfigMap already exists
	foundMap := &corev1.ConfigMap{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, foundMap)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), configMap)
		if err != nil {
			// Creation failed
			r.Log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "configMap.Name", configMap.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the secret not existing
		r.Log.Error(err, "Failed to get ConfigMap")
		return &reconcile.Result{}, err
	}
	// Creation was successful
	return nil, nil
}

func labels(v *f8av1alpha1.CodeReadyAnalytics, tier string) map[string]string {
	return map[string]string{
		"app":          "analytics",
		"analytics_cr": v.Name,
		"tier":         tier,
	}
}

func (r *CodeReadyAnalyticsReconciler) postgresSecret(cr *f8av1alpha1.CodeReadyAnalytics) *corev1.Secret {
	labels := map[string]string{
		"app": cr.Name,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "coreapi-postgres",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"database":         cr.Spec.Config.Database.DbName,
			"host":             cr.Spec.Config.Database.Host,
			"initial-database": cr.Spec.Config.Database.InitialDatabase,
			"password":         cr.Spec.Config.Database.Password,
			"port":             cr.Spec.Config.Database.Port,
			"username":         cr.Spec.Config.Database.Username,
		},
	}
	return secret
}

func (r *CodeReadyAnalyticsReconciler) workerSecret(cr *f8av1alpha1.CodeReadyAnalytics) *corev1.Secret {
	labels := map[string]string{
		"app": cr.Name,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "worker",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"github-token":       cr.Spec.Config.Common.GithubToken,
			"libraries-io-token": cr.Spec.Config.Common.LibrariesIoToken,
			"sentry_dsn":         "",
		},
	}
	return secret
}

func (r *CodeReadyAnalyticsReconciler) awsSecret(cr *f8av1alpha1.CodeReadyAnalytics) *corev1.Secret {
	labels := map[string]string{
		"app": cr.Name,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"aws_access_key_id":      []byte(os.Getenv("AWS_KEY")),
			"s3-access-key-id":       []byte(os.Getenv("AWS_KEY")),
			"aws_secret_access_key":  []byte(os.Getenv("AWS_SECRET")),
			"s3-secret-access-key":   []byte(os.Getenv("AWS_SECRET")),
			"sync-s3":                []byte("1"),
			"aws_region":             []byte(cr.Spec.Config.Common.AwsDefaultRegion),
			"s3-bucket-for-analyses": []byte("deepshar-hello"),
		},
	}
	return secret
}

func (r *CodeReadyAnalyticsReconciler) threeScaleSecret(cr *f8av1alpha1.CodeReadyAnalytics) *corev1.Secret {
	labels := map[string]string{
		"app": cr.Name,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "3scale",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"three_scale_account_secret": []byte(cr.Spec.Config.Common.ThreeScaleAccountSecret),
		},
	}
	return secret
}

func (r *CodeReadyAnalyticsReconciler) geminiSecret(cr *f8av1alpha1.CodeReadyAnalytics) *corev1.Secret {
	labels := map[string]string{
		"app": cr.Name,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gemini-server",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"gemini-sa-client-id":     `test`,
			"gemini-sa-client-secret": `secret`,
		},
	}
	return secret
}

func (r *CodeReadyAnalyticsReconciler) bayesianConfigMap(instance *f8av1alpha1.CodeReadyAnalytics) *corev1.ConfigMap {
	labels := map[string]string{
		"app": instance.Name,
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bayesian-config",
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"dynamodb-prefix":        instance.Spec.Config.Common.DynamodbPrefix,
			"auth-url":               instance.Spec.Config.Common.AuthUrl,
			"deployment-prefix":      instance.Spec.Config.Common.DeploymentPrefix,
			"notification-url":       "",
			"s3-bucket-for-analyses": "deepshar-bayesian-core-package-data",
		},
	}
}

// CreateVolume generates a PersistentVolume and Claim.
func (r *CodeReadyAnalyticsReconciler) ensurePV(request reconcile.Request,
	instance *f8av1alpha1.CodeReadyAnalytics,
	dep *corev1.PersistentVolume,
) (*reconcile.Result, error) {

	foundVolume := &corev1.PersistentVolume{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name: dep.Name,
	}, foundVolume)

	r.Log.Info("Instance Name")
	r.Log.Info(instance.Name)

	r.Log.Info("Getting Persistent Volume")
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Vol not found")
		pvDep := r.pvDeployment(instance)

		// Create the deployment
		r.Log.Info("Creating a new PV", "PV.Namespace", pvDep.Namespace, "PV.Name", pvDep.Name)
		err = r.Client.Create(context.TODO(), pvDep)

		if err != nil {
			// Deployment failed
			r.Log.Error(err, "Failed to create new PV", "PV.Namespace", pvDep.Namespace, "PV.Name", pvDep.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		r.Log.Error(err, "Failed to get Persistent Vol")
		return &reconcile.Result{}, err
	}
	// VOLUME was successful
	return nil, nil
}

func (r *CodeReadyAnalyticsReconciler) ensurePVC(request reconcile.Request,
	instance *f8av1alpha1.CodeReadyAnalytics,
	pvcDep *corev1.PersistentVolumeClaim,
) (*reconcile.Result, error) {
	// PVC
	foundPVC := &corev1.PersistentVolumeClaim{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      pvcDep.Name,
		Namespace: instance.Namespace,
	}, foundPVC)

	if err != nil && errors.IsNotFound(err) {
		// Create the PVC
		r.Log.Info("Creating a new PVC", "PVC.Namespace", pvcDep.Namespace, "PVC.Name", pvcDep.Name)
		err = r.Client.Create(context.TODO(), pvcDep)

		if err != nil {
			// PVC failed
			r.Log.Error(err, "Failed to create new PVC", "PVC.Namespace", pvcDep.Namespace, "PVC.Name", pvcDep.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the PVC not existing
		r.Log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}

	r.Log.Info("Successfull Everything")
	return nil, nil
}

func (r *CodeReadyAnalyticsReconciler) ensureBouncerDeployment(request reconcile.Request,
	instance *f8av1alpha1.CodeReadyAnalytics,
	dep *appsv1.StatefulSet,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.StatefulSet{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		r.Log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			r.Log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		}
		if err := r.Client.Update(context.TODO(), instance); err != nil {
			return &reconcile.Result{}, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		r.Log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}
	// Deployment was successful
	return nil, nil
}
