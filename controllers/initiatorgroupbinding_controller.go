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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
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

	log = log.WithValues("IQN", target.Spec.IQN)

	// your logic here

	return ctrl.Result{}, nil
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
		Watches(
			&source.Kind{Type: &corev1.Node{}},
			enqueueRequestsNodeForTarget(mgr.GetClient()),
		).
		Complete(r)
}

func enqueueRequestsNodeForTarget(c client.Client) *handler.EnqueueRequestsFromMapFunc {
	return &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
			targets := &tgtdv1alpha1.TargetList{}
			if err := c.List(context.Background(), targets); err != nil {
				return nil
			}
			reqs := make([]reconcile.Request, len(targets.Items))
			for i := range reqs {
				name := targets.Items[i].GetName()
				reqs[i].NamespacedName = types.NamespacedName{Name: name}
			}
			return reqs
		}),
	}
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
			ig := &tgtdv1alpha1.InitiatorGroup{}
			if err := c.Get(context.Background(), types.NamespacedName{Name: a.Meta.GetName()}, ig); err != nil {
				return nil
			}
			reqs := []reconcile.Request{}
			igbs := &tgtdv1alpha1.InitiatorGroupBindingList{}
			if err := c.List(context.Background(), igbs); err != nil {
				return nil
			}
			for i := range igbs.Items {
				igb := &igbs.Items[i]
				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: igb.Spec.TargetRef.Name}}
				reqs = append(reqs, req)
			}
			return reqs
		}),
	}
}
