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
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/scope"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/services"
	"github.com/nifcloud-labs/cluster-api-provider-nifcloud/pkg/cloud/services/computing"
	"golang.org/x/crypto/ssh"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NifcloudMachineReconciler reconciles a NifcloudMachine object
type NifcloudMachineReconciler struct {
	client.Client
	Log            logr.Logger
	Recorder       record.EventRecorder
	serviceFactory func(*scope.ClusterScope) services.NifcloudMachineInterface
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nifcloudmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nifcloudmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

func (r *NifcloudMachineReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.Background()
	logger := r.Log.WithValues("nifcloudmachine", req.NamespacedName)

	// fetch nifcloud machine
	nifcloudMachine := &infrav1alpha2.NifcloudMachine{}
	err := r.Get(ctx, req.NamespacedName, nifcloudMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch Machine
	machine, err := util.GetOwnerMachine(ctx, r.Client, nifcloudMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		logger.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("machine", machine.Name)

	// Fetch Cluster
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		logger.Info("Machine is missing cluster label or cluster does not exists")
		return ctrl.Result{}, err
	}

	logger = logger.WithValues("cluster", cluster.Name)

	nifcloudCluster := &infrav1alpha2.NifcloudCluster{}

	nifcloudClusterName := client.ObjectKey{
		Namespace: nifcloudMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, nifcloudClusterName, nifcloudCluster); err != nil {
		logger.Info("Nifcloud is not available yet")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("nifcloudCluster", nifcloudCluster.Name)

	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:          r.Client,
		Logger:          logger,
		Cluster:         cluster,
		NifcloudCluster: nifcloudCluster,
	})
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, err
	}

	machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
		Logger:          logger,
		Client:          r.Client,
		Cluster:         cluster,
		Machine:         machine,
		NifcloudCluster: nifcloudCluster,
		NifcloudMachine: nifcloudMachine,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create machine scope: %w", err)
	}

	defer func() {
		if err := machineScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !nifcloudMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(machineScope, clusterScope)
	}
	return r.reconcileMachine(ctx, machineScope, clusterScope)
}

func (r *NifcloudMachineReconciler) reconcileDelete(machineScope *scope.MachineScope, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	machineScope.Info("Handling delete NifcloudMachine")

	svc := r.getComputingService(clusterScope)
	instance, err := r.findInstance(machineScope, svc)
	if err != nil {
		return ctrl.Result{}, err
	}

	if instance == nil {
		machineScope.V(2).Info("Unable to locate Nifcloud Instance by ID")
		r.Recorder.Eventf(machineScope.NifcloudMachine, corev1.EventTypeWarning, "NoInstanceFound", "Unable to locate Nifcloud Instance by ID")
		machineScope.NifcloudMachine.Finalizers = util.Filter(machineScope.NifcloudMachine.Finalizers, infrav1alpha2.MachineFinalizer)
		return ctrl.Result{}, nil
	}

	machineScope.V(3).Info("Nifcloud server found matching deleted NifcludInstance", "instance-id", instance.ID)

	switch instance.State {
	case infrav1alpha2.InstanceStopped:
		machineScope.Info("Terminating Nifcloud server", "instance-id", instance.ID)
		if err := svc.TerminateInstance(instance.ID); err != nil {
			r.Recorder.Eventf(machineScope.NifcloudMachine, corev1.EventTypeWarning, "FailedTerminate", "FailedTerminate", "Failed to terminate server %q: %v", instance.ID, err)
			return ctrl.Result{}, fmt.Errorf("failed to terminate server: %w", err)
		}
		machineScope.Info("Nifcloud server successfully terminated", "instance-id", instance.ID)
		r.Recorder.Eventf(machineScope.NifcloudMachine, corev1.EventTypeNormal, "SuccessfullyTerminated", "Terminate instance %q", instance.ID)
	case infrav1alpha2.InstancePending:
		machineScope.Info("Nifcloud server is now shutting down, or terminating", "instance-id", instance.ID)
	default:
		machineScope.Info("Stopping and Terminating Nifcloud server", "instance-id", instance.ID)
		if err := svc.StopAndTerminateInstanceWithTimeout(instance.ID); err != nil {
			r.Recorder.Eventf(machineScope.NifcloudMachine, corev1.EventTypeWarning, "FailedStopAndTerminate", "Failed to stop and terminate server %q: %v", instance.ID, err)
			return ctrl.Result{}, fmt.Errorf("failed to stop and terminate server: %w", err)
		}
		machineScope.Info("Nifcloud server successfully terminated", "instance-id", instance.ID)
		r.Recorder.Eventf(machineScope.NifcloudMachine, corev1.EventTypeNormal, "SuccessfullyTerminated", "Terminate instance %q", instance.ID)
	}
	machineScope.NifcloudMachine.Finalizers = util.Filter(machineScope.NifcloudMachine.Finalizers, infrav1alpha2.MachineFinalizer)
	return ctrl.Result{}, nil
}

func (r *NifcloudMachineReconciler) reconcileMachine(ctx context.Context, machineScope *scope.MachineScope, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	machineScope.Info("Reconcile NifcloudMachine")
	if machineScope.NifcloudMachine.Status.ErrorReason != nil || machineScope.NifcloudMachine.Status.ErrorMessage != nil {
		machineScope.Info("Error state detected, skipping reconciliation")
		return ctrl.Result{}, nil
	}

	if !util.Contains(machineScope.NifcloudMachine.Finalizers, infrav1alpha2.MachineFinalizer) {
		machineScope.V(1).Info("Adding Cluster API Provider Nifcloud finalizer")
		machineScope.NifcloudMachine.Finalizers = append(machineScope.NifcloudMachine.Finalizers, infrav1alpha2.MachineFinalizer)
	}

	if !machineScope.Cluster.Status.InfrastructureReady {
		machineScope.Info("Cluster infrastructure is not ready yet")
		return ctrl.Result{}, nil
	}

	if machineScope.Machine.Spec.Bootstrap.Data == nil {
		machineScope.Info("Bootstrap data is not yet available")
		return ctrl.Result{}, nil
	}

	svc := r.getComputingService(clusterScope)

	instance, err := r.getOrCreate(machineScope, svc)
	if err != nil {
		return ctrl.Result{}, err
	}
	if instance == nil {
		machineScope.Info("Nifcloud instance cannot be found")
		machineScope.SetErrorReason(capierrors.UpdateMachineError)
		machineScope.SetErrorMessage(fmt.Errorf("Nifcloud instance cannot be found"))
		return ctrl.Result{}, nil
	}

	machineScope.SetProviderID(fmt.Sprintf("nifcloud:////%s", instance.UID))

	existingInstanceState := machineScope.GetInstanceState()
	machineScope.SetInstanceState(instance.State)

	instanceID := machineScope.GetInstanceID()
	if existingInstanceState == nil || *existingInstanceState != instance.State {
		machineScope.Info("Nifcloud instance state changed", "state", instance.State, "instance-id", *instanceID)
	}

	switch instance.State {
	case infrav1alpha2.InstancePending, infrav1alpha2.InstanceStopped:
		machineScope.UnsetSendBootstrap()
		break
	case infrav1alpha2.InstanceRunning:
		machineScope.SetReady()
	default:
		machineScope.SetNotReady()
		machineScope.UnsetSendBootstrap()
		machineScope.Info("nifcloud instance state is undefined", "state", instance.State, "instace-id", *instanceID)
		machineScope.SetErrorReason(capierrors.UpdateMachineError)
		machineScope.SetErrorMessage(errors.Errorf("nifcloud instance state %q is undefined", instance.State))
	}

	if instance.State == infrav1alpha2.InstanceStopped {
		machineScope.SetErrorReason(capierrors.UpdateMachineError)
		machineScope.SetErrorMessage(errors.Errorf("nifcloud instance state %q is unexpected", instance.State))
	}

	machineScope.SetAddresses(instance.Addresses)

	// send bootstrap data over ssh
	// because nifcldoud userData is limited 8KB
	if !machineScope.IsSendBootstrap() && machineScope.NifcloudMachine.Status.Ready {
		machineScope.Info("wait for remote machine provisioning")
		ip := instance.PublicIP
		if machineScope.IsControlPlane() {
			ip = machineScope.NifcloudCluster.Status.APIEndpoints[0].Host
		}
		time.Sleep(3 * time.Second)
		err := r.sendBootstrapDataWithSCP(machineScope, ip, "22")
		if err != nil {
			return ctrl.Result{}, err
		}
		machineScope.SetSendBootstrap()
		machineScope.Info("success to send bootstrap data to nifcloud server")
	}

	return ctrl.Result{}, nil
}

func (r *NifcloudMachineReconciler) getOrCreate(scope *scope.MachineScope, svc services.NifcloudMachineInterface) (*infrav1alpha2.Instance, error) {
	instance, err := r.findInstance(scope, svc)
	if err != nil {
		return nil, err
	}

	if instance == nil {
		scope.Info("Creating Nifcloud instance")
		instance, err = svc.CreateInstance(scope)
		if err != nil {
			return nil, fmt.Errorf("failed to create Nifcloud instance for NifcloudMachine %s/%s: %w", scope.Namespace(), scope.Name(), err)
		}
	}

	return instance, nil
}

func (r *NifcloudMachineReconciler) findInstance(scope *scope.MachineScope, svc services.NifcloudMachineInterface) (*infrav1alpha2.Instance, error) {
	// Parse the ProviderID
	_, err := noderefutil.NewProviderID(scope.GetProviderID())
	if err != nil && err != noderefutil.ErrEmptyProviderID {
		return nil, fmt.Errorf("failed to parse Spec.ProviderID: %w", err)
	}
	if err == nil {
		// ProviderID include instance UniqueID
		// we can only grep Instance by InstanceID, so use it instead of ProviderID
		instance, err := svc.InstanceIfExists(scope.GetInstanceID())
		if err != nil {
			return nil, fmt.Errorf("failed to query NifcloudMachine instance: %w", err)
		}
		return instance, nil
	}

	instance, err := svc.GetRunningInstanceByTag(scope)
	if err != nil {
		return nil, fmt.Errorf("failed to query NifcloudMachine instance by tags: %w", err)
	}

	return instance, nil
}

func (r *NifcloudMachineReconciler) getComputingService(scope *scope.ClusterScope) services.NifcloudMachineInterface {
	if r.serviceFactory != nil {
		return r.serviceFactory(scope)
	}
	return computing.NewService(scope)
}

func (r *NifcloudMachineReconciler) sendBootstrapDataWithSCP(scope *scope.MachineScope, host, port string) error {
	userData := scope.GetRawBootstrapData()
	reader := bytes.NewReader(userData)
	pass := os.Getenv("CLUSTER_API_PRIVATE_KEY_PASS")

	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	privateKey := os.Getenv("CLUSTER_API_SSH_KEY")
	clientConfig, err := auth.PrivateKeyWithPassphrase("root", []byte(pass), privateKey, ssh.InsecureIgnoreHostKey())
	if err != nil {
		return err
	}
	// Create a new SCP client
	client := scp.NewClient(host+":"+port, &clientConfig)
	// Connect to the remote server
	if err := client.Connect(); err != nil {
		return err
	}

	// Close client connection after the file has been copied
	defer client.Close()

	if err := client.Copy(reader, "/root/bootstrap.cfg", "0655", int64(len(userData))); err != nil {
		return err
	}
	return nil
}

func (r *NifcloudMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1alpha2.NifcloudMachine{}).
		Complete(r)
}
