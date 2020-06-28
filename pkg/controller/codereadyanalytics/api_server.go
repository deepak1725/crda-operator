package codereadyanalytics

import (
	"context"
	openshiftv1alpha1 "operator/crda-operator/pkg/apis/openshift/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	// "k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func apiDeploymentName(v *openshiftv1alpha1.CodeReadyAnalytics) string {
	return v.Spec.APIServerService.Name
}

// func (r *ReconcileCodeReadyAnalytics) handleApiServerChanges(v *openshiftv1alpha1.CodeReadyAnalytics) (*reconcile.Result, error) {
// 	found := &appsv1.Deployment{}
// 	err := r.client.Get(context.TODO(), types.NamespacedName{
// 		Name:      apiDeploymentName(v),
// 		Namespace: v.Namespace,
// 	}, found)
// 	if err != nil {
// 		// The deployment may not have been created yet, so requeue
// 		return &reconcile.Result{RequeueAfter: 5 * time.Second}, err
// 	}

// 	size := v.Spec.APIServerService.Size

// 	if size != *found.Spec.Replicas {
// 		found.Spec.Replicas = &size
// 		err = r.client.Update(context.TODO(), found)
// 		if err != nil {
// 			log.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
// 			return &reconcile.Result{}, err
// 		}
// 		// Spec updated - return and requeue
// 		return &reconcile.Result{Requeue: true}, nil
// 	}

// 	return nil, nil
// }

func serverAuthName() string {
	return "apiserver-auth"
}

func (r *ReconcileCodeReadyAnalytics) apiDeployment(v *openshiftv1alpha1.CodeReadyAnalytics) *appsv1.Deployment {
	labels := labels(v, "backend")
	size := v.Spec.APIServerService.Size

	deploymentPrefix := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "bayesian-config"},
			Key:                  "deployment-prefix",
		},
	}
	dbName := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "database",
		},
	}
	dbPassword := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "password",
		},
	}
	dbUsername := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "username",
		},
	}
	awsAccessKeyId := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "aws"},
			Key:                  "aws_access_key_id",
		},
	}
	awsSecretAccessKey := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "aws"},
			Key:                  "aws_secret_access_key",
		},
	}
	threeScaleSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "3scale"},
			Key:                  "three_scale_account_secret",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apiDeploymentName(v),
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           v.Spec.APIServerService.Image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            apiDeploymentName(v),
						Ports: []corev1.ContainerPort{{
							ContainerPort: 5000,
							Name:          "hello",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "BAYESIAN_COMPONENT_TAGGED_COUNT",
								Value: "2",
							},
							{
								Name:  "COMPONENT_ANALYSES_LIMIT",
								Value: "10",
							},
							{
								Name:      "DEPLOYMENT_PREFIX",
								ValueFrom: deploymentPrefix,
							},
							{
								Name:  "WORKER_ADMINISTRATION_REGION",
								Value: "api",
							},
							{
								Name:  "F8_API_BACKBONE_HOST",
								Value: "backbone",
							},
							{
								Name:  "METRICS_ACCUMULATOR_HOST",
								Value: "metrics-accumulator",
							},
							{
								Name:  "METRICS_ACCUMULATOR_PORT",
								Value: "5200",
							},
							{
								Name:  "FUTURES_SESSION_WORKER_COUNT",
								Value: "100",
							},
							{
								Name:  "PGBOUNCER_SERVICE_HOST",
								Value: "bayesian-pgbouncer",
							},
							{
								Name:  "OSIO_AUTH_URL",
								Value: "auth-url",
							},
							{
								Name:  "BAYESIAN_FETCH_PUBLIC_KEY",
								Value: "",
							},
							{
								Name:  "FABRIC8_ANALYTICS_JWT_AUDIENCE",
								Value: "fabric8-online-platform,openshiftio-public,https://prod-preview.openshift.io,https://openshift.io",
							},
							{
								Name:  "SENTRY_DSN",
								Value: "",
							},
							{
								Name:  "INVOKE_API_WORKERS",
								Value: "1",
							},
							{
								Name:  "DISABLE_UNKNOWN_PACKAGE_FLOW",
								Value: "0",
							},
							{
								Name:  "SHOW_TRANSITIVE_REPORT",
								Value: "false",
							},
							{
								Name:  "STACK_ANALYSIS_REQUEST_TIMEOUT",
								Value: "120",
							},
							{
								Name:  "STACK_ANALYSIS_REQUEST_TIMEOUT",
								Value: "0",
							},
							{
								Name:  "FLASK_LOGGING_LEVEL",
								Value: "INFO",
							},
							{
								Name:      "POSTGRESQL_DATABASE",
								ValueFrom: dbName,
							},
							{
								Name:      "POSTGRESQL_USER",
								ValueFrom: dbUsername,
							},
							{
								Name:      "POSTGRESQL_PASSWORD",
								ValueFrom: dbPassword,
							},
							{
								Name:      "AWS_SQS_ACCESS_KEY_ID",
								ValueFrom: awsAccessKeyId,
							},
							{
								Name:      "AWS_SQS_SECRET_ACCESS_KEY",
								ValueFrom: awsSecretAccessKey,
							},
							{
								Name:      "THREESCALE_ACCOUNT_SECRET",
								ValueFrom: threeScaleSecret,
							},
							{
								Name:  "STACK_ANALYSIS_REQUEST_TIMEOUT",
								Value: "120",
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.scheme)
	return dep
}

func (r *ReconcileCodeReadyAnalytics) apiService(v *openshiftv1alpha1.CodeReadyAnalytics) *corev1.Service {
	labels := labels(v, "backend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      apiDeploymentName(v),
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       5000,
				TargetPort: intstr.FromInt(5000),
				NodePort:   32000,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}

func (r *ReconcileCodeReadyAnalytics) updateBackendStatus(v *openshiftv1alpha1.CodeReadyAnalytics) error {
	v.Status.APIServerService.Image = v.Spec.APIServerService.Image
	err := r.client.Status().Update(context.TODO(), v)
	return err
}
