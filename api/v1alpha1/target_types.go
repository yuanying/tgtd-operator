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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TargetSpec defines the desired state of Target
type TargetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// TargetNodeName is a node name where the target will be placed.
	TargetNodeName string `json:"targetNodeName,omitempty"`

	// IQN is an iqn of the target
	// +kubebuilder:validation:Required
	IQN string `json:"iqn,omitempty"`

	// InitiatorAddresses is a list of initiator address. If "All" is specified, any addresses are allowed.
	InitiatorAddresses []string `json:"initiatorAddresses,omitempty"`

	// LUNs is a list of LUNs
	LUNs []TargetLUN `json:"targetLUNs,omitempty"`
}

// TargetLun is the specification of LUN
type TargetLUN struct {
	// LID is an id of the LUN
	// +kubebuilder:validation:Required
	LID int32 `json:"lid,omitempty"`

	// BackingStore is a path of the backing store
	// +kubebuilder:validation:Required
	BackingStore string `json:"backingStore,omitempty"`

	// BSType is a backing store type
	BSType *string `json:"bsType,omitempty"`

	// BSOpts is a options for backing store
	BSOpts *string `json:"bsOpts,omitempty"`
}

// TargetStatus defines the observed state of Target
type TargetStatus struct {
	// Conditions are the current state of Target
	Conditions []TargetCondition `json:"conditions,omitempty"`
	// ObservedGeneration is the last generation observed by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// TargetConditionType is a valid value for TargetCondition.Type
type TargetConditionType string

const (
	// TargetConditionTypeReady means that API server on the Target is ready for service.
	TargetConditionReady TargetConditionType = "Ready"
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
