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

package nifcloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"

	nc "github.com/aokumasan/nifcloud-sdk-go-v2/nifcloud"
	"github.com/aokumasan/nifcloud-sdk-go-v2/service/computing"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud"
)

type nifcloudClientFactory struct {
}

func (nf *nifcloudClientFactory) CreateClient(accessKey, secretKey, region string) (cloud.Client, error) {
	client, err := New(accessKey, secretKey, region)
	if err != nil {
		return nil, err
	}
	return client, err
}

func NewNifcloudClientFactory() cloud.ClientFactory {
	return &nifcloudClientFactory{}
}

type nifcloud struct {
	client *computing.Client
}

func New(accessKey, secretKey, region string) (cloud.Client, error) {
	return &nifcloud{
		client: getNifcloudComputingClient(accessKey, secretKey, region),
	}, nil
}

func getNifcloudComputingClient(accessKey, secretKey, region string) *computing.Client {
	cfg := nc.NewConfig(
		accessKey,
		secretKey,
		region,
	)
	return computing.New(cfg)
}

func (nc *nifcloud) AllocateAddress(ctx context.Context, input *computing.AllocateAddressInput) (*computing.AllocateAddressOutput, error) {
	request := nc.client.AllocateAddressRequest(input)
	res, err := request.Send(ctx)

	if err != nil {
		return nil, err
	}

	if res.PublicIp == nil {
		return nil, errors.New("fail to allocate public address")
	}
	return res.AllocateAddressOutput, nil
}

func (nc *nifcloud) ReleaseAddress(ctx context.Context, input *computing.ReleaseAddressInput) (*computing.ReleaseAddressOutput, error) {
	request := nc.client.ReleaseAddressRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}

	if *res.Return == false {
		return nil, errors.New("fail to release address")
	}
	return res.ReleaseAddressOutput, nil
}

func (nc *nifcloud) DescribeAddresses(ctx context.Context, input *computing.DescribeAddressesInput) (*computing.DescribeAddressesOutput, error) {
	request := nc.client.DescribeAddressesRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.DescribeAddressesOutput, nil
}

func (nc *nifcloud) DisassociateAddress(ctx context.Context, input *computing.DisassociateAddressInput) (*computing.DisassociateAddressOutput, error) {
	request := nc.client.DisassociateAddressRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.DisassociateAddressOutput, nil
}

func (nc *nifcloud) DescribeInstances(ctx context.Context, input *computing.DescribeInstancesInput) (*computing.DescribeInstancesOutput, error) {
	request := nc.client.DescribeInstancesRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.DescribeInstancesOutput, nil
}

func (nc *nifcloud) DescribeImages(ctx context.Context, input *computing.DescribeImagesInput) (*computing.DescribeImagesOutput, error) {
	request := nc.client.DescribeImagesRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.DescribeImagesOutput, nil
}

func (nc *nifcloud) RunInstances(ctx context.Context, input *computing.RunInstancesInput) (*computing.RunInstancesOutput, error) {
	request := nc.client.RunInstancesRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.RunInstancesOutput, nil
}

func (nc *nifcloud) StopInstances(ctx context.Context, input *computing.StopInstancesInput) (*computing.StopInstancesOutput, error) {
	request := nc.client.StopInstancesRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.StopInstancesOutput, nil
}

func (nc *nifcloud) TerminateInstances(ctx context.Context, input *computing.TerminateInstancesInput) (*computing.TerminateInstancesOutput, error) {
	request := nc.client.TerminateInstancesRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.TerminateInstancesOutput, nil
}

func (nc *nifcloud) CreateSecurityGroup(ctx context.Context, input *computing.CreateSecurityGroupInput) (*computing.CreateSecurityGroupOutput, error) {
	request := nc.client.CreateSecurityGroupRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.CreateSecurityGroupOutput, nil
}

func (nc *nifcloud) DeleteSecurityGroup(ctx context.Context, input *computing.DeleteSecurityGroupInput) (*computing.DeleteSecurityGroupOutput, error) {
	request := nc.client.DeleteSecurityGroupRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.DeleteSecurityGroupOutput, nil
}

func (nc *nifcloud) DescribeSecurityGroups(ctx context.Context, input *computing.DescribeSecurityGroupsInput) (*computing.DescribeSecurityGroupsOutput, error) {
	request := nc.client.DescribeSecurityGroupsRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.DescribeSecurityGroupsOutput, nil
}

func (nc *nifcloud) RegisterInstancesWithSecurityGroup(ctx context.Context, input *computing.RegisterInstancesWithSecurityGroupInput) (*computing.RegisterInstancesWithSecurityGroupOutput, error) {
	request := nc.client.RegisterInstancesWithSecurityGroupRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.RegisterInstancesWithSecurityGroupOutput, nil
}

func (nc *nifcloud) DeregisterInstancesFromSecurityGroup(ctx context.Context, input *computing.DeregisterInstancesFromSecurityGroupInput) (*computing.DeregisterInstancesFromSecurityGroupOutput, error) {
	request := nc.client.DeregisterInstancesFromSecurityGroupRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.DeregisterInstancesFromSecurityGroupOutput, nil
}

func (nc *nifcloud) AuthorizeSecurityGroupIngress(ctx context.Context, input *computing.AuthorizeSecurityGroupIngressInput) (*computing.AuthorizeSecurityGroupIngressOutput, error) {
	request := nc.client.AuthorizeSecurityGroupIngressRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.AuthorizeSecurityGroupIngressOutput, nil
}

func (nc *nifcloud) RevokeSecurityGroupIngress(ctx context.Context, input *computing.RevokeSecurityGroupIngressInput) (*computing.RevokeSecurityGroupIngressOutput, error) {
	request := nc.client.RevokeSecurityGroupIngressRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.RevokeSecurityGroupIngressOutput, nil
}

func (nc *nifcloud) AssociateAddress(ctx context.Context, input *computing.AssociateAddressInput) (*computing.AssociateAddressOutput, error) {
	request := nc.client.AssociateAddressRequest(input)
	res, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	return res.AssociateAddressOutput, nil
}

func (nc *nifcloud) WaitUntilInstanceStopped(ctx context.Context, input *computing.DescribeInstancesInput) error {
	return nc.client.WaitUntilInstanceStopped(ctx, input, []aws.WaiterOption{}...)
}

func (nc *nifcloud) WaitUntilInstanceDeleted(ctx context.Context, input *computing.DescribeInstancesInput) error {
	return nc.client.WaitUntilInstanceDeleted(ctx, input, []aws.WaiterOption{}...)
}

func (nc *nifcloud) WaitUntilInstanceRunning(ctx context.Context, input *computing.DescribeInstancesInput) error {
	return nc.client.WaitUntilInstanceRunning(ctx, input, []aws.WaiterOption{}...)
}
