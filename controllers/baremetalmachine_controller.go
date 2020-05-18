/*
Copyright 2020 Open Source Community.

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

	// Stdlib
	"context"

	// Community
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	// Local
	infrav1 "github.com/h0tbird/cluster-api-provider-metal/api/v1alpha3"
)

// BareMetalMachineReconciler reconciles a BareMetalMachine object
type BareMetalMachineReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=baremetalmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=baremetalmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch

// Reconcile ...
func (r *BareMetalMachineReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.TODO()
	log := r.Log.WithValues("namespace", req.Namespace, "baremetalMachine", req.Name)

	// Fetch the BareMetalMachine instance.
	bareMetalMachine := &infrav1.BareMetalMachine{}
	err := r.Get(ctx, req.NamespacedName, bareMetalMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Machine instance.
	machine, err := util.GetOwnerMachine(ctx, r.Client, bareMetalMachine.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}
	log = log.WithValues("machine", machine.Name)

	// Fetch the Cluster instance.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, bareMetalMachine.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	// Cluster is paused or the object has the 'paused' annotation.
	if util.IsPaused(cluster, bareMetalMachine) {
		log.Info("BareMetalMachine or linked Cluster is marked as paused. Won't reconcile")
		return reconcile.Result{}, nil
	}
	log = log.WithValues("cluster", cluster.Name)

	// Do not requeue.
	return reconcile.Result{}, nil
}

// SetupWithManager ...
func (r *BareMetalMachineReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.BareMetalMachine{}).
		Complete(r)
}
