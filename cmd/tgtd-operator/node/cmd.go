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

package node

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/yuanying/tgtd-operator/controllers"
	"github.com/yuanying/tgtd-operator/utils/tgtadm"
)

func NewNodeCmd(scheme *runtime.Scheme, metricsAddr *string, enableLeaderElection *bool) *cobra.Command {
	var (
		nodeName string
	)
	log := ctrl.Log.WithName("node")
	c := &cobra.Command{
		Use:   "node",
		Short: "Run node",
		Long:  `Run node`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Starting controller", "metricsAddr", metricsAddr, "enableLeaderElection", enableLeaderElection)

			if nodeName == "" {
				err := fmt.Errorf("Node name must be specified")
				log.Error(err, "Node name must be specified")
				return err
			}

			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Scheme:             scheme,
				MetricsBindAddress: *metricsAddr,
				Port:               9443,
				LeaderElection:     *enableLeaderElection,
				LeaderElectionID:   fmt.Sprintf("%s.9a544d75.unstable.cloud", nodeName),
			})
			if err != nil {
				log.Error(err, "unable to start manager")
				return err
			}

			if err = (&controllers.TargetReconciler{
				Client:   mgr.GetClient(),
				Log:      ctrl.Log.WithName("controllers").WithName("Target"),
				Scheme:   mgr.GetScheme(),
				Recorder: mgr.GetEventRecorderFor("target-controller"),
				TgtAdm:   &tgtadm.TgtAdmLonghornHelper{},
				NodeName: nodeName,
			}).SetupWithManager(mgr); err != nil {
				log.Error(err, "unable to create controller", "controller", "Target")
				return err
			}

			if err = (&controllers.InitiatorGroupBindingReconciler{
				Client:   mgr.GetClient(),
				Log:      ctrl.Log.WithName("controllers").WithName("InitiatorGroupBinding"),
				Scheme:   mgr.GetScheme(),
				Recorder: mgr.GetEventRecorderFor("initiator-group-binding-controller"),
				TgtAdm:   &tgtadm.TgtAdmLonghornHelper{},
				NodeName: nodeName,
			}).SetupWithManager(mgr); err != nil {
				log.Error(err, "unable to create controller", "controller", "InitiatorGroupBinding")
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

	c.Flags().StringVar(&nodeName, "node-name", "", "Node name where this agent is placed.")
	return c
}
