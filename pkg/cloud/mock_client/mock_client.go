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

// Code generated by MockGen. DO NOT EDIT.
// Source: clients.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	computing "github.com/aokumasan/nifcloud-sdk-go-v2/service/computing"
	gomock "github.com/golang/mock/gomock"
	cloud "github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// AllocateAddress mocks base method
func (m *MockClient) AllocateAddress(arg0 context.Context, arg1 *computing.AllocateAddressInput) (*computing.AllocateAddressOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllocateAddress", arg0, arg1)
	ret0, _ := ret[0].(*computing.AllocateAddressOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllocateAddress indicates an expected call of AllocateAddress
func (mr *MockClientMockRecorder) AllocateAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllocateAddress", reflect.TypeOf((*MockClient)(nil).AllocateAddress), arg0, arg1)
}

// ReleaseAddress mocks base method
func (m *MockClient) ReleaseAddress(arg0 context.Context, arg1 *computing.ReleaseAddressInput) (*computing.ReleaseAddressOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseAddress", arg0, arg1)
	ret0, _ := ret[0].(*computing.ReleaseAddressOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReleaseAddress indicates an expected call of ReleaseAddress
func (mr *MockClientMockRecorder) ReleaseAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseAddress", reflect.TypeOf((*MockClient)(nil).ReleaseAddress), arg0, arg1)
}

// DescribeAddresses mocks base method
func (m *MockClient) DescribeAddresses(arg0 context.Context, arg1 *computing.DescribeAddressesInput) (*computing.DescribeAddressesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeAddresses", arg0, arg1)
	ret0, _ := ret[0].(*computing.DescribeAddressesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeAddresses indicates an expected call of DescribeAddresses
func (mr *MockClientMockRecorder) DescribeAddresses(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeAddresses", reflect.TypeOf((*MockClient)(nil).DescribeAddresses), arg0, arg1)
}

// DisassociateAddress mocks base method
func (m *MockClient) DisassociateAddress(arg0 context.Context, arg1 *computing.DisassociateAddressInput) (*computing.DisassociateAddressOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisassociateAddress", arg0, arg1)
	ret0, _ := ret[0].(*computing.DisassociateAddressOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisassociateAddress indicates an expected call of DisassociateAddress
func (mr *MockClientMockRecorder) DisassociateAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisassociateAddress", reflect.TypeOf((*MockClient)(nil).DisassociateAddress), arg0, arg1)
}

// DescribeInstances mocks base method
func (m *MockClient) DescribeInstances(arg0 context.Context, arg1 *computing.DescribeInstancesInput) (*computing.DescribeInstancesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeInstances", arg0, arg1)
	ret0, _ := ret[0].(*computing.DescribeInstancesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeInstances indicates an expected call of DescribeInstances
func (mr *MockClientMockRecorder) DescribeInstances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeInstances", reflect.TypeOf((*MockClient)(nil).DescribeInstances), arg0, arg1)
}

// DescribeImages mocks base method
func (m *MockClient) DescribeImages(arg0 context.Context, arg1 *computing.DescribeImagesInput) (*computing.DescribeImagesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeImages", arg0, arg1)
	ret0, _ := ret[0].(*computing.DescribeImagesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeImages indicates an expected call of DescribeImages
func (mr *MockClientMockRecorder) DescribeImages(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeImages", reflect.TypeOf((*MockClient)(nil).DescribeImages), arg0, arg1)
}

// RunInstances mocks base method
func (m *MockClient) RunInstances(arg0 context.Context, arg1 *computing.RunInstancesInput) (*computing.RunInstancesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInstances", arg0, arg1)
	ret0, _ := ret[0].(*computing.RunInstancesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunInstances indicates an expected call of RunInstances
func (mr *MockClientMockRecorder) RunInstances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInstances", reflect.TypeOf((*MockClient)(nil).RunInstances), arg0, arg1)
}

// StopInstances mocks base method
func (m *MockClient) StopInstances(arg0 context.Context, arg1 *computing.StopInstancesInput) (*computing.StopInstancesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopInstances", arg0, arg1)
	ret0, _ := ret[0].(*computing.StopInstancesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StopInstances indicates an expected call of StopInstances
func (mr *MockClientMockRecorder) StopInstances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopInstances", reflect.TypeOf((*MockClient)(nil).StopInstances), arg0, arg1)
}

// TerminateInstances mocks base method
func (m *MockClient) TerminateInstances(arg0 context.Context, arg1 *computing.TerminateInstancesInput) (*computing.TerminateInstancesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TerminateInstances", arg0, arg1)
	ret0, _ := ret[0].(*computing.TerminateInstancesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TerminateInstances indicates an expected call of TerminateInstances
func (mr *MockClientMockRecorder) TerminateInstances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TerminateInstances", reflect.TypeOf((*MockClient)(nil).TerminateInstances), arg0, arg1)
}

// CreateSecurityGroup mocks base method
func (m *MockClient) CreateSecurityGroup(arg0 context.Context, arg1 *computing.CreateSecurityGroupInput) (*computing.CreateSecurityGroupOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSecurityGroup", arg0, arg1)
	ret0, _ := ret[0].(*computing.CreateSecurityGroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSecurityGroup indicates an expected call of CreateSecurityGroup
func (mr *MockClientMockRecorder) CreateSecurityGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecurityGroup", reflect.TypeOf((*MockClient)(nil).CreateSecurityGroup), arg0, arg1)
}

// DeleteSecurityGroup mocks base method
func (m *MockClient) DeleteSecurityGroup(arg0 context.Context, arg1 *computing.DeleteSecurityGroupInput) (*computing.DeleteSecurityGroupOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecurityGroup", arg0, arg1)
	ret0, _ := ret[0].(*computing.DeleteSecurityGroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSecurityGroup indicates an expected call of DeleteSecurityGroup
func (mr *MockClientMockRecorder) DeleteSecurityGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecurityGroup", reflect.TypeOf((*MockClient)(nil).DeleteSecurityGroup), arg0, arg1)
}

// DescribeSecurityGroups mocks base method
func (m *MockClient) DescribeSecurityGroups(arg0 context.Context, arg1 *computing.DescribeSecurityGroupsInput) (*computing.DescribeSecurityGroupsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeSecurityGroups", arg0, arg1)
	ret0, _ := ret[0].(*computing.DescribeSecurityGroupsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeSecurityGroups indicates an expected call of DescribeSecurityGroups
func (mr *MockClientMockRecorder) DescribeSecurityGroups(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeSecurityGroups", reflect.TypeOf((*MockClient)(nil).DescribeSecurityGroups), arg0, arg1)
}

// AuthorizeSecurityGroupIngress mocks base method
func (m *MockClient) AuthorizeSecurityGroupIngress(arg0 context.Context, arg1 *computing.AuthorizeSecurityGroupIngressInput) (*computing.AuthorizeSecurityGroupIngressOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthorizeSecurityGroupIngress", arg0, arg1)
	ret0, _ := ret[0].(*computing.AuthorizeSecurityGroupIngressOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthorizeSecurityGroupIngress indicates an expected call of AuthorizeSecurityGroupIngress
func (mr *MockClientMockRecorder) AuthorizeSecurityGroupIngress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthorizeSecurityGroupIngress", reflect.TypeOf((*MockClient)(nil).AuthorizeSecurityGroupIngress), arg0, arg1)
}

// RevokeSecurityGroupIngress mocks base method
func (m *MockClient) RevokeSecurityGroupIngress(arg0 context.Context, arg1 *computing.RevokeSecurityGroupIngressInput) (*computing.RevokeSecurityGroupIngressOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeSecurityGroupIngress", arg0, arg1)
	ret0, _ := ret[0].(*computing.RevokeSecurityGroupIngressOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RevokeSecurityGroupIngress indicates an expected call of RevokeSecurityGroupIngress
func (mr *MockClientMockRecorder) RevokeSecurityGroupIngress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeSecurityGroupIngress", reflect.TypeOf((*MockClient)(nil).RevokeSecurityGroupIngress), arg0, arg1)
}

// RegisterInstancesWithSecurityGroup mocks base method
func (m *MockClient) RegisterInstancesWithSecurityGroup(arg0 context.Context, arg1 *computing.RegisterInstancesWithSecurityGroupInput) (*computing.RegisterInstancesWithSecurityGroupOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterInstancesWithSecurityGroup", arg0, arg1)
	ret0, _ := ret[0].(*computing.RegisterInstancesWithSecurityGroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterInstancesWithSecurityGroup indicates an expected call of RegisterInstancesWithSecurityGroup
func (mr *MockClientMockRecorder) RegisterInstancesWithSecurityGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterInstancesWithSecurityGroup", reflect.TypeOf((*MockClient)(nil).RegisterInstancesWithSecurityGroup), arg0, arg1)
}

// DeregisterInstancesFromSecurityGroup mocks base method
func (m *MockClient) DeregisterInstancesFromSecurityGroup(arg0 context.Context, arg1 *computing.DeregisterInstancesFromSecurityGroupInput) (*computing.DeregisterInstancesFromSecurityGroupOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeregisterInstancesFromSecurityGroup", arg0, arg1)
	ret0, _ := ret[0].(*computing.DeregisterInstancesFromSecurityGroupOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeregisterInstancesFromSecurityGroup indicates an expected call of DeregisterInstancesFromSecurityGroup
func (mr *MockClientMockRecorder) DeregisterInstancesFromSecurityGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeregisterInstancesFromSecurityGroup", reflect.TypeOf((*MockClient)(nil).DeregisterInstancesFromSecurityGroup), arg0, arg1)
}

// AssociateAddress mocks base method
func (m *MockClient) AssociateAddress(arg0 context.Context, arg1 *computing.AssociateAddressInput) (*computing.AssociateAddressOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssociateAddress", arg0, arg1)
	ret0, _ := ret[0].(*computing.AssociateAddressOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssociateAddress indicates an expected call of AssociateAddress
func (mr *MockClientMockRecorder) AssociateAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssociateAddress", reflect.TypeOf((*MockClient)(nil).AssociateAddress), arg0, arg1)
}

// WaitUntilInstanceStopped mocks base method
func (m *MockClient) WaitUntilInstanceStopped(arg0 context.Context, arg1 *computing.DescribeInstancesInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilInstanceStopped", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitUntilInstanceStopped indicates an expected call of WaitUntilInstanceStopped
func (mr *MockClientMockRecorder) WaitUntilInstanceStopped(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilInstanceStopped", reflect.TypeOf((*MockClient)(nil).WaitUntilInstanceStopped), arg0, arg1)
}

// WaitUntilInstanceDeleted mocks base method
func (m *MockClient) WaitUntilInstanceDeleted(arg0 context.Context, arg1 *computing.DescribeInstancesInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilInstanceDeleted", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitUntilInstanceDeleted indicates an expected call of WaitUntilInstanceDeleted
func (mr *MockClientMockRecorder) WaitUntilInstanceDeleted(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilInstanceDeleted", reflect.TypeOf((*MockClient)(nil).WaitUntilInstanceDeleted), arg0, arg1)
}

// WaitUntilInstanceRunning mocks base method
func (m *MockClient) WaitUntilInstanceRunning(arg0 context.Context, arg1 *computing.DescribeInstancesInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitUntilInstanceRunning", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitUntilInstanceRunning indicates an expected call of WaitUntilInstanceRunning
func (mr *MockClientMockRecorder) WaitUntilInstanceRunning(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitUntilInstanceRunning", reflect.TypeOf((*MockClient)(nil).WaitUntilInstanceRunning), arg0, arg1)
}

// MockClientFactory is a mock of ClientFactory interface
type MockClientFactory struct {
	ctrl     *gomock.Controller
	recorder *MockClientFactoryMockRecorder
}

// MockClientFactoryMockRecorder is the mock recorder for MockClientFactory
type MockClientFactoryMockRecorder struct {
	mock *MockClientFactory
}

// NewMockClientFactory creates a new mock instance
func NewMockClientFactory(ctrl *gomock.Controller) *MockClientFactory {
	mock := &MockClientFactory{ctrl: ctrl}
	mock.recorder = &MockClientFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClientFactory) EXPECT() *MockClientFactoryMockRecorder {
	return m.recorder
}

// CreateClient mocks base method
func (m *MockClientFactory) CreateClient(arg0, arg1, arg2 string) (cloud.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateClient", arg0, arg1, arg2)
	ret0, _ := ret[0].(cloud.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateClient indicates an expected call of CreateClient
func (mr *MockClientFactoryMockRecorder) CreateClient(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateClient", reflect.TypeOf((*MockClientFactory)(nil).CreateClient), arg0, arg1, arg2)
}
