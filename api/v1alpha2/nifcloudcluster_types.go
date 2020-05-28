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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Finalizer is a label to allow ReconcileCluster to clean up
	// all resouces associated with NifcloudClustr before removing
	// from api server
	ClusterFinalizer = "nifcloudcluster.infrastructure.cluster.x-k8s.io"
)

// NifcloudClusterSpec defines the desired state of NifcloudCluster
type NifcloudClusterSpec struct {
	// NetworkSpec includes nifcloud network configurations
	NetworkSpec NetworkSpec `json:"networkSpec,omitempty"`

	// Zone is a nifcloud zone which cluster lives on
	Zone string `json:"zone,omitempty"`

	// Region ins a nifcloud region
	Region string `json:"region,omitempty"`

	// SSHKeyName is the name of ssh key to attach to the bastion
	SSHKeyName string `json:"sshKeyName,omitempty"`
}

// NifcloudClusterStatus defines the observed state of NifcloudCluster
type NifcloudClusterStatus struct {
	// cluster network configurations
	Network Network `json:"network,omitempty"`

	// bastion instatnce information
	Bastion *Instance `json:"bastion,omitempty"`

	// cluster resource is ready to available or not
	Ready bool `json:"ready,omitempty"`

	APIEndpoints []APIEndpoint `json:"apiEndpoints,omitempty"`

	// +optional
	ErrorReason string `json:"failureReason,omitempty"`
	// +optional
	ErrorMessage string `json:"failureMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=nifcloudclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status

// NifcloudCluster is the Schema for the nifcloudclusters API
type NifcloudCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NifcloudClusterSpec   `json:"spec,omitempty"`
	Status NifcloudClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NifcloudClusterList contains a list of NifcloudCluster
type NifcloudClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NifcloudCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NifcloudCluster{}, &NifcloudClusterList{})
}
