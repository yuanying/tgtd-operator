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
)

// InitiatorGroupReconciler reconciles a InitiatorGroupBinding object
type InitiatorGroupReconciler struct {
	client.Client
	Log                 logr.Logger
	Scheme              *runtime.Scheme
	Recorder            record.EventRecorder
	InitiatorNamePrefix string
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

	inis, err := r.fetchInitiators(log, ig)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, r.updateStatus(log, ig, inis)
}

func (r *InitiatorGroupReconciler) fetchInitiators(log logr.Logger, ig *tgtdv1alpha1.InitiatorGroup) ([]string, error) {
	nodes, err := r.fetchNodes(log, ig.Spec.NodeSelector)
	if err != nil {
		return nil, err
	}
	inits := make([]string, 0, len(nodes.Items))

	fn := r.genInitiatorNameFunc(log, ig)
	for i := range nodes.Items {
		n := &nodes.Items[i]
		name := fn(n)
		if name != nil {
			inits = append(inits, *name)
		}
	}
	sort.Strings(inits)
	return inits, nil
}

func (r *InitiatorGroupReconciler) genInitiatorNameFunc(log logr.Logger, ig *tgtdv1alpha1.InitiatorGroup) func(*corev1.Node) *string {
	var f func(*corev1.Node) *string
	if ig.Spec.InitiatorNameStrategy.Type == tgtdv1alpha1.AnnotationInitiatorNameStrategy {
		f = func(node *corev1.Node) *string {
			key := ig.Spec.InitiatorNameStrategy.AnnotationKey
			if key == nil {
				msg := "Annotation key (.Spec.InitiatorNameStrategy.AnnotationKey) must be specified"
				err := fmt.Errorf(msg)
				log.Error(err, msg)
				r.Recorder.Event(ig, corev1.EventTypeWarning, "AnnotationKeyMissing", msg)
				return nil
			}
			ans := node.Annotations
			if ans != nil {
				if v, ok := ans[*key]; ok {
					return &v
				}
			}
			return nil
		}
	} else {
		f = func(node *corev1.Node) *string {
			var prefix string
			if ig.Spec.InitiatorNameStrategy.InitiatorNamePrefix != nil {
				prefix = *ig.Spec.InitiatorNameStrategy.InitiatorNamePrefix
			} else {
				prefix = r.InitiatorNamePrefix
			}
			name := fmt.Sprintf("%s:%s", prefix, node.Name)
			return &name
		}
	}
	return f
}

func (r *InitiatorGroupReconciler) fetchNodes(log logr.Logger, selector map[string]string) (*corev1.NodeList, error) {
	ctx := context.Background()
	opts, err := r.newListOptions(selector)
	if err != nil {
		return nil, err
	}
	nodes := &corev1.NodeList{}
	if err := r.List(ctx, nodes, opts...); err != nil {
		log.Error(err, "Failed to list nodes")
		return nil, err
	}
	return nodes, nil
}

func (r *InitiatorGroupReconciler) newListOptions(selector map[string]string) ([]client.ListOption, error) {
	opts := make([]client.ListOption, 0)
	if selector == nil {
		return opts, nil
	}
	opt := client.MatchingLabels{}
	for k, v := range selector {
		opt[k] = v
	}
	opts = append(opts, opt)
	return opts, nil
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
