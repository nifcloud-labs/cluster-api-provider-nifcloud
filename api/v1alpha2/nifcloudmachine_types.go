/*
Copyright 2020 FUJITSU CLOUD TECHNOLOGIES LIMITED. All Rights Reserved.

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

package v1alpha2

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/errors"
)

const (
	// MachineFinalizer allow reconciler to clean up Nifcloud Machine resources
	// before removing these instances
	MachineFinalizer = "nifcloudmachine.infrastructure.cluster.x-k8s.io"
)

// NifcloudMachineSpec defines the desired state of NifcloudMachine
type NifcloudMachineSpec struct {
	// the identifier for the provider's machine instance
	ProviderID *string `json:"providerID,omitempty"`

	// InstanceID is corresponding to nifcloud `instance id`
	InstanceID string `json:"instanceID,omitempty"`

	// ImageID is instance os image
	ImageID string `json:"imageID,omitempty"`

	// AvailabilityZone is reference to nifcloud availability zone for this instance
	AvailabilityZone *string `json:"availabilityZone,omitempty"`

	// KeyName is a ssh key name to attach to this instance
	KeyName string `json:"keyName,omitempty"`

	// InstanceType is reference to nifcloud instance type
	InstanceType string `json:"instanceType,omitempty"`

	// PublicType specifies whether this machine get public IP address or not
	// +optional
	PublicType string `json:"publicType,omitempty"`

	// NetworkInterfaces is a list of nifcloud networkInterfaceSet
	// max 2 entry : public,private
	// +optional
	// +kubebuilder:validation:MaxItems=2
	NetworkInterfaces []string `json:"networkInterfaces,omitempty"`
}

// NifcloudMachineStatus defines the observed state of NifcloudMachine
type NifcloudMachineStatus struct {
	// Ready is a flag whether this resouce is available or not
	Ready bool `json:"ready"`

	// Address contains apiserver endpoints
	Address []v1.NodeAddress `json:"address,omitempty"`

	// InstanceState is the state of the nifcloud instance
	InstanceState *InstanceState `json:"instanceState,omitempty"`

	// Bootstrap data has been sended to server
	SendBootstrap bool `json:"sendBootstrap,omitempty"`

	// +optional
	ErrorReason *errors.MachineStatusError `json:"errorReason,omitempty"`
	// +optional
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resouce:path=nifcloudmachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresouce:status

// NifcloudMachine is the Schema for the nifcloudmachines API
type NifcloudMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NifcloudMachineSpec   `json:"spec,omitempty"`
	Status NifcloudMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NifcloudMachineList contains a list of NifcloudMachine
type NifcloudMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NifcloudMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NifcloudMachine{}, &NifcloudMachineList{})
}
