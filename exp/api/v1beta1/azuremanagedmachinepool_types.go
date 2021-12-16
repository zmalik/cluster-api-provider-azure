/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capierrors "sigs.k8s.io/cluster-api/errors"
)

const (
	// LabelAgentPoolMode represents mode of an agent pool. Possible values include: System, User.
	LabelAgentPoolMode = "azuremanagedmachinepool.infrastructure.cluster.x-k8s.io/agentpoolmode"

	// NodePoolModeSystem represents mode system for azuremachinepool.
	NodePoolModeSystem NodePoolMode = "System"

	// NodePoolModeUser represents mode user for azuremachinepool.
	NodePoolModeUser NodePoolMode = "User"
)

// NodePoolMode enumerates the values for agent pool mode.
type NodePoolMode string

// AzureManagedMachinePoolSpec defines the desired state of AzureManagedMachinePool.
type AzureManagedMachinePoolSpec struct {

	// Name - name of the agent pool. If not specified, CAPZ uses the name of the CR as the agent pool name.
	// +optional
	Name *string `json:"name,omitempty"`

	// Mode - represents mode of an agent pool. Possible values include: System, User.
	// +kubebuilder:validation:Enum=System;User
	Mode string `json:"mode"`

	// SKU is the size of the VMs in the node pool.
	SKU string `json:"sku"`

	// OSDiskSizeGB is the disk size for every machine in this agent pool.
	// If you specify 0, it will apply the default osDisk size according to the vmSize specified.
	// +optional
	OSDiskSizeGB *int32 `json:"osDiskSizeGB,omitempty"`

	// AvailabilityZones - Availability zones for nodes. Must use VirtualMachineScaleSets AgentPoolType.
	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// Taints specifies the taints for nodes present in this agent pool.
	// +optional
	Taints Taints `json:"taints,omitempty"`

	// ProviderIDList is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`

	// Scaling specifies the autoscaling parameters for the node pool.
	// +optional
	Scaling *ManagedMachinePoolScaling `json:"scaling,omitempty"`

	// MaxPods specifies the kubelet --max-pods configuration for the node pool.
	// +optional
	MaxPods *int32 `json:"maxPods,omitempty"`
}

// ManagedMachinePoolScaling specifies scaling options.
type ManagedMachinePoolScaling struct {
	MinSize *int32 `json:"minSize,omitempty"`
	MaxSize *int32 `json:"maxSize,omitempty"`
}

// TaintEffect is the effect for a Kubernetes taint.
type TaintEffect string

var (
	// TaintEffectNoSchedule  is a taint that does not allow new pods to schedule onto the node unless
	// they tolerate the taint, but allow all pods submitted to Kubelet without going through the scheduler
	// to start, and allow all already-running pods to continue running.
	// Enforced by the scheduler.
	TaintEffectNoSchedule = TaintEffect("no-schedule")
	// TaintEffectNoExecute will evict any already-running pods that do not tolerate the taint.
	// Currently enforced by NodeController.
	TaintEffectNoExecute = TaintEffect("no-execute")
	// TaintEffectPreferNoSchedule is Like TaintEffectNoSchedule, but the scheduler tries not to schedule
	// new pods onto the node, rather than prohibiting new pods from scheduling
	// onto the node entirely. Enforced by the scheduler.
	TaintEffectPreferNoSchedule = TaintEffect("prefer-no-schedule")
)

type Taint struct {
	// Effect specifies the effect for the taint
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=no-schedule;no-execute;prefer-no-schedule
	Effect TaintEffect `json:"effect"`
	// Key is the key of the taint
	// +kubebuilder:validation:Required
	Key string `json:"key"`
	// Value is the value of the taint
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// Taints is an array of Taints.
type Taints []Taint

// AzureManagedMachinePoolStatus defines the observed state of AzureManagedMachinePool.
type AzureManagedMachinePoolStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Replicas is the most recently observed number of replicas.
	// +optional
	Replicas int32 `json:"replicas"`

	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	ErrorReason *capierrors.MachineStatusError `json:"errorReason,omitempty"`

	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=azuremanagedmachinepools,scope=Namespaced,categories=cluster-api,shortName=ammp
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// AzureManagedMachinePool is the Schema for the azuremanagedmachinepools API.
type AzureManagedMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureManagedMachinePoolSpec   `json:"spec,omitempty"`
	Status AzureManagedMachinePoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AzureManagedMachinePoolList contains a list of AzureManagedMachinePools.
type AzureManagedMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureManagedMachinePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureManagedMachinePool{}, &AzureManagedMachinePoolList{})
}
