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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/scope"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/services/computing"
)

// NifcloudClusterReconciler reconciles a NifcloudCluster object
type NifcloudClusterReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nifcloudclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nifcloudclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

func (r *NifcloudClusterReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.Background()
	log := r.Log.WithValues("nifcloudcluster", req.NamespacedName)

	// Fetch cluster resources
	log.Info("fetching Cluster Resources")
	nifcloudCluster := &infrav1alpha2.NifcloudCluster{}
	if err := r.Get(ctx, req.NamespacedName, nifcloudCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, nifcloudCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:          r.Client,
		Logger:          log,
		Cluster:         cluster,
		NifcloudCluster: nifcloudCluster,
	})
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create scope")
	}

	defer func() {
		if err := clusterScope.Close(); err != nil {
			reterr = err
		}
	}()

	if !nifcloudCluster.DeletionTimestamp.IsZero() {
		return reconcileDelete(clusterScope)
	}

	return reconcileCluster(clusterScope)
}

func reconcileCluster(clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	clusterScope.Info("Reconciling Cluster")

	nifcloudCluster := clusterScope.NifcloudCluster

	// Add Finalizer
	if !util.Contains(nifcloudCluster.Finalizers, infrav1alpha2.ClusterFinalizer) {
		nifcloudCluster.Finalizers = append(nifcloudCluster.Finalizers, infrav1alpha2.ClusterFinalizer)
	}

	svc := computing.NewService(clusterScope)

	if err := svc.ReconcileNetwork(); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to reconcile network for NifcloudCluster %s/%s", nifcloudCluster.Namespace, nifcloudCluster.Name)
	}

	nifcloudCluster.Status.Ready = true

	clusterScope.Info("Reconciled Cluster successfully")
	return ctrl.Result{}, nil
}

func reconcileDelete(clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	clusterScope.Info("Reconciling NifcloudCluster delete")

	svc := computing.NewService(clusterScope)

	if err := svc.DeleteEndpoint(); err != nil {
		return ctrl.Result{}, err
	}

	if err := svc.DeleteNetwork(); err != nil {
		return ctrl.Result{}, err
	}

	clusterScope.NifcloudCluster.Finalizers = util.Filter(clusterScope.NifcloudCluster.Finalizers, infrav1alpha2.ClusterFinalizer)

	return ctrl.Result{}, nil
}

func (r *NifcloudClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1alpha2.NifcloudCluster{}).
		Complete(r)
}
