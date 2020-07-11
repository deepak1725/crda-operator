package codereadyanalytics

import (
	"context"
	openshiftv1alpha1 "operator/crda-operator/pkg/apis/openshift/v1alpha1"
	"os"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileCodeReadyAnalytics) ensureDeployment(request reconcile.Request,
	instance *openshiftv1alpha1.CodeReadyAnalytics,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}
	// Deployment was successful
	return nil, nil
}

func (r *ReconcileCodeReadyAnalytics) ensureService(request reconcile.Request,
	instance *openshiftv1alpha1.CodeReadyAnalytics,
	s *corev1.Service,
) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Service")
		return &reconcile.Result{}, err
	}
	// Creation was successful
	return nil, nil
}

func (r *ReconcileCodeReadyAnalytics) ensureSecret(request reconcile.Request,
	instance *openshiftv1alpha1.CodeReadyAnalytics,
	secret *corev1.Secret,
) (*reconcile.Result, error) {
	found := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      secret.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the secret
		log.Info("Creating a new secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.client.Create(context.TODO(), secret)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the secret not existing
		log.Error(err, "Failed to get Secret")
		return &reconcile.Result{}, err
	}
	// Creation was successful
	return nil, nil
}

func (r *ReconcileCodeReadyAnalytics) ensureConfigMap(request reconcile.Request,
	instance *openshiftv1alpha1.CodeReadyAnalytics,
	configMap *corev1.ConfigMap,
) (*reconcile.Result, error) {

	// Check if this ConfigMap already exists
	foundMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, foundMap)
	if err != nil && errors.IsNotFound(err) {
		err = r.client.Create(context.TODO(), configMap)
		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "configMap.Name", configMap.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the secret not existing
		log.Error(err, "Failed to get ConfigMap")
		return &reconcile.Result{}, err
	}
	// Creation was successful
	return nil, nil
}

func labels(v *openshiftv1alpha1.CodeReadyAnalytics, tier string) map[string]string {
	return map[string]string{
		"app":          "analytics",
		"analytics_cr": v.Name,
		"tier":         tier,
	}
}

func (r *ReconcileCodeReadyAnalytics) postgresSecret(cr *openshiftv1alpha1.CodeReadyAnalytics) *corev1.Secret {
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

func (r *ReconcileCodeReadyAnalytics) workerSecret(cr *openshiftv1alpha1.CodeReadyAnalytics) *corev1.Secret {
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
			"github-token":         cr.Spec.Config.Common.GithubToken,
			"libraries-io-token":   cr.Spec.Config.Common.LibrariesIoToken,
			"sentry_dsn":            "",
		},
	}
	return secret
}


func (r *ReconcileCodeReadyAnalytics) awsSecret(cr *openshiftv1alpha1.CodeReadyAnalytics) *corev1.Secret {
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
			"aws_access_key_id":     	[]byte(os.Getenv("AWS_KEY")),
			"s3-access-key-id":     	[]byte(os.Getenv("AWS_KEY")),
			"aws_secret_access_key": 	[]byte(os.Getenv("AWS_SECRET")),
			"s3-secret-access-key": 	[]byte(os.Getenv("AWS_SECRET")),
			"sync-s3": 					[]byte("1"),
			"aws_region": 				[]byte(cr.Spec.Config.Common.AwsDefaultRegion),
			"s3-bucket-for-analyses": 	[]byte("deepshar-hello"),
		},
	}
	return secret
}


func (r *ReconcileCodeReadyAnalytics) threeScaleSecret(cr *openshiftv1alpha1.CodeReadyAnalytics) *corev1.Secret {
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

func (r *ReconcileCodeReadyAnalytics) geminiSecret(cr *openshiftv1alpha1.CodeReadyAnalytics) *corev1.Secret {
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
			"gemini-sa-client-id": `test`,
			"gemini-sa-client-secret": `secret`,
		},
	}
	return secret
}

func (r *ReconcileCodeReadyAnalytics) bayesianConfigMap(instance *openshiftv1alpha1.CodeReadyAnalytics) *corev1.ConfigMap {
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
			"dynamodb-prefix":   instance.Spec.Config.Common.DynamodbPrefix,
			"auth-url":          instance.Spec.Config.Common.AuthUrl,
			"deployment-prefix": instance.Spec.Config.Common.DeploymentPrefix,
			"notification-url": "",
			"s3-bucket-for-analyses": "deepshar-bayesian-core-package-data",
		},
	}
}

// CreateVolume generates a PersistentVolume and Claim.
func (r *ReconcileCodeReadyAnalytics) ensurePV(request reconcile.Request, 
	instance *openshiftv1alpha1.CodeReadyAnalytics,
	dep *corev1.PersistentVolume,
	) (*reconcile.Result, error) {

	foundVolume := &corev1.PersistentVolume{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
	}, foundVolume)
	
	log.Info("Instance Name")
	log.Info(instance.Name)

	log.Info("Getting Persistent Volume")
	if err != nil && errors.IsNotFound(err) {
		log.Info("Vol not found")
		pvDep := r.pvDeployment(instance)

		// Create the deployment
		log.Info("Creating a new PV", "PV.Namespace", pvDep.Namespace, "PV.Name", pvDep.Name)
		err = r.client.Create(context.TODO(), pvDep)

		if err != nil {
			// Deployment failed
			log.Error(err, "Failed to create new PV", "PV.Namespace", pvDep.Namespace, "PV.Name", pvDep.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get Persistent Vol")
		return &reconcile.Result{}, err
	}
	// VOLUME was successful
	return nil, nil
}

func (r *ReconcileCodeReadyAnalytics) ensurePVC(request reconcile.Request,
	instance *openshiftv1alpha1.CodeReadyAnalytics,
	pvcDep *corev1.PersistentVolumeClaim,
) (*reconcile.Result, error) {
	// PVC
	foundPVC := &corev1.PersistentVolumeClaim{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      pvcDep.Name,
		Namespace: instance.Namespace,
	}, foundPVC)

	if err != nil && errors.IsNotFound(err) {
		// Create the PVC
		log.Info("Creating a new PVC", "PVC.Namespace", pvcDep.Namespace, "PVC.Name", pvcDep.Name)
		err = r.client.Create(context.TODO(), pvcDep)

		if err != nil {
			// PVC failed
			log.Error(err, "Failed to create new PVC", "PVC.Namespace", pvcDep.Namespace, "PVC.Name", pvcDep.Name)
			return &reconcile.Result{}, err
		}
	} else if err != nil {
		// Error that isn't due to the PVC not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}
	
	log.Info("Successfull Everything")
	return nil, nil
}





func (r *ReconcileCodeReadyAnalytics) ensureBouncerDeployment(request reconcile.Request,
	instance *openshiftv1alpha1.CodeReadyAnalytics,
	dep *appsv1.StatefulSet,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.StatefulSet{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		}
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return &reconcile.Result{}, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}
	// Deployment was successful
	return nil, nil
}
