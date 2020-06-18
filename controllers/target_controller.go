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
	// "fmt"

	"github.com/go-logr/logr"
	// corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/yuanying/tgtd-operator/api"
	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
	"github.com/yuanying/tgtd-operator/utils/tgtadm"
)

// TargetReconciler reconciles a Target object
type TargetReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	TgtAdm   tgtadm.TgtAdm
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

	if !target.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.deleteTarget(log, target)
	}

	if err := r.createOrUpdateTarget(log, target); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *TargetReconciler) deleteTarget(log logr.Logger, target *tgtdv1alpha1.Target) error {
	return nil
}

func (r *TargetReconciler) createOrUpdateTarget(log logr.Logger, target *tgtdv1alpha1.Target) error {
	if !containsFinalizer(target, api.TargetCleanupFinalizer) {
		controllerutil.AddFinalizer(target, api.TargetCleanupFinalizer)
		if err := r.Update(context.Background(), target); err != nil {
			log.Error(err, "Unable to update")
			return err
		}
	}
	tid, err := r.TgtAdm.GetTargetTid(target.Spec.IQN)
	if err != nil {
		log.Error(err, "Can't retrieve Target tid", "IQN", target.Spec.IQN)
		return err
	}
	if tid == -1 {
		log.Info("Target doesn't exist, try to create")
		tid, err := r.TgtAdm.FindNextAvailableTargetID()
		if err != nil {
			log.Error(err, "Can't get available target id")
			return err
		}
		if err = r.TgtAdm.CreateTarget(tid, target.Spec.IQN); err != nil {
			log.Error(err, "Can't create target")
			return err
		}
	}
	return nil
}

func containsFinalizer(target *tgtdv1alpha1.Target, finalizer string) bool {
	f := target.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return true
		}
	}
	return false
}

func (r *TargetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tgtdv1alpha1.Target{}).
		Complete(r)
}
