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
	"reflect"
	"testing"

	"github.com/aokumasan/nifcloud-sdk-go-v2/service/computing"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	nferrors "github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/errors"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/mock_client"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
)

func TestService_InstanceIfExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.TODO()

	tests := []struct {
		name       string
		instanceID string
		expect     func(m *mock_client.MockClientMockRecorder)
		check      func(instance *infrav1alpha2.Instance, err error)
	}{
		{
			name:       "does not exists",
			instanceID: "hoge",
			expect: func(m *mock_client.MockClientMockRecorder) {
				m.DescribeInstances(ctx, &computing.DescribeInstancesInput{
					InstanceId: []string{"hoge"},
				}).
					Return(nil, nferrors.NewNotFound(errors.New("not found")))
			},
			check: func(instance *infrav1alpha2.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if instance != nil {
					t.Fatalf("Did not expected result, but got something: %+v", instance)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := mock_client.NewMockClient(mockCtrl)

			scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Cluster: &clusterv1.Cluster{},
				NifcloudClients: scope.NifcloudClients{
					Computing: mockSvc,
				},
				NifcloudCluster: &infrav1alpha2.NifcloudCluster{
					Spec: infrav1alpha2.NifcloudClusterSpec{},
				},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tt.expect(mockSvc.EXPECT())

			service := NewService(scope)
			instance, err := service.InstanceIfExists(&tt.instanceID)
			tt.check(instance, err)
		})
	}
}

func TestService_GetRunningInstanceByTag(t *testing.T) {
	type fields struct {
		scope *scope.ClusterScope
	}
	type args struct {
		scope *scope.MachineScope
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *infrav1alpha2.Instance
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				scope: tt.fields.scope,
			}
			got, err := s.GetRunningInstanceByTag(tt.args.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetRunningInstanceByTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetRunningInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CreateInstance(t *testing.T) {
	type fields struct {
		scope *scope.ClusterScope
	}
	type args struct {
		scope *scope.MachineScope
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *infrav1alpha2.Instance
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				scope: tt.fields.scope,
			}
			got, err := s.CreateInstance(tt.args.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CreateInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_runInstance(t *testing.T) {
	type fields struct {
		scope *scope.ClusterScope
	}
	type args struct {
		role string
		i    *infrav1alpha2.Instance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *infrav1alpha2.Instance
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				scope: tt.fields.scope,
			}
			got, err := s.runInstance(tt.args.role, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.runInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.runInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SDKToInstance(t *testing.T) {
	type fields struct {
		scope *scope.ClusterScope
	}
	type args struct {
		v computing.InstancesSetItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *infrav1alpha2.Instance
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				scope: tt.fields.scope,
			}
			got, err := s.SDKToInstance(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SDKToInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.SDKToInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_defaultImageLookup(t *testing.T) {
	type fields struct {
		scope *scope.ClusterScope
	}
	type args struct {
		kubernetesVersion string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				scope: tt.fields.scope,
			}
			got, err := s.defaultImageLookup()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.defaultImageLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.defaultImageLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}
