// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackboneServiceType) DeepCopyInto(out *BackboneServiceType) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackboneServiceType.
func (in *BackboneServiceType) DeepCopy() *BackboneServiceType {
	if in == nil {
		return nil
	}
	out := new(BackboneServiceType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CodeReadyAnalytics) DeepCopyInto(out *CodeReadyAnalytics) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CodeReadyAnalytics.
func (in *CodeReadyAnalytics) DeepCopy() *CodeReadyAnalytics {
	if in == nil {
		return nil
	}
	out := new(CodeReadyAnalytics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CodeReadyAnalytics) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CodeReadyAnalyticsList) DeepCopyInto(out *CodeReadyAnalyticsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CodeReadyAnalytics, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CodeReadyAnalyticsList.
func (in *CodeReadyAnalyticsList) DeepCopy() *CodeReadyAnalyticsList {
	if in == nil {
		return nil
	}
	out := new(CodeReadyAnalyticsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CodeReadyAnalyticsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CodeReadyAnalyticsSpec) DeepCopyInto(out *CodeReadyAnalyticsSpec) {
	*out = *in
	out.Config = in.Config
	out.BackboneService = in.BackboneService
	out.APIServerService = in.APIServerService
	out.Pgbouncer = in.Pgbouncer
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CodeReadyAnalyticsSpec.
func (in *CodeReadyAnalyticsSpec) DeepCopy() *CodeReadyAnalyticsSpec {
	if in == nil {
		return nil
	}
	out := new(CodeReadyAnalyticsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CodeReadyAnalyticsStatus) DeepCopyInto(out *CodeReadyAnalyticsStatus) {
	*out = *in
	out.BackboneService = in.BackboneService
	out.APIServerService = in.APIServerService
	out.Pgbouncer = in.Pgbouncer
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CodeReadyAnalyticsStatus.
func (in *CodeReadyAnalyticsStatus) DeepCopy() *CodeReadyAnalyticsStatus {
	if in == nil {
		return nil
	}
	out := new(CodeReadyAnalyticsStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonType) DeepCopyInto(out *CommonType) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonType.
func (in *CommonType) DeepCopy() *CommonType {
	if in == nil {
		return nil
	}
	out := new(CommonType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigType) DeepCopyInto(out *ConfigType) {
	*out = *in
	out.Common = in.Common
	out.Database = in.Database
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigType.
func (in *ConfigType) DeepCopy() *ConfigType {
	if in == nil {
		return nil
	}
	out := new(ConfigType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DatabaseType) DeepCopyInto(out *DatabaseType) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DatabaseType.
func (in *DatabaseType) DeepCopy() *DatabaseType {
	if in == nil {
		return nil
	}
	out := new(DatabaseType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PgbouncerType) DeepCopyInto(out *PgbouncerType) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PgbouncerType.
func (in *PgbouncerType) DeepCopy() *PgbouncerType {
	if in == nil {
		return nil
	}
	out := new(PgbouncerType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerServiceType) DeepCopyInto(out *ServerServiceType) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerServiceType.
func (in *ServerServiceType) DeepCopy() *ServerServiceType {
	if in == nil {
		return nil
	}
	out := new(ServerServiceType)
	in.DeepCopyInto(out)
	return out
}
