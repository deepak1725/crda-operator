package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type BackboneServiceType struct {
	// Backbone Service Specs
	Image string `json:"image,omitempty"`
	Size  int32  `json:"size,omitempty"`
}

type ServerServiceType struct {
	// API Server Service Specs
	Image string `json:"image,omitempty"`
	Size  int32  `json:"size,omitempty"`
}

// CodeReadyAnalyticsSpec defines the desired state of CodeReadyAnalytics
type CodeReadyAnalyticsSpec struct {
	// Fields Required for Operator Functioning.
	BackboneService  BackboneServiceType `json:"backbone,omitempty"`
	APIServerService ServerServiceType   `json:"api-server,omitempty"`
}

// CodeReadyAnalyticsStatus defines the observed state of CodeReadyAnalytics
type CodeReadyAnalyticsStatus struct {
	// Status fields
	BackboneService  BackboneServiceType `json:"backbone,omitempty"`
	APIServerService ServerServiceType   `json:"api-server,omitempty"`
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