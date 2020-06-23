package codereadyanalytics

import (
	openshiftv1alpha1 "operator/crda-operator/pkg/apis/openshift/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	// "k8s.io/apimachinery/pkg/util/intstr"
)

func bouncerDeploymentName(v *openshiftv1alpha1.CodeReadyAnalytics) string {
	return v.Spec.Pgbouncer.Name
}

func (r *ReconcileCodeReadyAnalytics) bouncerDeployment(v *openshiftv1alpha1.CodeReadyAnalytics) *appsv1.Deployment {
	labels := map[string]string{
		"app": v.Name,
	}
	size := v.Spec.Pgbouncer.Size

	database := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "database",
		},
	}
	host := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "host",
		},
	}
	initialDatabase := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "initial-database",
		},
	}
	username := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "username",
		},
	}
	password := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "password",
		},
	}
	port := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "port",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bouncerDeploymentName(v),
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
						Image:           "quay.io/openshiftio/bayesian-coreapi-pgbouncer:latest",
						ImagePullPolicy: corev1.PullAlways,
						Name:            apiDeploymentName(v),
						Env: []corev1.EnvVar{
							{
								Name:      "POSTGRESQL_DATABASE",
								ValueFrom: database,
							},
							{
								Name:      "POSTGRESQL_INITIAL_DATABASE",
								ValueFrom: initialDatabase,
							},
							{
								Name:      "POSTGRESQL_PASSWORD",
								ValueFrom: password,
							},
							{
								Name:      "POSTGRESQL_USER",
								ValueFrom: username,
							},
							{
								Name:      "POSTGRES_SERVICE_HOST",
								ValueFrom: host,
							},
							{
								Name:      "POSTGRES_SERVICE_PORT",
								ValueFrom: port,
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

func (r *ReconcileCodeReadyAnalytics) pgBouncerService(v *openshiftv1alpha1.CodeReadyAnalytics) *corev1.Service {
	labels := map[string]string{
		"app": v.Name,
	}

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bouncerDeploymentName(v),
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       5432,
				TargetPort: intstr.FromInt(5432),
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(v, s, r.scheme)
	return s
}
