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
	"fmt"
	"sort"

	corev1 "k8s.io/api/core/v1"
)

type APIEndpoint struct {
	// the hostname on which the API server is serving
	Host string `json:"host"`

	// the port on which the API server is serving
	Port int32 `json:"port"`
}

type Network struct {
	// SecurityGroups is a map from a name of role/kind to spesific role filewall
	SecurityGroups map[SecurityGroupRole]SecurityGroup `json:"securityGroups,omitempty"`
}

type NetworkSpec struct {
}

// InstanceState describes the state of an nifcloud instance.
type InstanceState string

var (
	InstancePending = InstanceState("pending")
	InstanceRunning = InstanceState("running")
	InstanceStopped = InstanceState("stopped")
	InstanceWaiting = InstanceState("waiting")
)

// SecurityGroupRole defines the unique role of a security group.
type SecurityGroupRole string

var (
	// SSH entry point role
	SecurityGroupBastion = SecurityGroupRole("bastion")
	// kubernets controleplane node role
	SecurityGroupControlPlane = SecurityGroupRole("controlplane")
	// kubernetes workload node role
	SecurityGroupNode = SecurityGroupRole("node")
)

// SecurityGroup defines nifcloud firewall group
type SecurityGroup struct {
	// ID is an identifier
	ID string `json:"id"`
	// security(firewall) group name
	Name string `json:"name"`
	// ingress rules of the group
	// +optional
	IngressRules IngressRules `json:"ingressRules"`
}

func (s *SecurityGroup) String() string {
	return fmt.Sprintf("id=%s/name=%s", s.ID, s.Name)
}

// SecurityGroupProtocol defines the protocol type for a security group rule.
type SecurityGroupProtocol string

var (
	SecurityGroupProtocolAny = SecurityGroupProtocol("ANY")
	SecurityGroupProtocolTCP = SecurityGroupProtocol("TCP")
	SecurityGroupProtocolUDP = SecurityGroupProtocol("UDP")
)

type IngressRule struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description,omitempty"`
	Protocol    SecurityGroupProtocol `json:"protocol"`
	FromPort    int64                 `json:"fromPort"`
	ToPort      int64                 `json:"toPort"`

	// List of CIDR blocks to allow access from. Cannot be specified with SourceSecurityGroupID.
	// +optional
	CidrBlocks []string `json:"cidrBlocks,omitempty"`
	// The security group id to allow access from. Cannot be specified with CidrBlocks.
	// +optional
	SourceSecurityGroupName []string `json:"sourceSecurityGroupName,omitempty"`
}

func (i IngressRule) String() string {
	return fmt.Sprintf("protocol=%s/range[%d-%d]/description=%s", i.Protocol, i.FromPort, i.ToPort, i.Description)
}

type IngressRules []*IngressRule

func (i IngressRules) Difference(o IngressRules) (out IngressRules) {
	for _, x := range i {
		found := false
		for _, y := range o {
			if x.Equals(y) {
				found = true
				break
			}
		}
		if !found {
			out = append(out, x)
		}
	}
	return
}

func (i *IngressRule) Equals(o *IngressRule) bool {
	if len(i.CidrBlocks) != len(i.CidrBlocks) {
		return false
	}
	if len(i.SourceSecurityGroupName) != len(o.SourceSecurityGroupName) {
		return false
	}

	if len(i.CidrBlocks) != 0 {
		sort.Strings(i.CidrBlocks)
		sort.Strings(o.CidrBlocks)
		for ii, v := range i.CidrBlocks {
			if v != o.CidrBlocks[ii] {
				return false
			}
		}
	}
	sort.Strings(i.SourceSecurityGroupName)
	sort.Strings(o.SourceSecurityGroupName)
	for ii, v := range i.SourceSecurityGroupName {
		if v != o.SourceSecurityGroupName[ii] {
			return false
		}
	}

	if i.Description != o.Description || i.Protocol != o.Protocol {
		return false
	}

	switch i.Protocol {
	case SecurityGroupProtocolTCP, SecurityGroupProtocolUDP:
		return i.FromPort == o.FromPort && i.ToPort == o.ToPort
	}

	return true
}

// DisableApiTermination should be false to delete server with NifcloudAPI
const ApiTermination = false

type Instance struct {
	// UID is an instance identifier
	UID string `json:"uid"`
	// ID is a name of nifcloud instance
	ID string `json:"id,omitemptuy"`
	// Zone is machine location
	Zone string `json:"zone,omitempty"`
	// State is current state of nicloud instance
	State InstanceState `json:"state,omitempty"`
	// Type is machine type of nicloud instance
	Type string `json:"type,omitempty"`
	// ImageID is an image running on nicloud instance
	ImageID string `json:"imageID,omitempty"`
	// UserData is cloud-init script
	UserData *string `json:"userData,omitempty"`
	// security group names
	SecurityGroups []string `json:"securityGroups,omitempty"`
	// A name of SSh key pair
	SSHKeyName string `json:"sshKeyName,omitempty"`
	// tags in instance
	Tag Tag `json:"tag,omitempty"`
	// The public IPv4 address assigned to the instance
	PublicIP string `json:"publicIP,omitempty"`
	// The private IPv4 address assigned to the instance
	PrivateIP string `json:"privateIP,omitempty"`
	// Address containes a list of apiserver endpoints
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`
	// a list of networkinterface which attached the instance
	NetworkInterfaces []string `json:"networkInterfaces,omitempty"`
}
