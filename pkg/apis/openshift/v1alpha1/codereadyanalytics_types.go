package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CodeReadyAnalyticsSpec defines the desired state of CodeReadyAnalytics
type CodeReadyAnalyticsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// CodeReadyAnalyticsStatus defines the observed state of CodeReadyAnalytics
type CodeReadyAnalyticsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CodeReadyAnalytics is the Schema for the codereadyanalytics API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=codereadyanalytics,scope=Namespaced
type CodeReadyAnalytics struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CodeReadyAnalyticsSpec   `json:"spec,omitempty"`
	Status CodeReadyAnalyticsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CodeReadyAnalyticsList contains a list of CodeReadyAnalytics
type CodeReadyAnalyticsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CodeReadyAnalytics `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CodeReadyAnalytics{}, &CodeReadyAnalyticsList{})
}
