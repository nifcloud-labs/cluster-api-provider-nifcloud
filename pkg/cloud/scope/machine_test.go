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

package scope

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/aokumasan/nifcloud-sdk-go-v2/nifcloud"
	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func setupScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := infrav1alpha2.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := clusterv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return scheme, nil
}

func newMachine(clusterName, machineName string) *clusterv1.Machine {
	return &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"cluster.x-k8s.io/cluster-name": clusterName,
			},
			ClusterName: clusterName,
			Name:        machineName,
			Namespace:   "default",
		},
		Spec: clusterv1.MachineSpec{
			Bootstrap: clusterv1.Bootstrap{
				Data: pointer.StringPtr(machineName),
			},
		},
	}
}

func newCluster(clusterName string) *clusterv1.Cluster {
	return &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterName,
			Namespace: "default",
		},
	}
}

func newNifcloudCluster(clusterName string) *infrav1alpha2.NifcloudCluster {
	return &infrav1alpha2.NifcloudCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterName,
			Namespace: "default",
		},
	}
}

func newNifcloudMachine(clusterName, machineName string) *infrav1alpha2.NifcloudMachine {
	return &infrav1alpha2.NifcloudMachine{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"cluster.x-k8s.io/cluster-name": clusterName,
			},
			Name:      clusterName,
			Namespace: "default",
		},
		Spec: infrav1alpha2.NifcloudMachineSpec{
			ProviderID: nifcloud.String("nifcloud:///test-instance-0"),
		},
	}
}

func setupMachineScope() (*MachineScope, error) {
	scheme, err := setupScheme()
	if err != nil {
		return nil, err
	}
	clusterName := "test-cluster"
	cluster := newCluster(clusterName)
	machine := newMachine(clusterName, "test-machine-0")
	nifcloudClsuter := newNifcloudCluster(clusterName)
	nifcloudMachine := newNifcloudMachine(clusterName, "test-machine-0")

	initObjects := []runtime.Object{
		cluster, machine, nifcloudClsuter, nifcloudMachine,
	}

	client := fake.NewFakeClientWithScheme(scheme, initObjects...)
	return NewMachineScope(
		MachineScopeParams{
			Client:          client,
			Machine:         machine,
			Cluster:         cluster,
			NifcloudCluster: nifcloudClsuter,
			NifcloudMachine: nifcloudMachine,
		},
	)
}

func TestGetUserDataIsBase64Encoded(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	userData, err := scope.GetUserData()
	if err != nil {
		t.Fatal(err)
	}
	_, err = base64.StdEncoding.DecodeString(userData)
	if err != nil {
		t.Fatalf("GetUserData is not base64 encoded: %v", err)
	}
}

func TestGetUserDataIsTemplateFilled(t *testing.T) {
	scope, err := setupMachineScope()
	if err != nil {
		t.Fatal(err)
	}

	userData, err := scope.GetUserData()
	if err != nil {
		t.Fatal(err)
	}
	d, err := base64.StdEncoding.DecodeString(userData)
	if err != nil {
		t.Fatal("some error occuered other test")
	}
	if len(d) == 0 {
		t.Fatalf("userdata is null: %v", d)
	}
	instanceID := scope.GetInstanceID()
	if !strings.Contains(string(d), *instanceID) {
		t.Fatalf("composed string does not contain instance-id [%v]", *instanceID)
	}
	fmt.Println(string(d))
}
