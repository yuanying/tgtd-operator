/*
Copyright 2020 O.Yuanying <yuanying@fraction.jp>

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
	"flag"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
	"github.com/yuanying/tgtd-operator/cmd/tgtd-operator/controller"
	"github.com/yuanying/tgtd-operator/cmd/tgtd-operator/node"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = tgtdv1alpha1.AddToScheme(scheme)
}

func main() {

	var metricsAddr string
	var enableLeaderElection bool

	rootCmd := &cobra.Command{
		Use: "tgtd-operator",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
			log := ctrl.Log.WithName("main")
			log.Info("staring tgtd-operator")
		},
	}

	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	rootCmd.PersistentFlags().StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	rootCmd.PersistentFlags().BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	rootCmd.AddCommand(controller.NewControllerCmd(scheme, &metricsAddr, &enableLeaderElection))
	rootCmd.AddCommand(node.NewNodeCmd(scheme, &metricsAddr, &enableLeaderElection))
	rootCmd.Execute()
}
