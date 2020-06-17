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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
)

// TargetReconciler reconciles a Target object
type TargetReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	NodeName string
}

// +kubebuilder:rbac:groups=tgtd.unstable.cloud,resources=targets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tgtd.unstable.cloud,resources=targets/status,verbs=get;update;patch

func (r *TargetReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("target", req.NamespacedName)

	target := &tgtdv1alpha1.Target{}
	if err := r.Get(ctx, req.NamespacedName, target); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Unable to fetch Target - skipping")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to fetch Target")
		return ctrl.Result{}, err
	}
	if target.Spec.TargetNodeName != r.NodeName {
		log.Info("TargetNodeName is different -- skipping", "TargetNodeName", r.NodeName)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *TargetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tgtdv1alpha1.Target{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(e event.DeleteEvent) bool {
				return false
			},
		}).
		Complete(r)
}
