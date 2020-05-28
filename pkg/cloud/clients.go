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

package cloud

import (
	"context"

	"github.com/aokumasan/nifcloud-sdk-go-v2/service/computing"
)

//go:generate mockgen -source clients.go -destination mock_client/mock_client.go -package mock_client

type Client interface {
	AllocateAddress(context.Context, *computing.AllocateAddressInput) (*computing.AllocateAddressOutput, error)
	ReleaseAddress(context.Context, *computing.ReleaseAddressInput) (*computing.ReleaseAddressOutput, error)
	DescribeAddresses(context.Context, *computing.DescribeAddressesInput) (*computing.DescribeAddressesOutput, error)
	DisassociateAddress(context.Context, *computing.DisassociateAddressInput) (*computing.DisassociateAddressOutput, error)
	DescribeInstances(context.Context, *computing.DescribeInstancesInput) (*computing.DescribeInstancesOutput, error)
	DescribeImages(context.Context, *computing.DescribeImagesInput) (*computing.DescribeImagesOutput, error)
	RunInstances(context.Context, *computing.RunInstancesInput) (*computing.RunInstancesOutput, error)
	StopInstances(context.Context, *computing.StopInstancesInput) (*computing.StopInstancesOutput, error)
	TerminateInstances(context.Context, *computing.TerminateInstancesInput) (*computing.TerminateInstancesOutput, error)
	CreateSecurityGroup(context.Context, *computing.CreateSecurityGroupInput) (*computing.CreateSecurityGroupOutput, error)
	DeleteSecurityGroup(context.Context, *computing.DeleteSecurityGroupInput) (*computing.DeleteSecurityGroupOutput, error)
	DescribeSecurityGroups(context.Context, *computing.DescribeSecurityGroupsInput) (*computing.DescribeSecurityGroupsOutput, error)
	AuthorizeSecurityGroupIngress(context.Context, *computing.AuthorizeSecurityGroupIngressInput) (*computing.AuthorizeSecurityGroupIngressOutput, error)
	RevokeSecurityGroupIngress(context.Context, *computing.RevokeSecurityGroupIngressInput) (*computing.RevokeSecurityGroupIngressOutput, error)
	RegisterInstancesWithSecurityGroup(context.Context, *computing.RegisterInstancesWithSecurityGroupInput) (*computing.RegisterInstancesWithSecurityGroupOutput, error)
	DeregisterInstancesFromSecurityGroup(context.Context, *computing.DeregisterInstancesFromSecurityGroupInput) (*computing.DeregisterInstancesFromSecurityGroupOutput, error)
	AssociateAddress(context.Context, *computing.AssociateAddressInput) (*computing.AssociateAddressOutput, error)
	WaitUntilInstanceStopped(context.Context, *computing.DescribeInstancesInput) error
	WaitUntilInstanceDeleted(context.Context, *computing.DescribeInstancesInput) error
	WaitUntilInstanceRunning(context.Context, *computing.DescribeInstancesInput) error
}

type ClientFactory interface {
	CreateClient(string, string, string) (Client, error)
}
