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

package main

import (

	// Stdlib
	"flag"
	"os"

	// Community
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/klogr"
	ctrl "sigs.k8s.io/controller-runtime"

	// Local
	infrastructurev1alpha3 "github.com/h0tbird/cluster-api-provider-metal/api/v1alpha3"
	"github.com/h0tbird/cluster-api-provider-metal/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = infrastructurev1alpha3.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {

	var (
		metricsAddr          string
		enableLeaderElection bool
		watchNamespace       string
	)

	//---------------------
	// Command line flags.
	//---------------------

	flag.StringVar(
		&metricsAddr,
		"metrics-addr",
		":8080",
		"The address the metric endpoint binds to.",
	)

	flag.BoolVar(
		&enableLeaderElection,
		"enable-leader-election",
		false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.",
	)

	flag.StringVar(
		&watchNamespace,
		"namespace",
		"",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.",
	)

	flag.Parse()

	//----------------------------------------
	// Sets the klogr logging implementation.
	//----------------------------------------

	ctrl.SetLogger(klogr.New())

	if watchNamespace != "" {
		setupLog.Info("Watching cluster-api objects only in namespace for reconciliation", "namespace", watchNamespace)
	}

	//-------------------------------
	// Setup the controller manager.
	//-------------------------------

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "f938d593.cluster.x-k8s.io",
		Namespace:          watchNamespace,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	//---------------------------------------------------------
	// Setup the BareMetalCluster controller with the manager.
	//---------------------------------------------------------

	if err = (&controllers.BareMetalClusterReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("BareMetalCluster"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "BareMetalCluster")
		os.Exit(1)
	}

	//---------------------------------------------------------
	// Setup the BareMetalMachine controller with the manager.
	//---------------------------------------------------------

	if err = (&controllers.BareMetalMachineReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("BareMetalMachine"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "BareMetalMachine")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	//--------------------
	// Start the manager.
	//--------------------

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
