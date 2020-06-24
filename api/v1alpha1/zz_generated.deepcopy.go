// +build !ignore_autogenerated

/*
Copyright 2020 O.Yuanying <yuanying@fraction.jp>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroup) DeepCopyInto(out *InitiatorGroup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroup.
func (in *InitiatorGroup) DeepCopy() *InitiatorGroup {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InitiatorGroup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroupBinding) DeepCopyInto(out *InitiatorGroupBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroupBinding.
func (in *InitiatorGroupBinding) DeepCopy() *InitiatorGroupBinding {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroupBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InitiatorGroupBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroupBindingList) DeepCopyInto(out *InitiatorGroupBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InitiatorGroupBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroupBindingList.
func (in *InitiatorGroupBindingList) DeepCopy() *InitiatorGroupBindingList {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroupBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InitiatorGroupBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroupBindingSpec) DeepCopyInto(out *InitiatorGroupBindingSpec) {
	*out = *in
	out.TargetRef = in.TargetRef
	out.InitiatorGroupRef = in.InitiatorGroupRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroupBindingSpec.
func (in *InitiatorGroupBindingSpec) DeepCopy() *InitiatorGroupBindingSpec {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroupBindingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroupBindingStatus) DeepCopyInto(out *InitiatorGroupBindingStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroupBindingStatus.
func (in *InitiatorGroupBindingStatus) DeepCopy() *InitiatorGroupBindingStatus {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroupBindingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroupInitiatorNameStrategy) DeepCopyInto(out *InitiatorGroupInitiatorNameStrategy) {
	*out = *in
	if in.InitiatorNamePrefix != nil {
		in, out := &in.InitiatorNamePrefix, &out.InitiatorNamePrefix
		*out = new(string)
		**out = **in
	}
	if in.AnnotationKey != nil {
		in, out := &in.AnnotationKey, &out.AnnotationKey
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroupInitiatorNameStrategy.
func (in *InitiatorGroupInitiatorNameStrategy) DeepCopy() *InitiatorGroupInitiatorNameStrategy {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroupInitiatorNameStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroupList) DeepCopyInto(out *InitiatorGroupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InitiatorGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroupList.
func (in *InitiatorGroupList) DeepCopy() *InitiatorGroupList {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InitiatorGroupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroupReference) DeepCopyInto(out *InitiatorGroupReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroupReference.
func (in *InitiatorGroupReference) DeepCopy() *InitiatorGroupReference {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroupReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroupSpec) DeepCopyInto(out *InitiatorGroupSpec) {
	*out = *in
	if in.Addresses != nil {
		in, out := &in.Addresses, &out.Addresses
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.InitiatorNameStrategy.DeepCopyInto(&out.InitiatorNameStrategy)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroupSpec.
func (in *InitiatorGroupSpec) DeepCopy() *InitiatorGroupSpec {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InitiatorGroupStatus) DeepCopyInto(out *InitiatorGroupStatus) {
	*out = *in
	if in.Addresses != nil {
		in, out := &in.Addresses, &out.Addresses
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Initiators != nil {
		in, out := &in.Initiators, &out.Initiators
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InitiatorGroupStatus.
func (in *InitiatorGroupStatus) DeepCopy() *InitiatorGroupStatus {
	if in == nil {
		return nil
	}
	out := new(InitiatorGroupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Target) DeepCopyInto(out *Target) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Target.
func (in *Target) DeepCopy() *Target {
	if in == nil {
		return nil
	}
	out := new(Target)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Target) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetActual) DeepCopyInto(out *TargetActual) {
	*out = *in
	if in.LUNs != nil {
		in, out := &in.LUNs, &out.LUNs
		*out = make([]TargetLUN, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Accounts != nil {
		in, out := &in.Accounts, &out.Accounts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ACLs != nil {
		in, out := &in.ACLs, &out.ACLs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetActual.
func (in *TargetActual) DeepCopy() *TargetActual {
	if in == nil {
		return nil
	}
	out := new(TargetActual)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetCondition) DeepCopyInto(out *TargetCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetCondition.
func (in *TargetCondition) DeepCopy() *TargetCondition {
	if in == nil {
		return nil
	}
	out := new(TargetCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetLUN) DeepCopyInto(out *TargetLUN) {
	*out = *in
	if in.BSType != nil {
		in, out := &in.BSType, &out.BSType
		*out = new(string)
		**out = **in
	}
	if in.BSOpts != nil {
		in, out := &in.BSOpts, &out.BSOpts
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetLUN.
func (in *TargetLUN) DeepCopy() *TargetLUN {
	if in == nil {
		return nil
	}
	out := new(TargetLUN)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetList) DeepCopyInto(out *TargetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Target, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetList.
func (in *TargetList) DeepCopy() *TargetList {
	if in == nil {
		return nil
	}
	out := new(TargetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TargetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetReference) DeepCopyInto(out *TargetReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetReference.
func (in *TargetReference) DeepCopy() *TargetReference {
	if in == nil {
		return nil
	}
	out := new(TargetReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetSpec) DeepCopyInto(out *TargetSpec) {
	*out = *in
	if in.LUNs != nil {
		in, out := &in.LUNs, &out.LUNs
		*out = make([]TargetLUN, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetSpec.
func (in *TargetSpec) DeepCopy() *TargetSpec {
	if in == nil {
		return nil
	}
	out := new(TargetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetStatus) DeepCopyInto(out *TargetStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]TargetCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ObservedState != nil {
		in, out := &in.ObservedState, &out.ObservedState
		*out = new(TargetActual)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetStatus.
func (in *TargetStatus) DeepCopy() *TargetStatus {
	if in == nil {
		return nil
	}
	out := new(TargetStatus)
	in.DeepCopyInto(out)
	return out
}
