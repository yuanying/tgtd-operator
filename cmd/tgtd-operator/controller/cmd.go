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

package controller

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/yuanying/tgtd-operator/controllers"
)

func NewControllerCmd(scheme *runtime.Scheme, metricsAddr *string, enableLeaderElection *bool) *cobra.Command {
	var (
		initiatorNamePrefix string
	)
	log := ctrl.Log.WithName("controller")
	c := &cobra.Command{
		Use:   "controller",
		Short: "Run controller",
		Long:  `Run controller`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Starting controller", "metricsAddr", metricsAddr, "enableLeaderElection", enableLeaderElection)
			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Scheme:             scheme,
				MetricsBindAddress: *metricsAddr,
				Port:               9443,
				LeaderElection:     *enableLeaderElection,
				LeaderElectionID:   "9a544d75.unstable.cloud",
			})
			if err != nil {
				log.Error(err, "unable to start manager")
				return err
			}

			if err = (&controllers.InitiatorGroupReconciler{
				Client:              mgr.GetClient(),
				Log:                 ctrl.Log.WithName("controllers").WithName("InitiatorGroup"),
				Scheme:              mgr.GetScheme(),
				Recorder:            mgr.GetEventRecorderFor("initiator-group-binding-controller"),
				InitiatorNamePrefix: initiatorNamePrefix,
			}).SetupWithManager(mgr); err != nil {
				log.Error(err, "unable to create controller", "controller", "InitiatorGroup")
				return err
			}

			log.Info("starting manager")
			if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
				log.Error(err, "problem running manager")
				return err
			}
			return nil
		},
	}

	c.Flags().StringVar(&initiatorNamePrefix, "initiator-name-prefix", "iqn.2020-06.cloud.unstable.tgtd", "Initiator name prefix to generate initiator name.")
	return c
}
