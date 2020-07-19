package controllers

import (
	f8av1alpha1 "github.com/deepak1725/crda-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func gremlinDeploymentName(v *f8av1alpha1.CodeReadyAnalytics) string {
	return v.Spec.Gremlin.Name
}

func (r *CodeReadyAnalyticsReconciler) gremlinDeployment(v *f8av1alpha1.CodeReadyAnalytics) *appsv1.Deployment {
	labels := map[string]string{
		"service": v.Spec.Gremlin.Name,
	}
	size := v.Spec.Gremlin.Size

	dynamoPrefix := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "bayesian-config"},
			Key:                  "dynamodb-prefix",
		},
	}
	awsDefaultRegion := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "aws"},
			Key:                  "aws_region",
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

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gremlinDeploymentName(v),
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
						Image:           v.Spec.Gremlin.Image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            gremlinDeploymentName(v),
						Lifecycle: &corev1.Lifecycle{
							PostStart: &corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{"post-hook.sh"},
								},
							},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								TCPSocket: &corev1.TCPSocketAction{
									Port: intstr.FromInt(8182),
								},
							},

							InitialDelaySeconds: 60,
							PeriodSeconds:       60,
							TimeoutSeconds:      30,
							FailureThreshold:    3,
							SuccessThreshold:    1,
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								TCPSocket: &corev1.TCPSocketAction{
									Port: intstr.FromInt(8182),
								},
							},
							InitialDelaySeconds: 60,
							PeriodSeconds:       60,
							TimeoutSeconds:      30,
							FailureThreshold:    3,
							SuccessThreshold:    1,
						},
						Env: []corev1.EnvVar{
							{
								Name:  "REST",
								Value: "1",
							},
							{
								Name:  "DYNAMODB_CLIENT_ENDPOINT",
								Value: v.Spec.Gremlin.DynamoDbEndpoint,
							},
							{
								Name:      "DYNAMODB_PREFIX",
								ValueFrom: dynamoPrefix,
							},
							{
								Name:  "DYNAMODB_CLIENT_CREDENTIALS_CLASS_NAME",
								Value: "com.amazonaws.auth.DefaultAWSCredentialsProviderChain",
							},
							{
								Name:  "JAVA_OPTIONS",
								Value: "",
							},
							{
								Name:      "AWS_ACCESS_KEY_ID",
								ValueFrom: awsAccessKeyId,
							},
							{
								Name:      "AWS_SECRET_ACCESS_KEY",
								ValueFrom: awsSecretAccessKey,
							},
							{
								Name:      "AWS_DEFAULT_REGION",
								ValueFrom: awsDefaultRegion,
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

func (r *CodeReadyAnalyticsReconciler) gremlinService(v *f8av1alpha1.CodeReadyAnalytics) *corev1.Service {
	labels := map[string]string{
		"service": v.Spec.Gremlin.Name,
	}
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gremlinDeploymentName(v),
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       8182,
				TargetPort: intstr.FromInt(8182),
				NodePort:   32500,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(v, s, r.Scheme)
	return s
}
