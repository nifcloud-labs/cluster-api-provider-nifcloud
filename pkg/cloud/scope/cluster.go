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
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/scope/nifcloud"
	"k8s.io/klog/klogr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterScopeParams include input paramater to create new scope for cluster
type ClusterScopeParams struct {
	NifcloudClients
	Client client.Client
	Logger logr.Logger

	Cluster         *clusterv1.Cluster
	NifcloudCluster *infrav1alpha2.NifcloudCluster
}

// NewClusterScope creates a new cluster scope from a specified prams
// This params gave from each Reconcile
func NewClusterScope(params ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, fmt.Errorf("fail to generate new scope from nil Cluster")
	}
	if params.NifcloudCluster == nil {
		return nil, fmt.Errorf("fail to generate new scope form nil NifcloudCluster")
	}
	if params.Logger == nil {
		params.Logger = klogr.New()
	}
	if params.NifcloudClients.Computing == nil {
		accKey := os.Getenv("NIFCLOUD_ACCESS_KEY")
		secKey := os.Getenv("NIFCLOUD_SECRET_KEY")
		region := os.Getenv("NIFCLOUD_REGION")
		cmpClient, err := nifcloud.New(accKey, secKey, region)
		if err != nil {
			return nil, fmt.Errorf("failed to create nifcloud client: %w", err)
		}
		params.NifcloudClients.Computing = cmpClient
	}

	// helper need to close scope
	helper, err := patch.NewHelper(params.NifcloudCluster, params.Client)
	if err != nil {
		return nil, fmt.Errorf("fail to init patch helper")
	}

	return &ClusterScope{
		client:          params.Client,
		Cluster:         params.Cluster,
		Logger:          params.Logger,
		NifcloudCluster: params.NifcloudCluster,
		NifcloudClients: params.NifcloudClients,
		patchHelper:     helper,
	}, nil
}

// ClusterScope
type ClusterScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	NifcloudClients
	NifcloudCluster *infrav1alpha2.NifcloudCluster
	Cluster         *clusterv1.Cluster
}

func (s *ClusterScope) Network() *infrav1alpha2.Network {
	return &s.NifcloudCluster.Status.Network
}

func (s *ClusterScope) SecurityGroups() map[infrav1alpha2.SecurityGroupRole]infrav1alpha2.SecurityGroup {
	return s.NifcloudCluster.Status.Network.SecurityGroups
}

func (s *ClusterScope) Name() string {
	return s.Cluster.Name
}

func (s *ClusterScope) Close() error {
	return s.patchHelper.Patch(context.TODO(), s.NifcloudCluster)
}
