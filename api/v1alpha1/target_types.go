/*


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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TargetSpec defines the desired state of Target
type TargetSpec struct {
	// NodeName is a node name where the target will be placed.
	// +kubebuilder:validation:Required
	NodeName string `json:"targetNodeName,omitempty"`

	// IQN is an iqn of the target
	// +kubebuilder:validation:Required
	IQN string `json:"iqn,omitempty"`

	// LUNs is a list of LUNs
	LUNs []TargetLUN `json:"luns,omitempty"`
}

// TargetLun is the specification of LUN
type TargetLUN struct {
	// LUN is an id of the LUN
	LUN int32 `json:"lun,omitempty"`

	// BackingStore is a path of the backing store
	// +kubebuilder:validation:Required
	BackingStore string `json:"backingStore,omitempty"`

	// BSType is a backing store type
	BSType *string `json:"bsType,omitempty"`

	// BSOpts is a options for backing store
	BSOpts *string `json:"bsOpts,omitempty"`
}

// TargetActual is the observed information of Target
type TargetActual struct {
	// TID is the observed tid
	TID int32 `json:"tid,omitempty"`

	// IQN is the observed IQN
	IQN string `json:"iqn,omitempty"`

	// LUNs is the observed LUNs
	LUNs []TargetLUN `json:"luns,omitempty"`

	// Accounts is the observed Accounts
	Accounts []string `json:"accounts,omitempty"`

	// ACLs is the observed ACLs
	ACLs []string `json:"acls,omitempty"`
}

// TargetStatus defines the observed state of Target
type TargetStatus struct {
	// Conditions are the current state of Target
	Conditions []TargetCondition `json:"conditions,omitempty"`
	// ObservedGeneration is the last generation observed by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// ObservedState is the actual target information
	ObservedState *TargetActual `json:"observedState,omitempty"`
}

// TargetConditionType is a valid value for TargetCondition.Type
type TargetConditionType string

const (
	// TargetConditionTypeReady means that API server on the Target is ready for service.
	TargetConditionReady TargetConditionType = "Ready"

	// TargetTargetFailed means that reconciling target is failed on node.
	TargetTargetFailed TargetConditionType = "TargetFailed"

	// TargetLUNFailed means that reconcilation of LUNs is failed on node.
	TargetLUNFailed TargetConditionType = "LUNFailed"
)

type TargetCondition struct {
	// Type is the type of this condition.
	Type TargetConditionType `json:"type,omitempty"`
	// Status is the status of this condition.
	Status corev1.ConditionStatus `json:"status,omitempty"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is the one-word, CamelCase reason about the last transition.
	Reason string `json:"reason,omitempty"`
	// Message is human readable message about the last transition.
	Message string `json:"message,omitempty"`
}

type TargetReference struct {
	// Name is unique to reference a target resource.
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
}

// +kubebuilder:object:root=true

// Target is the Schema for the targets API
// +kubebuilder:resource:shortName=tgt,scope=Cluster
// +kubebuilder:subresource:status
type Target struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TargetSpec   `json:"spec,omitempty"`
	Status TargetStatus `json:"status,omitempty"`
}

// SetCondition sets condition of type condType with empty reason and message.
func (target *Target) SetCondition(condType TargetConditionType, status corev1.ConditionStatus, t metav1.Time) {
	target.SetConditionReason(condType, status, "", "", t)
}

// SetConditionReason is similar to setCondition, but it takes reason and message.
func (target *Target) SetConditionReason(condType TargetConditionType, status corev1.ConditionStatus, reason, msg string, t metav1.Time) {
	cond := target.GetCondition(condType)
	if cond == nil {
		target.Status.Conditions = append(target.Status.Conditions, TargetCondition{
			Type: condType,
		})
		cond = &target.Status.Conditions[len(target.Status.Conditions)-1]
	}

	if cond.Status != status {
		cond.Status = status
		cond.LastTransitionTime = t
	}

	cond.Reason = reason
	cond.Message = msg
}

// GetCondition returns condition of type condType if it exists.  Otherwise returns nil.
func (target *Target) GetCondition(condType TargetConditionType) *TargetCondition {
	for i := range target.Status.Conditions {
		cond := &target.Status.Conditions[i]
		if cond.Type == condType {
			return cond
		}
	}
	return nil
}

func (target *Target) Ready() bool {
	r := target.GetCondition(TargetConditionReady)
	if r != nil {
		if r.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// +kubebuilder:object:root=true

// TargetList contains a list of Target
type TargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Target `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Target{}, &TargetList{})
}
