package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BackboneServiceType defines Backbone Service Specs
type BackboneServiceType struct {
	Image         string `json:"image,omitempty"`
	Size          int32  `json:"size,omitempty"`
	ContainerPort int32  `json:"containerPort,omitempty"`
	Name          string `json:"name,omitempty"`
}

// ServerServiceType defines API Server Service Specs
type ServerServiceType struct {
	Image         string `json:"image,omitempty"`
	Size          int32  `json:"size,omitempty"`
	ContainerPort int32  `json:"containerPort,omitempty"`
	Name          string `json:"name,omitempty"`
}

// PgbouncerType defines Specs for Pgbouncer Service
type PgbouncerType struct {
	Name string `json:"name,omitempty"`
	Size int32  `json:"size,omitempty"`
	Image string `json:"image,omitempty"`
}

// GremlinType defines Specs for Gremlin Service
type GremlinType struct {
	Name string `json:"name,omitempty"`
	Size int32  `json:"size,omitempty"`
	Image string `json:"image,omitempty"`
	DynamoDbEndpoint string `json:"dynamoDbEndpoint,omitempty"`
	Resources ResourceType `json:"resources,omitempty"`
}

// DatabaseType defines Database Config Specs
type DatabaseType struct {
	DbName          string `json:"dbName,omitempty"`
	Host            string `json:"host,omitempty"`
	InitialDatabase string `json:"initialDatabase,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	Port            string `json:"port,omitempty"`
}

// CommonType defines common Config Specs
type CommonType struct {
	AuthUrl                 string `json:"authUrl,omitempty"`
	DeploymentPrefix        string `json:"deploymentPrefix,omitempty"`
	DynamodbPrefix          string `json:"dynamodbPrefix,omitempty"`
	ThreeScaleAccountSecret string `json:"threeScaleAccountSecret,omitempty"`
	AwsAccessKeyId          string `json:"awsAccessKeyId,omitempty"`
	AwsSecretAccessKey      string `json:"awsSecretAccessKey,omitempty"`
	AwsDefaultRegion      	string `json:"awsDefaultRegion,omitempty"`
	GithubToken      		string `json:"githubToken,omitempty"`
	LibrariesIoToken      	string `json:"librariesIoToken,omitempty"`
}

// ConfigType defines Configuration Specs
type ConfigType struct {
	Common   CommonType   `json:"common,omitempty"`
	Database DatabaseType `json:"database,omitempty"`
}

// LimitType defines Resource Types of Gremlin
type LimitType struct {
	Memory   string   		`json:"memory,omitempty"`
	CPU 	 string 		`json:"cpu,omitempty"`
}

// ResourceType defines Resource Types of Gremlin
type ResourceType struct {
	Requests   LimitType	    `json:"requests,omitempty"`
	Limits 	   LimitType 		`json:"limits,omitempty"`
}

// WorkerType defines Worker Types
type WorkerType struct {
	Name   		string	     `json:"name,omitempty"`
	Image 	   	string 		 `json:"image,omitempty"`
	Size        int32  		 `json:"size,omitempty"`
	Resources	ResourceType `json:"resources,omitempty"`
}

// CodeReadyAnalyticsSpec defines the desired state of CodeReadyAnalytics
type CodeReadyAnalyticsSpec struct {
	// Fields Required for Operator Functioning.
	Config           ConfigType          `json:"config,omitempty"`
	BackboneService  BackboneServiceType `json:"backbone,omitempty"`
	APIServerService ServerServiceType   `json:"api-server,omitempty"`
	Pgbouncer        PgbouncerType       `json:"pgbouncer,omitempty"`
	Gremlin          GremlinType       	`json:"gremlin,omitempty"`
	Worker           WorkerType       	`json:"worker,omitempty"`
}

// CodeReadyAnalyticsStatus defines the observed state of CodeReadyAnalytics
type CodeReadyAnalyticsStatus struct {
	// Status fields
	BackboneService  BackboneServiceType `json:"backbone,omitempty"`
	APIServerService ServerServiceType   `json:"api-server,omitempty"`
	Pgbouncer        PgbouncerType       `json:"pgbouncer,omitempty"`
	Gremlin          GremlinType       	`json:"gremlin,omitempty"`
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
