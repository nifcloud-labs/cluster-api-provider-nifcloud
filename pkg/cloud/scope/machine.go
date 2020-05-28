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
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"text/template"

	capierrors "sigs.k8s.io/cluster-api/errors"

	"github.com/aokumasan/nifcloud-sdk-go-v2/nifcloud"
	"github.com/aokumasan/nifcloud-sdk-go-v2/service/computing"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/userdata"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/klogr"
	"k8s.io/utils/pointer"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// nifcloud requirements
	maxInstanceIDName = 15
)

// MachineScopeParams include input paramater to create new scope for machine
type MachineScopeParams struct {
	NifcloudClient  *computing.Client
	Client          client.Client
	Logger          logr.Logger
	Cluster         *clusterv1.Cluster
	Machine         *clusterv1.Machine
	NifcloudCluster *infrav1alpha2.NifcloudCluster
	NifcloudMachine *infrav1alpha2.NifcloudMachine
}

// NewMachineScope creates a new machine scope from a specified prams
// This params gave from each Reconcile
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	if params.Client == nil {
		return nil, fmt.Errorf("client is required when creating a MachineScope")
	}
	if params.Machine == nil {
		return nil, fmt.Errorf("machine is required when creating a MachineScope")
	}
	if params.Cluster == nil {
		return nil, fmt.Errorf("cluster is required when creating a MachineScope")
	}
	if params.NifcloudMachine == nil {
		return nil, fmt.Errorf("nifcloud machine is required when creating a MachineScope")
	}
	if params.NifcloudCluster == nil {
		return nil, fmt.Errorf("nifcloud cluster is required when creating a MachineScope")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	helper, err := patch.NewHelper(params.NifcloudMachine, params.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to init patch helper: %w", err)
	}
	return &MachineScope{
		Logger:      params.Logger,
		client:      params.Client,
		patchHelper: helper,

		Cluster:         params.Cluster,
		Machine:         params.Machine,
		NifcloudCluster: params.NifcloudCluster,
		NifcloudMachine: params.NifcloudMachine,
	}, nil
}

// MachineScope
type MachineScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	Cluster         *clusterv1.Cluster
	Machine         *clusterv1.Machine
	NifcloudCluster *infrav1alpha2.NifcloudCluster
	NifcloudMachine *infrav1alpha2.NifcloudMachine
}

func (m *MachineScope) Name() string {
	return m.NifcloudMachine.Name
}

func (m *MachineScope) Namespace() string {
	return m.NifcloudMachine.Namespace
}

func (m *MachineScope) IsControlPlane() bool {
	return util.IsControlPlaneMachine(m.Machine)
}

func (m *MachineScope) Role() string {
	if util.IsControlPlaneMachine(m.Machine) {
		return "control-plane"
	}
	return "node"
}

func (m *MachineScope) GetProviderID() string {
	if m.NifcloudMachine.Spec.ProviderID != nil {
		return *m.NifcloudMachine.Spec.ProviderID
	}
	return ""
}

func (m *MachineScope) SetProviderID(v string) {
	m.NifcloudMachine.Spec.ProviderID = pointer.StringPtr(v)
}

// GetInstanceIDConved returns the expression of InstanceID in nifcloud
func (m *MachineScope) GetInstanceIDConved() string {
	instanceID := m.Name()
	hashed := md5.Sum([]byte(instanceID))
	tmp := fmt.Sprintf("%s", hex.EncodeToString(hashed[:]))
	return tmp[:maxInstanceIDName]
}

func (m *MachineScope) GetInstanceID() *string {
	return nifcloud.String(m.GetInstanceIDConved())
}

// GetInstanceUID returns Instance UID for ProviderID
func (m *MachineScope) GetInstanceUID() (*string, error) {
	parsed, err := noderefutil.NewProviderID(m.GetProviderID())
	if err != nil {
		return nil, err
	}
	return pointer.StringPtr(parsed.ID()), nil
}

func (m *MachineScope) GetInstanceState() *infrav1alpha2.InstanceState {
	return m.NifcloudMachine.Status.InstanceState
}

func (m *MachineScope) SetInstanceState(i infrav1alpha2.InstanceState) {
	m.NifcloudMachine.Status.InstanceState = &i
}

func (m *MachineScope) SetAddresses(addrs []corev1.NodeAddress) {
	m.NifcloudMachine.Status.Address = addrs
}

func (m *MachineScope) SetReady() {
	m.NifcloudMachine.Status.Ready = true
}

func (m *MachineScope) SetNotReady() {
	m.NifcloudMachine.Status.Ready = false
}

func (m *MachineScope) GetRawBootstrapData() []byte {
	return []byte(*m.Machine.Spec.Bootstrap.Data)
}

// GetBootstrapData returns bootstrap data from machines
func (m *MachineScope) GetBootstrapData() (string, error) {
	if m.Machine.Spec.Bootstrap.Data == nil {
		return "", fmt.Errorf("error retrieving bootstrap data: machine's bootstrap.dta is nil")
	}

	return base64.StdEncoding.EncodeToString(m.GetRawBootstrapData()), nil
}

// GetRawUserData returns userdata of each instance from script template
func (m *MachineScope) GetRawUserData() ([]byte, error) {
	tpl, err := template.New("userdata").Parse(userdata.ScriptTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create template")
	}
	// be careful to set instance-id: same name at `service/computing`
	instanceID := m.GetInstanceIDConved()
	var out bytes.Buffer
	if err := tpl.Execute(&out, map[string]string{
		"default_hostname": "{{ ds.meta_data.hostname }}",
		"instance_id":      instanceID,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to replace template signs")
	}
	return out.Bytes(), nil
}

func (m *MachineScope) GetUserData() (string, error) {
	d, err := m.GetRawUserData()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(d), nil
}

func (m *MachineScope) SetErrorReason(v capierrors.MachineStatusError) {
	m.NifcloudMachine.Status.ErrorReason = &v
}

func (m *MachineScope) SetErrorMessage(v error) {
	m.NifcloudMachine.Status.ErrorMessage = pointer.StringPtr(v.Error())
}

func (m *MachineScope) SetSendBootstrap() {
	m.NifcloudMachine.Status.SendBootstrap = true
}

func (m *MachineScope) UnsetSendBootstrap() {
	m.NifcloudMachine.Status.SendBootstrap = false
}

func (m *MachineScope) IsSendBootstrap() bool {
	return m.NifcloudMachine.Status.SendBootstrap
}

func (m *MachineScope) Close() error {
	return m.patchHelper.Patch(context.TODO(), m.NifcloudMachine)
}
