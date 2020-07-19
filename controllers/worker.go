package controllers

import (
	f8av1alpha1 "github.com/deepak1725/crda-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func workerDeploymentName(v *f8av1alpha1.CodeReadyAnalytics) string {
	return v.Spec.Worker.Name
}

func (r *CodeReadyAnalyticsReconciler) workerDeployment(v *f8av1alpha1.CodeReadyAnalytics) *appsv1.Deployment {
	labels := labels(v, "worker")
	size := v.Spec.Worker.Size

	depPrefix := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "bayesian-config"},
			Key:                  "deployment-prefix",
		},
	}
	sentryDsn := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "worker"},
			Key:                  "sentry_dsn",
		},
	}
	postgresDb := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "database",
		},
	}
	postgresPwd := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "password",
		},
	}
	postgresUser := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "username",
		},
	}
	pgInitDb := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "coreapi-postgres"},
			Key:                  "initial-database",
		},
	}
	awsSqsKeyId := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "aws"},
			Key:                  "aws_access_key_id",
		},
	}
	awsSqsSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "aws"},
			Key:                  "aws_secret_access_key",
		},
	}
	awsS3KeyId := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "aws"},
			Key:                  "aws_access_key_id",
		},
	}
	awsS3Secret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "aws"},
			Key:                  "aws_secret_access_key",
		},
	}
	syncS3 := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "aws"},
			Key:                  "sync-s3",
		},
	}
	gitToken := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "worker"},
			Key:                  "github-token",
		},
	}
	libIoToken := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "worker"},
			Key:                  "libraries-io-token",
		},
	}
	geminiSaClientId := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "gemini-server"},
			Key:                  "gemini-sa-client-id",
		},
	}
	geminiSaClientSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "gemini-server"},
			Key:                  "gemini-sa-client-secret",
		},
	}
	authServiceHost := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "bayesian-config"},
			Key:                  "auth-url",
		},
	}
	notificationUrl := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "bayesian-config"},
			Key:                  "notification-url",
		},
	}
	s3AnalysesBucket := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "bayesian-config"},
			Key:                  "s3-bucket-for-analyses",
		},
	}
	f8aAuthHost := &corev1.EnvVarSource{
		ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "bayesian-config"},
			Key:                  "auth-url",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workerDeploymentName(v),
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
						Image:           v.Spec.Worker.Image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            workerDeploymentName(v),
						Lifecycle: &corev1.Lifecycle{
							PostStart: &corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{"worker-pre-hook.sh"},
								},
							},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{"worker-liveness.sh"},
								},
							},
							InitialDelaySeconds: 60,
							PeriodSeconds:       60,
							TimeoutSeconds:      30,
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								Exec: &corev1.ExecAction{
									Command: []string{"worker-readiness.sh"},
								},
							},
							InitialDelaySeconds: 60,
							PeriodSeconds:       60,
							TimeoutSeconds:      30,
						},
						Env: []corev1.EnvVar{
							{
								Name:  "OPENSHIFT_DEPLOYMENT",
								Value: "1",
							},
							{
								Name:      "DEPLOYMENT_PREFIX",
								ValueFrom: depPrefix,
							},
							{
								Name:  "WORKER_ADMINISTRATION_REGION",
								Value: "api",
							},
							{
								Name:  "WORKER_EXCLUDE_QUEUES",
								Value: "GraphImporterTask",
							},
							{
								Name:  "WORKER_INCLUDE_QUEUES",
								Value: "",
							},
							{
								Name:  "WORKER_RUN_DB_MIGRATIONS",
								Value: "1",
							},
							{
								Name:      "SENTRY_DSN",
								ValueFrom: sentryDsn,
							},
							{
								Name:      "POSTGRESQL_DATABASE",
								ValueFrom: postgresDb,
							},
							{
								Name:      "POSTGRESQL_PASSWORD",
								ValueFrom: postgresPwd,
							},
							{
								Name:      "POSTGRESQL_USER",
								ValueFrom: postgresUser,
							},
							{
								Name:      "POSTGRESQL_INITIAL_DATABASE",
								ValueFrom: pgInitDb,
							},
							{
								Name:      "AWS_SQS_ACCESS_KEY_ID",
								ValueFrom: awsSqsKeyId,
							},
							{
								Name:      "AUTH_SERVICE_HOST",
								ValueFrom: authServiceHost,
							},
							{
								Name:      "NOTIFICATION_SERVICE_HOST",
								ValueFrom: notificationUrl,
							},
							{
								Name:      "AWS_SQS_SECRET_ACCESS_KEY",
								ValueFrom: awsSqsSecret,
							},
							{
								Name:      "AWS_S3_ACCESS_KEY_ID",
								ValueFrom: awsS3KeyId,
							},
							{
								Name:      "AWS_S3_SECRET_ACCESS_KEY",
								ValueFrom: awsS3Secret,
							},
							{
								Name:      "BAYESIAN_SYNC_S3",
								ValueFrom: syncS3,
							},
							{
								Name:      "GITHUB_TOKEN",
								ValueFrom: gitToken,
							},
							{
								Name:      "LIBRARIES_IO_TOKEN",
								ValueFrom: libIoToken,
							},
							{
								Name:  "PGBOUNCER_SERVICE_HOST",
								Value: "bayesian-pgbouncer",
							},
							{
								Name:  "PGM_SERVICE_HOST",
								Value: "bayesian-kronos",
							},
							{
								Name:  "MAX_COMPANION_PACKAGES",
								Value: "4",
							},
							{
								Name:  "MAX_ALTERNATE_PACKAGES",
								Value: "2",
							},
							{
								Name:  "OUTLIER_THRESHOLD",
								Value: "0.88",
							},
							{
								Name:  "UNKNOWN_PACKAGES_THRESHOLD",
								Value: "0.3",
							},
							{
								Name:  "LICENSE_SERVICE_HOST",
								Value: "f8a-license-analysis",
							},
							{
								Name:  "PGM_SERVICE_PORT",
								Value: "6006",
							},
							{
								Name:  "LICENSE_SERVICE_PORT",
								Value: "6162",
							},
							{
								Name:  "F8A_SERVER_SERVICE_HOST",
								Value: "bayesian-api",
							},
							{
								Name:  "GIT_COMMITTER_NAME",
								Value: "rhdt-dep-analytics",
							},
							{
								Name:  "GIT_COMMITTER_EMAIL",
								Value: "rhdt-dep-analytics@example.com",
							},
							{
								Name:  "F8_API_BACKBONE_HOST",
								Value: "http://f8a-server-backbone:5000",
							},
							{
								Name:  "AWS_S3_BUCKET_NAME",
								Value: "http://f8a-gemini-server:5000",
							},
							{
								Name:      "AWS_S3_BUCKET_NAME",
								ValueFrom: s3AnalysesBucket,
							},
							{
								Name:  "RABBITMQ_SERVICE_SERVICE_HOST",
								Value: "bayesian-broker",
							},
							{
								Name:      "GEMINI_SA_CLIENT_ID",
								ValueFrom: geminiSaClientId,
							},
							{
								Name:      "GEMINI_SA_CLIENT_SECRET",
								ValueFrom: geminiSaClientSecret,
							},
							{
								Name:      "F8A_AUTH_SERVICE_HOST",
								ValueFrom: f8aAuthHost,
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
