package codereadyanalytics

import (
	openshiftv1alpha1 "operator/crda-operator/pkg/apis/openshift/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/api/resource"
	"fmt"
)

func bouncerDeploymentName(v *openshiftv1alpha1.CodeReadyAnalytics) string {
	return v.Spec.Pgbouncer.Name
}

const (
	DiskSize            = 1 * 1000 * 1000 * 10  //10 MB
	AppVolumeName       = "pgdata"
	AppVolumeMountPath  = "/pgdata"
	HostProvisionerPath = "/tmp/hostpath-provisioner"
	ImagePullPolicy     = corev1.PullAlways
)	

var (
	storageClassName              = "standard"
	diskSize                      = *resource.NewQuantity(DiskSize, resource.DecimalSI)
	terminationGracePeriodSeconds = int64(10)
	accessMode                    = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
	resourceList                  = corev1.ResourceList{corev1.ResourceStorage: diskSize}
)



func (r *ReconcileCodeReadyAnalytics) pvcDeployment(cr *openshiftv1alpha1.CodeReadyAnalytics) (*corev1.PersistentVolumeClaim) {
	labels := map[string]string{
		"app": cr.Name,
	}
	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      AppVolumeName,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: accessMode,
			VolumeName:  AppVolumeName,
			Resources:   corev1.ResourceRequirements{Requests: resourceList},
		},
	}
	addOwnerRefToObject(pvc, asOwner(cr))
	return pvc
}

func (r *ReconcileCodeReadyAnalytics) pvDeployment(cr *openshiftv1alpha1.CodeReadyAnalytics) (*corev1.PersistentVolume) {
	log.Info("Creating Deployment")

	labels := map[string]string{
		"app": "pv",
	}
	dep := &corev1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      AppVolumeName,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: storageClassName,
			AccessModes:      accessMode,
			Capacity:         resourceList,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: fmt.Sprintf("%s/%s", HostProvisionerPath, cr.ObjectMeta.Name),
				},
			},
		},
	}
	addOwnerRefToObject(dep, asOwner(cr))
	return dep
}

func (r *ReconcileCodeReadyAnalytics) bouncerDeployment(cr *openshiftv1alpha1.CodeReadyAnalytics) *appsv1.StatefulSet {
	log.Info("Creating a new Statefulset")
	labels := labels(cr, "database")
	size := cr.Spec.Pgbouncer.Size

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
	statefulset := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      bouncerDeploymentName(cr),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:   &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					Containers: []corev1.Container{
						corev1.Container{
							Name:            cr.Spec.Pgbouncer.Name,
							Image:           cr.Spec.Pgbouncer.Image,
							ImagePullPolicy: ImagePullPolicy,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      AppVolumeName,
									MountPath: AppVolumeMountPath,
							},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"health-check-probe.sh"},
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       60,
								TimeoutSeconds: 5,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{"health-check-probe.sh"},
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       60,
								TimeoutSeconds: 5,
							},
							Env: []corev1.EnvVar{
								{
									Name:      "POSTGRES_DB",
									ValueFrom: database,
								},
								{
									Name:      "PGDATA",
									Value: "/var/lib/postgresql/data/pgdata",
								},
								{
									Name:      "POSTGRESQL_INITIAL_DATABASE",
									ValueFrom: initialDatabase,
								},
								{
									Name:      "POSTGRES_PASSWORD",
									ValueFrom: password,
								},
								{
									Name:      "POSTGRES_USER",
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
						},
					},
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: AppVolumeName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: AppVolumeName,
								},
							},
						},
					},
				},
			},
		},
	}
	addOwnerRefToObject(statefulset, asOwner(cr))
	return statefulset
}
// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// asOwner returns an OwnerReference set as the stateful CR
func asOwner(hs *openshiftv1alpha1.CodeReadyAnalytics) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: hs.APIVersion,
		Kind:       hs.Kind,
		Name:       hs.Name,
		UID:        hs.UID,
		Controller: &trueVar,
	}
}

func (r *ReconcileCodeReadyAnalytics) bouncerService(cr *openshiftv1alpha1.CodeReadyAnalytics) (*corev1.Service) {
	labels := labels(cr, "database")
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Pgbouncer.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Selector:  labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       5432,
				TargetPort: intstr.FromInt(5432),
				NodePort:   31500,
			}},
		},
	}
	addOwnerRefToObject(service, asOwner(cr))
	return service
}