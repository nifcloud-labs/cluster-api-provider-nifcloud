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

package services

import (
	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/scope"
)

type NifcloudMachineInterface interface {
	InstanceIfExists(id *string) (*infrav1alpha2.Instance, error)
	CreateInstance(scope *scope.MachineScope) (*infrav1alpha2.Instance, error)
	GetRunningInstanceByTag(scope *scope.MachineScope) (*infrav1alpha2.Instance, error)
	StopAndTerminateInstanceWithTimeout(id string) error
	TerminateInstance(id string) error
}
