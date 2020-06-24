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
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	if target.Spec.NodeName != r.NodeName {
		log.Info("NodeName is different -- skipping", "NodeName", r.NodeName)
		return ctrl.Result{}, nil
	}

	log = log.WithValues("IQN", target.Spec.IQN)

	if !target.DeletionTimestamp.IsZero() {
		if err := r.deleteTarget(log, target); err != nil {
			return ctrl.Result{}, err
		}
		controllerutil.RemoveFinalizer(target, api.TargetCleanupFinalizer)
		if err := r.Update(ctx, target); err != nil {
			log.Error(err, "Failed to update Target")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	t := metav1.Now()

	if err := r.createOrUpdateTarget(log, target); err != nil {
		target.SetConditionReason(tgtdv1alpha1.TargetTargetFailed, corev1.ConditionTrue, "UpdateFailed", err.Error(), t)
		r.updateStatus(log, target, t)
		return ctrl.Result{}, err
	}
	target.SetCondition(tgtdv1alpha1.TargetTargetFailed, corev1.ConditionFalse, t)

	actual, err := r.TgtAdm.GetTarget(target.Spec.IQN)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileLUNs(log, target, actual); err != nil {
		target.SetConditionReason(tgtdv1alpha1.TargetLUNFailed, corev1.ConditionTrue, "UpdateFailed", err.Error(), t)
		r.updateStatus(log, target, t)
		return ctrl.Result{}, err
	}
	target.SetCondition(tgtdv1alpha1.TargetLUNFailed, corev1.ConditionFalse, t)

	if err := r.updateStatus(log, target, t); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *TargetReconciler) deleteTarget(log logr.Logger, target *tgtdv1alpha1.Target) error {
	tid, err := r.TgtAdm.GetTargetTid(target.Spec.IQN)
	if err != nil {
		return err
	}
	if tid == -1 {
		log.V(4).Info("Target is already removed")
		return nil
	}
	err = r.TgtAdm.DeleteTarget(tid)
	if err != nil {
		log.Error(err, "Failed to remove Target")
		return err
	}
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

func (r *TargetReconciler) getLUN(lun int32, luns []tgtdv1alpha1.TargetLUN) *tgtdv1alpha1.TargetLUN {
	for i := range luns {
		l := &luns[i]
		if l.LUN == lun {
			return l
		}
	}
	return nil
}

func (r *TargetReconciler) reconcileLUNs(log logr.Logger, target *tgtdv1alpha1.Target, actual *tgtdv1alpha1.TargetActual) error {
	for i := range target.Spec.LUNs {
		l := &target.Spec.LUNs[i]
		al := r.getLUN(l.LUN, actual.LUNs)
		if al != nil {
			if al.BackingStore != l.BackingStore {
				log.V(1).Info("BackingStore is missmatch",
					"LUN", l.LUN,
					"DesiredBackingStore", l.BackingStore,
					"ObservedBackingStore", al.BackingStore,
				)
			}
		} else {
			if l.BSType != nil {
				bsopts := ""
				if l.BSOpts != nil {
					bsopts = *l.BSOpts
				}
				if err := r.TgtAdm.AddLun(int(actual.TID), int(l.LUN), l.BackingStore, *l.BSType, bsopts); err != nil {
					return err
				}
			} else {
				if err := r.TgtAdm.AddLunBackedByFile(int(actual.TID), int(l.LUN), l.BackingStore); err != nil {
					return err
				}
			}
		}
	}

	return r.deleteStaledLUNs(log, target, actual)
}

func (r *TargetReconciler) deleteStaledLUNs(log logr.Logger, target *tgtdv1alpha1.Target, actual *tgtdv1alpha1.TargetActual) error {
	for i := range actual.LUNs {
		al := &actual.LUNs[i]
		if al.LUN == 0 {
			log.V(4).Info("LID 0 should be ignored")
			continue
		}
		l := r.getLUN(al.LUN, target.Spec.LUNs)
		if l == nil {
			if err := r.TgtAdm.DeleteLun(int(actual.TID), int(al.LUN)); err != nil {
				log.Error(err, "Failed to delete LUN: %v, Path: %v", al.LUN, al.BackingStore)
				return err
			}
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

// updateStatus updates target.Status.
func (r *TargetReconciler) updateStatus(log logr.Logger, target *tgtdv1alpha1.Target, t metav1.Time) error {
	targetFailed := target.GetCondition(tgtdv1alpha1.TargetTargetFailed)
	lunsFailed := target.GetCondition(tgtdv1alpha1.TargetLUNFailed)

	ready := (targetFailed == nil || targetFailed.Status == corev1.ConditionFalse) &&
		(lunsFailed == nil || lunsFailed.Status == corev1.ConditionFalse)
	status := corev1.ConditionFalse
	if ready {
		status = corev1.ConditionTrue
	}

	target.SetCondition(tgtdv1alpha1.TargetConditionReady, status, t)

	target.Status.ObservedGeneration = target.Generation
	observedTarget, err := r.TgtAdm.GetTarget(target.Spec.IQN)
	if err != nil {
		log.Error(err, "Unable to get actual state of target")
		return err
	}
	target.Status.ObservedState = observedTarget

	if err := r.Status().Update(context.Background(), target); err != nil {
		log.Error(err, "Unable to update Target status", "observedTarget", observedTarget)
		return err
	}

	return nil
}
