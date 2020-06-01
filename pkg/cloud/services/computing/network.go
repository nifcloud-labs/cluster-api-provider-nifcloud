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
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/aokumasan/nifcloud-sdk-go-v2/nifcloud"
	"github.com/aokumasan/nifcloud-sdk-go-v2/service/computing"
	"github.com/chyeh/pubip"
	"github.com/pkg/errors"
	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	nferrors "github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/errors"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/services/wait"
	"go.uber.org/multierr"
	"sigs.k8s.io/cluster-api/util/record"
)

const (
	IPProtocolTCP        = "TCP"
	IPProtocolUDP        = "UDP"
	maxSecurityGroupName = 15
	anyIPv4CidrBlock     = "0.0.0.0/0"
	apiEndpointPort      = 6443
	endPoint             = "ENDPOINT"
)

type SecurityGroupRole struct {
	Cluster string `json:"cluster"`
	Role    string `json:"role"`
}

func (s *Service) ReconcileNetwork() error {
	s.scope.V(2).Info("Reconciling network for cluster", "cluster-name", s.scope.Cluster.Name, "cluster-namespace", s.scope.Cluster.Namespace)

	if err := s.reconcileSecurityGroups(); err != nil {
		return err
	}

	if err := s.reconcileEndpoint(apiEndpointPort); err != nil {
		return err
	}

	s.scope.V(2).Info("Reconcile network complated successfully")
	return nil
}

func (s *Service) DeleteNetwork() error {
	s.scope.V(2).Info("Deleting network")

	if err := s.deleteSecurityGroups(); err != nil {
		return err
	}

	s.scope.V(2).Info("Delete network completed successfully")
	return nil
}

func (s *Service) DeleteEndpoint() error {
	s.scope.V(2).Info("Delete Endpoint")

	if err := s.releaseAddress(); err != nil {
		return err
	}

	s.scope.V(2).Info("Delete endpoint completed successfully")
	return nil
}

func (s *Service) reconcileSecurityGroups() error {
	s.scope.V(2).Info("Reconciling security groups")

	if s.scope.Network().SecurityGroups == nil {
		s.scope.Network().SecurityGroups = make(map[infrav1alpha2.SecurityGroupRole]infrav1alpha2.SecurityGroup)
	}

	sgs, err := s.describeSecurityGroupsByName()
	if err != nil {
		return err
	}
	// security group roles to handle with reconciliation loop
	roles := []infrav1alpha2.SecurityGroupRole{
		infrav1alpha2.SecurityGroupControlPlane,
	}
	// make sure that security groups are valid or created
	for _, role := range roles {
		sg := s.getDefaultSecurityGroup(role)
		exists, ok := sgs[*sg.GroupName]
		if !ok {
			if err := s.createSecurityGroupWithTag(role, sg); err != nil {
				return err
			}
			s.scope.SecurityGroups()[role] = infrav1alpha2.SecurityGroup{
				Name: *sg.GroupName,
			}
			s.scope.V(2).Info("Created security group for role", "role", role, "security-group", s.scope.SecurityGroups()[role])
			continue
		}

		s.scope.SecurityGroups()[role] = exists
	}

	// update security group to attouch ingress rules
	for i := range s.scope.SecurityGroups() {
		sg := s.scope.SecurityGroups()[i]
		current := sg.IngressRules
		want, err := s.getSecurityGroupIngressRules(i)
		if err != nil {
			return err
		}

		toRevoke := current.Difference(want)
		if len(toRevoke) > 0 {
			if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
				if err := s.revokeSecurityGroupIngressRules(sg.Name, toRevoke); err != nil {
					return false, err
				}
				return true, nil
			}, nferrors.SecurityGroupProcessing); err != nil {
				return errors.Wrapf(err, "failed to revoke security group ingress rules for %q", sg.Name)
			}

			s.scope.V(2).Info("revoked ingress rules from security group", "revoked-ingress-rules", toRevoke, "security-group-name", sg.Name)
		}
		toAuthorize := want.Difference(current)
		if len(toAuthorize) > 0 {
			if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
				if err := s.authorizeSecurityGroupIngressRules(sg.Name, toAuthorize); err != nil {
					return false, err
				}
				return true, nil
			}, nferrors.SecurityGroupProcessing); err != nil {
				return err
			}

			s.scope.V(2).Info("Authorized ingress rules in security group", "authorized-ingress-rules", toAuthorize, "security-group-name", sg.Name)
		}
	}

	return nil
}

func (s *Service) reconcileEndpoint(endpointPort int) error {
	s.scope.V(2).Info("Reconcile API endpoint")

	ip, err := s.getOrAllocateAddress(endPoint)
	if err != nil {
		return fmt.Errorf("failed to create IP addres: %w", err)
	}

	if len(s.scope.NifcloudCluster.Status.APIEndpoints) == 0 {
		s.scope.NifcloudCluster.Status.APIEndpoints = []infrav1alpha2.APIEndpoint{
			{
				Host: ip,
				Port: int32(endpointPort),
			},
		}
	}

	return nil
}

func (s *Service) getOrAllocateAddress(role string) (string, error) {
	out, err := s.scope.NifcloudClients.Computing.DescribeAddresses(context.TODO(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to describe address: %w", err)
	}

	for _, addr := range out.AddressesSet {
		if addr.PublicIp != nil {
			return nifcloud.StringValue(addr.PublicIp), nil
		}
	}
	return s.allocateAddress(role)
}

func (s *Service) allocateAddress(role string) (string, error) {
	out, err := s.scope.NifcloudClients.Computing.AllocateAddress(context.TODO(), &computing.AllocateAddressInput{
		Placement: &computing.RequestPlacementStruct{
			AvailabilityZone: &s.scope.NifcloudCluster.Spec.Zone,
			RegionName:       &s.scope.NifcloudCluster.Spec.Region,
		},
	})
	if err != nil {
		return "", err
	}

	return *out.PublicIp, nil
}

func (s *Service) releaseAddress() error {
	for _, e := range s.scope.NifcloudCluster.Status.APIEndpoints {
		params := &computing.ReleaseAddressInput{
			PublicIp: &e.Host,
		}
		if _, err := s.scope.NifcloudClients.Computing.ReleaseAddress(context.TODO(), params); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) deleteSecurityGroups() error {
	for _, sg := range s.scope.SecurityGroups() {
		current := sg.IngressRules
		if err := s.revokeSecurityGroupIngressRules(sg.Name, current); nferrors.IsIgnorableSecurityGroupError(err) != nil {
			return err
		}
		s.scope.V(2).Info("Revoke ingress rules from security group", "revoked-ingress-rules", current, "security-group-name", sg.Name)
	}

	var errs error
	for _, sg := range s.scope.SecurityGroups() {
		if err := s.deleteSecurityGroup(&sg); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (s *Service) deleteSecurityGroup(sg *infrav1alpha2.SecurityGroup) error {
	input := &computing.DeleteSecurityGroupInput{
		GroupName: nifcloud.String(sg.Name),
	}

	if _, err := s.scope.Computing.DeleteSecurityGroup(context.TODO(), input); nferrors.IsIgnorableSecurityGroupError(err) != nil {
		record.Warnf(s.scope.NifcloudCluster, "FailedDeleteSecurityGroup", "Failed to dlete security group %q: %v", sg.Name, err)
		s.scope.V(2).Info("Deleted security group", "security-group-name", sg.Name)
	}

	return nil
}

func (s *Service) describeSecurityGroupsByName() (map[string]infrav1alpha2.SecurityGroup, error) {
	input := &computing.DescribeSecurityGroupsInput{
		Filter: []computing.RequestFilterStruct{
			computing.RequestFilterStruct{
				Name:         nifcloud.String("group-name"),
				RequestValue: []string{s.scope.Name()},
			},
		},
	}
	out, err := s.scope.Computing.DescribeSecurityGroups(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to describe security groups")
	}

	res := make(map[string]infrav1alpha2.SecurityGroup)
	for _, sgi := range out.SecurityGroupInfo {
		sg := infrav1alpha2.SecurityGroup{
			ID:   *sgi.OwnerId,
			Name: *sgi.GroupName,
		}
		for _, rule := range sgi.IpPermissions {
			sg.IngressRules = append(sg.IngressRules, ingressRuleFromSDKType(&rule))
		}
		res[sg.Name] = sg
	}
	return res, nil
}

func (s *Service) createSecurityGroupWithTag(role infrav1alpha2.SecurityGroupRole, input *computing.SecurityGroupInfoSetItem) error {
	_, err := s.scope.NifcloudClients.Computing.CreateSecurityGroup(context.TODO(), &computing.CreateSecurityGroupInput{
		GroupName:        input.GroupName,
		GroupDescription: nifcloud.String(fmt.Sprintf("{\"cluster\":\"%s\",\"role\":\"%s\"}", s.scope.Name(), role)),
	})
	if err != nil {
		record.Warnf(s.scope.NifcloudCluster, "FailedCreateSecurityGroup", "Failed to create managed SecurityGroup for Role %q:%v", role, err)
		return errors.Wrapf(err, "failed to creat security group %q", role)
	}
	record.Eventf(s.scope.NifcloudCluster, "SuccessfulCreateGroup", "Created managed SecuirtyGroup %q for role %q", input.GroupName, role)

	return nil
}

func (s *Service) getSecurityGroupName(clusterName string, role infrav1alpha2.SecurityGroupRole) string {
	hashed := md5.Sum([]byte(role))
	tmp := fmt.Sprintf("%s%v", clusterName, hex.EncodeToString(hashed[:]))
	return tmp[:maxSecurityGroupName]
}

func (s *Service) getDefaultSecurityGroup(role infrav1alpha2.SecurityGroupRole) *computing.SecurityGroupInfoSetItem {
	name := s.getSecurityGroupName(s.scope.Name(), role)
	return &computing.SecurityGroupInfoSetItem{
		GroupName: nifcloud.String(name),
	}
}

func (s *Service) authorizeSecurityGroupIngressRules(name string, rules infrav1alpha2.IngressRules) error {
	input := &computing.AuthorizeSecurityGroupIngressInput{GroupName: nifcloud.String(name)}
	for _, rule := range rules {
		// to adopt nifcloud requirement
		sanitized := s.sanitizeRole(name, *rule)
		input.IpPermissions = append(input.IpPermissions, *ingressRuleToSDKType(sanitized))
	}
	if _, err := s.scope.Computing.AuthorizeSecurityGroupIngress(context.TODO(), input); err != nil {
		record.Warnf(s.scope.NifcloudCluster, "FailedAuthorizeSecurityGroupIngressRules", "Failed to authorize security group ingress rules %v for SecurityGroup %q: %v", rules, name, err)
		return errors.Wrapf(err, "failed to authorize security group %q ingress rules: %v", name, rules)
	}

	record.Eventf(s.scope.NifcloudCluster, "SuccessfulAuthorizeSecurityGroupIngressRules", "Authorize security group ingress rules %v for SecurityGrup %q", rules, name)
	return nil
}

func (s *Service) revokeSecurityGroupIngressRules(name string, rules infrav1alpha2.IngressRules) error {
	input := &computing.RevokeSecurityGroupIngressInput{GroupName: nifcloud.String(name)}
	for _, rule := range rules {
		input.IpPermissions = append(input.IpPermissions, *ingressRuleToSDKType(rule))
	}
	if _, err := s.scope.Computing.RevokeSecurityGroupIngress(context.TODO(), input); err != nil {
		record.Warnf(s.scope.NifcloudCluster, "FailedRevokeSecurityGroupIngressRules", "Failed to revoke security group ingress rules %v for SecurityGroup %q: %v", rules, name, err)
		return errors.Wrapf(err, "failed to revoke security group %q ingress rules: %v", name, rules)
	}

	record.Eventf(s.scope.NifcloudCluster, "SuccessfulRevokeSecurityGroupIngressRules", "Revoke security group ingress rules %v for SecurityGrup %q", rules, name)
	return nil
}

func (s *Service) sanitizeRole(name string, i infrav1alpha2.IngressRule) *infrav1alpha2.IngressRule {
	var groups []string
	for _, group := range i.SourceSecurityGroupName {
		if group != name {
			groups = append(groups, group)
		}
	}
	i.SourceSecurityGroupName = groups
	return &i
}

func ingressRuleToSDKType(i *infrav1alpha2.IngressRule) (res *computing.RequestIpPermissionsStruct) {
	switch i.Protocol {
	case infrav1alpha2.SecurityGroupProtocolTCP, infrav1alpha2.SecurityGroupProtocolUDP:
		res = &computing.RequestIpPermissionsStruct{
			IpProtocol: nifcloud.String(string(i.Protocol)),
			FromPort:   nifcloud.Int64(i.FromPort),
			ToPort:     nifcloud.Int64(i.ToPort),
		}
	default:
		res = &computing.RequestIpPermissionsStruct{
			IpProtocol: nifcloud.String(string(i.Protocol)),
		}
	}

	res.Description = nifcloud.String(i.Description)

	for _, cidr := range i.CidrBlocks {
		ipRange := &computing.RequestIpRangesStruct{
			CidrIp: nifcloud.String(cidr),
		}
		res.RequestIpRanges = append(res.RequestIpRanges, *ipRange)
	}
	for _, group := range i.SourceSecurityGroupName {
		groupNames := &computing.RequestGroupsStruct{
			GroupName: nifcloud.String(group),
		}
		res.RequestGroups = append(res.RequestGroups, *groupNames)
	}

	return res
}

func ingressRuleFromSDKType(v *computing.IpPermissionsSetItem) (res *infrav1alpha2.IngressRule) {
	switch *v.IpProtocol {
	case IPProtocolTCP, IPProtocolUDP:
		res = &infrav1alpha2.IngressRule{
			Protocol: infrav1alpha2.SecurityGroupProtocol(*v.IpProtocol),
			FromPort: *v.FromPort,
			ToPort:   *v.ToPort,
		}
	default:
		res = &infrav1alpha2.IngressRule{
			Protocol: infrav1alpha2.SecurityGroupProtocol(*v.IpProtocol),
		}
	}

	if v.Description != nil && *v.Description != "" {
		res.Description = *v.Description
	}
	for _, ranges := range v.IpRanges {
		res.CidrBlocks = append(res.CidrBlocks, *ranges.CidrIp)
	}
	for _, pair := range v.Groups {
		if pair.GroupName == nil {
			continue
		}
		res.SourceSecurityGroupName = append(res.SourceSecurityGroupName, *pair.GroupName)
	}
	return res
}

func (s *Service) defaultSSHAsIP(ip string) *infrav1alpha2.IngressRule {
	return &infrav1alpha2.IngressRule{
		Description: "SSH",
		Protocol:    infrav1alpha2.SecurityGroupProtocolTCP,
		FromPort:    22,
		ToPort:      22,
		CidrBlocks:  []string{ip},
	}
}

func (s *Service) getSecurityGroupIngressRules(role infrav1alpha2.SecurityGroupRole) (infrav1alpha2.IngressRules, error) {
	hostIP, err := pubip.Get()
	if err != nil {
		return nil, err
	}

	switch role {
	// TODO: divide groups to each role
	case infrav1alpha2.SecurityGroupControlPlane, infrav1alpha2.SecurityGroupNode:
		return infrav1alpha2.IngressRules{
			s.defaultSSHAsIP(hostIP.String()),
			{
				Description: "Kubernetes API",
				Protocol:    infrav1alpha2.SecurityGroupProtocolTCP,
				FromPort:    6443,
				ToPort:      6443,
				CidrBlocks:  []string{hostIP.String()},
			},
		}, nil
	}

	return nil, errors.Errorf("Cannot determine ingress rules for unknown security group role %q", role)
}
