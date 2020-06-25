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
	"fmt"
	"net"
	"sort"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
	sliceutil "github.com/yuanying/tgtd-operator/utils/slice"
	"github.com/yuanying/tgtd-operator/utils/tgtadm"
)

// InitiatorGroupBindingReconciler reconciles a InitiatorGroupBinding object
type InitiatorGroupBindingReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	TgtAdm   tgtadm.TgtAdm
	NodeName string
}

// +kubebuilder:rbac:groups=tgtd.unstable.cloud,resources=initiatorgroupbindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tgtd.unstable.cloud,resources=initiatorgroupbindings/status,verbs=get;update;patch

func (r *InitiatorGroupBindingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
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

	if !target.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	log = log.WithValues("IQN", target.Spec.IQN)

	initiators, addresses, err := r.fetchInitiators(log, target)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileInitiators(log, target, initiators, addresses); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *InitiatorGroupBindingReconciler) reconcileInitiators(log logr.Logger, target *tgtdv1alpha1.Target, initiators, addresses []string) error {
	var err error
	iqn := target.Spec.IQN
	actual, err := r.TgtAdm.GetTarget(iqn)
	if err != nil {
		return fmt.Errorf("Failed to retrieve targets info: %v", iqn)
	}
	if actual == nil {
		log.Info("Target IQN has't registerd yet", "IQN", iqn)
		return nil
	}
	acls := map[string][]string{
		"IQN":     initiators,
		"address": addresses,
	}
	log = log.WithValues("Desired", acls, "Actual", actual.ACLs)
	log.V(4).Info("reconcileInitiators")

	for t, initiators := range acls {
		for _, i := range initiators {
			if r.containsInitiators(i, actual.ACLs) {
				log.V(1).Info("Already registerd, skipping", "initiator", i)
			} else {
				switch t {
				case "IQN":
					err = r.TgtAdm.BindInitiator(int(actual.TID), i)
				case "address":
					err = r.TgtAdm.BindInitiatorByAddress(int(actual.TID), i)
				}
				if err != nil {
					msg := "Failed to bind initiator"
					log.Error(err, msg, "initiator", i)
					r.Recorder.Eventf(target, corev1.EventTypeWarning, "BindInitiatorFailed", msg, "initiator", i)
					return err
				}
			}
		}
	}
	return r.deleteStaledInitiators(log, target, acls)
}

func (r *InitiatorGroupBindingReconciler) deleteStaledInitiators(log logr.Logger, target *tgtdv1alpha1.Target, acls map[string][]string) error {
	var err error
	iqn := target.Spec.IQN
	actual, err := r.TgtAdm.GetTarget(iqn)
	if err != nil {
		return fmt.Errorf("Failed to retrieve targets info: %v", iqn)
	}
	for _, acl := range actual.ACLs {
		if r.containsInitiators(acl, acls["IQN"]) {
			log.V(1).Info("Required", "initiator", acl)
			continue
		}
		if r.containsInitiators(acl, acls["address"]) {
			log.V(1).Info("Required", "initiator", acl)
			continue
		}
		var unbindErr error
		log.V(1).Info("Try to unbind", "initiator", acl)
		if acl == "ALL" {
			unbindErr = r.TgtAdm.UnbindInitiatorByAddress(int(actual.TID), acl)
		} else if _, _, err := net.ParseCIDR(acl); err == nil {
			unbindErr = r.TgtAdm.UnbindInitiatorByAddress(int(actual.TID), acl)
		} else if ip := net.ParseIP(acl); ip != nil {
			unbindErr = r.TgtAdm.UnbindInitiatorByAddress(int(actual.TID), acl)
		} else {
			unbindErr = r.TgtAdm.BindInitiator(int(actual.TID), acl)
		}
		if unbindErr != nil {
			msg := "Failed to unbind initiator"
			log.Error(err, msg, "initiator", acl)
			r.Recorder.Eventf(target, corev1.EventTypeWarning, "UnBindInitiatorFailed", msg, "initiator", acl)
			return unbindErr
		}
	}
	return err
}

func (r *InitiatorGroupBindingReconciler) containsInitiators(init string, acls []string) bool {
	for i := range acls {
		acl := &acls[i]
		if init == *acl {
			return true
		}
	}
	return false
}

func (r *InitiatorGroupBindingReconciler) fetchInitiators(log logr.Logger, target *tgtdv1alpha1.Target) (initiators []string, addresses []string, err error) {
	ctx := context.Background()
	igbs := &tgtdv1alpha1.InitiatorGroupBindingList{}
	if err := r.List(ctx, igbs); err != nil {
		log.Error(err, "Failed to list InitiatorGroupBinding")
		return nil, nil, err
	}
	for i := range igbs.Items {
		igb := &igbs.Items[i]
		ig := &tgtdv1alpha1.InitiatorGroup{}
		if err := r.Get(ctx, types.NamespacedName{Name: igb.Spec.InitiatorGroupRef.Name}, ig); err != nil {
			log.Error(err, fmt.Sprintf("Failed to get InitiatorGroup: %v", igb.Spec.InitiatorGroupRef.Name))
			return nil, nil, err
		}
		if ig.Status.Initiators != nil {
			initiators = append(initiators, ig.Status.Initiators...)
		}
		if ig.Status.Addresses != nil {
			addresses = append(addresses, ig.Status.Addresses...)
		}
	}

	sort.Strings(initiators)
	initiators = sliceutil.UniqueStrings(initiators)
	sort.Strings(addresses)
	addresses = sliceutil.UniqueStrings(addresses)
	return initiators, addresses, nil
}

func (r *InitiatorGroupBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tgtdv1alpha1.Target{}).
		Watches(
			&source.Kind{Type: &tgtdv1alpha1.InitiatorGroupBinding{}},
			enqueueRequestsIGBForTarget(mgr.GetClient()),
		).
		Watches(
			&source.Kind{Type: &tgtdv1alpha1.InitiatorGroup{}},
			enqueueRequestsIGForTarget(mgr.GetClient()),
		).
		Complete(r)
}

func enqueueRequestsIGBForTarget(c client.Client) *handler.EnqueueRequestsFromMapFunc {
	return &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
			igb := &tgtdv1alpha1.InitiatorGroupBinding{}
			if err := c.Get(context.Background(), types.NamespacedName{Name: a.Meta.GetName()}, igb); err != nil {
				return nil
			}
			return []reconcile.Request{{NamespacedName: types.NamespacedName{Name: igb.Spec.TargetRef.Name}}}
		}),
	}
}

func enqueueRequestsIGForTarget(c client.Client) *handler.EnqueueRequestsFromMapFunc {
	return &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
			reqs := []reconcile.Request{}
			igbs := &tgtdv1alpha1.InitiatorGroupBindingList{}
			if err := c.List(context.Background(), igbs); err != nil {
				return nil
			}
			for i := range igbs.Items {
				igb := &igbs.Items[i]
				if a.Meta.GetName() == igb.Spec.InitiatorGroupRef.Name {
					req := reconcile.Request{NamespacedName: types.NamespacedName{Name: igb.Spec.TargetRef.Name}}
					reqs = append(reqs, req)
				}
			}
			return reqs
		}),
	}
}
