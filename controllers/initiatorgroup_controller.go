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
)

// InitiatorGroupReconciler reconciles a InitiatorGroupBinding object
type InitiatorGroupReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=tgtd.unstable.cloud,resources=initiatorgroupbindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tgtd.unstable.cloud,resources=initiatorgroupbindings/status,verbs=get;update;patch

func (r *InitiatorGroupReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// var err error
	ctx := context.Background()
	log := r.Log.WithValues("InitiatorGroup", req.NamespacedName)

	ig := &tgtdv1alpha1.InitiatorGroup{}
	if err := r.Get(ctx, req.NamespacedName, ig); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Unable to fetch InitiatorGroup - skipping")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to fetch InitiatorGroup")
		return ctrl.Result{}, err
	}

	if !ig.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	// t := metav1.Now()

	return ctrl.Result{}, r.updateStatus(log, ig, nil)
}

func (r *InitiatorGroupReconciler) updateStatus(log logr.Logger, igs *tgtdv1alpha1.InitiatorGroup, initiators []string) error {
	igs.Status.Addresses = igs.Spec.Addresses
	igs.Status.Initiators = initiators

	if err := r.Status().Update(context.Background(), igs); err != nil {
		log.Error(err, "Unable to update status")
		return err
	}
	return nil
}

func (r *InitiatorGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tgtdv1alpha1.InitiatorGroup{}).
		Watches(
			&source.Kind{Type: &corev1.Node{}},
			enqueueRequestsNodeForIG(mgr.GetClient()),
		).
		Complete(r)
}

func enqueueRequestsNodeForIG(c client.Client) *handler.EnqueueRequestsFromMapFunc {
	return &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
			igs := &tgtdv1alpha1.InitiatorGroupList{}
			if err := c.List(context.Background(), igs); err != nil {
				return nil
			}
			reqs := make([]reconcile.Request, len(igs.Items))
			for i := range igs.Items {
				reqs[i].NamespacedName = types.NamespacedName{Name: igs.Items[i].GetName()}
			}
			return reqs
		}),
	}
}
