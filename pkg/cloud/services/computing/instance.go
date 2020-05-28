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

package computing

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/cluster-api/util/record"

	corev1 "k8s.io/api/core/v1"

	"github.com/aokumasan/nifcloud-sdk-go-v2/nifcloud"
	"github.com/aokumasan/nifcloud-sdk-go-v2/service/computing"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	nferrors "github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/errors"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/scope"
)

const (
	defaultSSHKeyName     = "default"
	defaultMachineOwnerID = "niftycloud"
	defaultMachineBaseOS  = ""
)

func (s *Service) InstanceIfExists(id *string) (*infrav1alpha2.Instance, error) {
	if id == nil {
		s.scope.Info("Instance does not have an instance id")
		return nil, nil
	}

	s.scope.V(2).Info("Looking for instance by id", "instance-id", *id)

	input := &computing.DescribeInstancesInput{
		InstanceId: []string{nifcloud.StringValue(id)},
	}
	out, err := s.scope.NifcloudClients.Computing.DescribeInstances(context.TODO(), input)
	switch {
	case nferrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed to describe instance[%q]: %w", *id, err)
	}

	for _, rs := range out.ReservationSet {
		if len(rs.InstancesSet) == 0 {
			break
		}
		instance := rs.InstancesSet[0]
		if nifcloud.StringValue(id) != nifcloud.StringValue(instance.InstanceId) {
			continue
		}
		return s.SDKToInstance(out.ReservationSet[0].InstancesSet[0])
	}
	return nil, nil
}

func (s *Service) GetRunningInstanceByTag(scope *scope.MachineScope) (*infrav1alpha2.Instance, error) {
	s.scope.V(2).Info("Looking for existing machine instance by tags")

	input := &computing.DescribeInstancesInput{}
	out, err := s.scope.NifcloudClients.Computing.DescribeInstances(context.TODO(), input)
	switch {
	case nferrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed to describe running instance: %w", err)
	}

	filtered := s.FilterInstancesByTag(out.ReservationSet, map[string]string{
		"cluster": s.scope.Name(),
		"role":    scope.Role(),
	})
	for _, res := range filtered {
		for _, instance := range res.InstancesSet {
			// filter by name
			if nifcloud.StringValue(instance.InstanceId) != nifcloud.StringValue(scope.GetInstanceID()) {
				continue
			}
			return s.SDKToInstance(instance)
		}
	}

	return nil, nil
}

func (s *Service) FilterInstancesByTag(vs []computing.ReservationSetItem, tags map[string]string) []computing.ReservationSetItem {
	var filtered []computing.ReservationSetItem
	for _, v := range vs {
		iTag := v1alpha2.ParseTags(nifcloud.StringValue(v.Description))
		ok := true
		for key, val := range iTag {
			if tags[key] != val {
				ok = false
				break
			}
		}
		if ok {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func (s *Service) CreateInstance(scope *scope.MachineScope) (*infrav1alpha2.Instance, error) {
	s.scope.V(2).Info("Creating an instance for a machine")

	instanceID := scope.GetInstanceIDConved()
	input := &infrav1alpha2.Instance{
		ID:                instanceID,
		Type:              scope.NifcloudMachine.Spec.InstanceType,
		NetworkInterfaces: scope.NifcloudMachine.Spec.NetworkInterfaces,
	}

	// create tags
	input.Tag = v1alpha2.BuildTags(v1alpha2.BuildParams{
		ClusterName: s.scope.Name(),
		Role:        nifcloud.String(scope.Role()),
	})
	var err error
	// set image from the machine configuration
	if scope.NifcloudMachine.Spec.ImageID != "" {
		input.ImageID = scope.NifcloudMachine.Spec.ImageID
	} else {
		input.ImageID, err = s.defaultImageLookup()
		if err != nil {
			return nil, err
		}
	}

	// set userdata
	userData, err := scope.GetUserData()
	if err != nil {
		scope.Info("failed to get bootstrap data")
		return nil, err
	}
	input.UserData = pointer.StringPtr(userData)

	ids, err := s.GetCoreSecurityGroup(scope)
	if err != nil {
		return nil, err
	}
	input.SecurityGroups = append(input.SecurityGroups, ids...)

	// set SSH key
	input.SSHKeyName = defaultSSHKeyName
	if scope.NifcloudMachine.Spec.KeyName != "" {
		input.SSHKeyName = scope.NifcloudMachine.Spec.KeyName
	} else {
		input.SSHKeyName = defaultSSHKeyName
	}

	s.scope.V(2).Info("Running instance", "machine-role", scope.Role())
	out, err := s.runInstance(scope.Role(), input)
	if err != nil {
		record.Warnf(scope.NifcloudMachine, "FailedCreate", "Failed to create instance: %v", err)
		return nil, err
	}

	if len(input.NetworkInterfaces) > 0 {
		for _, id := range input.NetworkInterfaces {
			// TODO: attach interface
			s.scope.V(2).Info("Attaching security groups to provide network interface", "groups", "[TODO]", "interface", id)
		}
	}

	// set reserved ip addr to controlplane
	labels := scope.Machine.GetLabels()
	if labels["cluster.x-k8s.io/control-plane"] == "true" {
		if err := s.attachAddress(out.ID, scope.NifcloudCluster.Status.APIEndpoints[0].Host); err != nil {
			return out, err
		}
	}

	record.Eventf(scope.NifcloudMachine, "SuccessfulCreate", "Created new instance [%s/%s]", scope.Role(), out.ID)
	return out, nil
}

func (s *Service) GetCoreSecurityGroup(scope *scope.MachineScope) ([]string, error) {
	sgRoles := []infrav1alpha2.SecurityGroupRole{}

	switch scope.Role() {
	case "node", "control-plane":
		sgRoles = append(sgRoles, infrav1alpha2.SecurityGroupControlPlane)
	default:
		return nil, errors.Errorf("Unknown node role %q", scope.Role())
	}
	ids := make([]string, 0, len(sgRoles))
	for _, sg := range sgRoles {
		if _, ok := s.scope.SecurityGroups()[sg]; !ok {
			return nil, nferrors.NewFailedDependency(
				errors.Errorf("%s security group not available", sg),
			)
		}
		ids = append(ids, s.scope.SecurityGroups()[sg].Name)
	}
	return ids, nil
}

func (s *Service) StopAndTerminateInstanceWithTimeout(instanceID string) error {
	ctx := context.TODO()
	input := &computing.DescribeInstancesInput{
		InstanceId: []string{instanceID},
	}

	// stopping server before terminating
	if err := s.StopInstance(instanceID); err != nil {
		return err
	}
	s.scope.V(2).Info("Waiting for Nifcloud server to stop", "instance-id", instanceID)

	if err := s.scope.NifcloudClients.Computing.WaitUntilInstanceStopped(ctx, input); err != nil {
		return fmt.Errorf("failed to wait for instance %q stopping: %w", instanceID, err)
	}

	if err := s.TerminateInstance(instanceID); err != nil {
		return err
	}
	s.scope.V(2).Info("Waiting for Nifcloud server to terminate", "intance-id", instanceID)

	if err := s.scope.NifcloudClients.Computing.WaitUntilInstanceDeleted(ctx, input); err != nil {
		return fmt.Errorf("failed to wait for instance %q termination: %w", instanceID, err)
	}

	return nil
}

func (s *Service) TerminateInstance(instanceID string) error {
	s.scope.V(2).Info("Try to terminate instance", "instance-id", instanceID)

	input := &computing.TerminateInstancesInput{
		InstanceId: []string{instanceID},
	}
	if _, err := s.scope.NifcloudClients.Computing.TerminateInstances(context.TODO(), input); err != nil {
		return fmt.Errorf("failed to termiante instance with id %q: %w", instanceID, err)
	}

	s.scope.V(2).Info("Terminated instance", "instance-id", instanceID)
	return nil
}

func (s *Service) StopInstance(instanceID string) error {
	s.scope.V(2).Info("Try to stop instance", "instance-id", instanceID)

	input := &computing.StopInstancesInput{
		InstanceId: []string{instanceID},
	}
	if _, err := s.scope.NifcloudClients.Computing.StopInstances(context.TODO(), input); err != nil {
		return fmt.Errorf("failed to stop instance with id %q: %w", instanceID, err)
	}

	s.scope.V(2).Info("Stoped instance", "instance-id", instanceID)
	return nil
}

func (s *Service) runInstance(role string, i *infrav1alpha2.Instance) (*infrav1alpha2.Instance, error) {
	apiTermination := infrav1alpha2.ApiTermination
	input := &computing.RunInstancesInput{
		InstanceId:            &i.ID,
		InstanceType:          &i.Type,
		ImageId:               &i.ImageID,
		KeyName:               &i.SSHKeyName,
		DisableApiTermination: &apiTermination,
	}
	if i.UserData != nil {
		input.UserData = i.UserData
		s.scope.Info("userData size", "bytes", len(nifcloud.StringValue(input.UserData)), "role", role)
	}

	if len(i.NetworkInterfaces) > 0 {
		netInterfaces := make([]computing.RequestNetworkInterfaceStruct, 0, len(i.NetworkInterfaces))
		for index, id := range i.NetworkInterfaces {
			idx := int64(index)
			netInterfaces = append(netInterfaces, computing.RequestNetworkInterfaceStruct{
				DeviceIndex: &idx,
				NetworkId:   &id,
			})
		}
		input.NetworkInterface = netInterfaces
	} else {
		if len(i.SecurityGroups) > 0 {
			input.SecurityGroup = i.SecurityGroups
		}
	}

	// tag to instance Description
	input.Description = i.Tag.ConvToString()

	ctx := context.TODO()
	creating, err := s.scope.NifcloudClients.Computing.RunInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to run instance: %w", err)
	}

	instanceID := *creating.InstancesSet[0].InstanceId
	s.scope.V(2).Info("Waiting for instance to be in running state", "instance-id", instanceID)

	if len(creating.InstancesSet) == 0 {
		return nil, fmt.Errorf("no instance returned for reservation: %v", creating.String())
	}

	describeInput := &computing.DescribeInstancesInput{
		InstanceId: []string{instanceID},
	}
	if err := s.scope.NifcloudClients.Computing.WaitUntilInstanceRunning(ctx, describeInput); err != nil {
		return nil, fmt.Errorf("failed to wait for instance %q running: %w", instanceID, err)
	}

	running, err := s.scope.NifcloudClients.Computing.DescribeInstances(ctx, describeInput)
	switch {
	case nferrors.IsNotFound(err):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed to describe instance[%q]: %w", instanceID, err)
	}

	return s.SDKToInstance(running.ReservationSet[0].InstancesSet[0])
}

func (s *Service) SDKToInstance(v computing.InstancesSetItem) (*infrav1alpha2.Instance, error) {
	i := &infrav1alpha2.Instance{
		UID:        *v.InstanceUniqueId,
		ID:         *v.InstanceId,
		Zone:       *v.Placement.AvailabilityZone,
		State:      infrav1alpha2.InstanceState(*v.InstanceState.Name),
		Type:       *v.InstanceType,
		ImageID:    *v.ImageId,
		SSHKeyName: *v.KeyName,
		PublicIP:   *v.IpAddress,
		PrivateIP:  *v.PrivateIpAddress,
	}

	i.Tag = v1alpha2.ParseTags(*v.Description)

	// TODO: security groups

	i.Addresses = s.getInstanceAddresses(&v)

	return i, nil
}

func (s *Service) getInstanceAddresses(instance *computing.InstancesSetItem) []corev1.NodeAddress {
	addresses := []corev1.NodeAddress{}
	for _, ni := range instance.NetworkInterfaceSet {
		privateDNSAddress := corev1.NodeAddress{
			Type:    corev1.NodeInternalDNS,
			Address: *ni.PrivateDnsName,
		}
		privateIPAddress := corev1.NodeAddress{
			Type: corev1.NodeInternalIP,
		}
		addresses = append(addresses, privateDNSAddress, privateIPAddress)

		if ni.Association != nil {
			publicDNSAddress := corev1.NodeAddress{
				Type:    corev1.NodeExternalDNS,
				Address: *ni.Association.PublicDnsName,
			}
			publicIPAddress := corev1.NodeAddress{
				Type:    corev1.NodeExternalIP,
				Address: *ni.Association.PublicIp,
			}
			addresses = append(addresses, publicDNSAddress, publicIPAddress)
		}
	}
	return addresses
}

func (s *Service) attachAddress(instanceID, ip string) error {
	ctx := context.TODO()
	_, err := s.scope.NifcloudClients.Computing.AssociateAddress(ctx, &computing.AssociateAddressInput{
		PublicIp:   nifcloud.String(ip),
		InstanceId: nifcloud.String(instanceID),
	})
	if err != nil {
		return err
	}

	describeInput := &computing.DescribeInstancesInput{
		InstanceId: []string{instanceID},
	}
	if err := s.scope.NifcloudClients.Computing.WaitUntilInstanceRunning(ctx, describeInput); err != nil {
		return errors.Wrapf(err, "failed to wait for instance %q running", instanceID)
	}

	return nil
}

func (s *Service) defaultImageLookup() (string, error) {
	baseOS := defaultMachineBaseOS
	ownerID := defaultMachineOwnerID
	input := &computing.DescribeImagesInput{
		ImageName: []string{baseOS},
		Owner:     []string{ownerID},
	}
	out, err := s.scope.NifcloudClients.Computing.DescribeImages(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("fail to find image[%q]: %w", baseOS, err)
	}
	if len(out.ImagesSet) == 0 {
		return "", fmt.Errorf("no Images found: %q", baseOS)
	}
	s.scope.V(2).Info("Found and using an existing Image", "image-id", out.ImagesSet[0].ImageId)
	return *out.ImagesSet[0].ImageId, nil
}
