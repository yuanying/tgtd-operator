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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
	"github.com/yuanying/tgtd-operator/utils/tgtadm"
)

var _ = Describe("InitiatorGroupBindingController", func() {

	var (
		tgtadm = tgtadm.TgtAdmLonghornHelper{}
		ig1Key = types.NamespacedName{Name: "address1"}
		ig1    *tgtdv1alpha1.InitiatorGroup
	)

	BeforeEach(func() {
		var err error
		ctx := context.Background()
		ig1 = newInitiatorGroupWithAddress(ig1Key.Name, []string{"192.168.1.0/24"})
		err = k8sClient.Create(ctx, ig1)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		var err error
		ctx := context.Background()
		err = k8sClient.DeleteAllOf(ctx, &tgtdv1alpha1.InitiatorGroup{})
		Expect(err).ToNot(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &tgtdv1alpha1.InitiatorGroupBinding{})
		Expect(err).ToNot(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &tgtdv1alpha1.Target{})
		Expect(err).ToNot(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &corev1.Node{})
		Expect(err).ToNot(HaveOccurred())
	})

	Context("InitiatorGroupBinding with no Target", func() {

		It("should reconcile InitiatorGroupBinding but do nothing", func() {
			spec := tgtdv1alpha1.InitiatorGroupBindingSpec{
				TargetRef:         tgtdv1alpha1.TargetReference{Name: "dont-exists"},
				InitiatorGroupRef: tgtdv1alpha1.InitiatorGroupReference{Name: "address1"},
			}

			key := types.NamespacedName{
				Name: "igb1",
			}

			beforeState, err := tgtadm.GetTargets()
			Expect(err).ToNot(HaveOccurred())

			toCreate := &tgtdv1alpha1.InitiatorGroupBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: key.Name,
				},
				Spec: spec,
			}

			Expect(k8sClient.Create(context.Background(), toCreate)).Should(Succeed())
			time.Sleep(time.Second * 5)

			afterState, err := tgtadm.GetTargets()
			Expect(err).ToNot(HaveOccurred())

			Expect(afterState).To(Equal(beforeState))
		})
	})

	Context("InitiatorGroupBinding with Target", func() {
		var (
			testTargetIQN = "iqn-2020-06.cloud.unstable:target1"
			targetKey     = types.NamespacedName{Name: "target1"}

			initialNodes []corev1.Node
			target       *tgtdv1alpha1.Target
		)

		BeforeEach(func() {
			target = newTarget(targetKey.Name, testTargetIQN)
			Expect(k8sClient.Create(context.Background(), target)).Should(Succeed())
			initialNodes = []corev1.Node{
				*newNode("node1"),
			}
			createNodes(initialNodes)
		})

		AfterEach(func() {
		})

		It("should bind/unbind initiators to Target", func() {
			var err error
			ctx := context.Background()
			key := types.NamespacedName{Name: "igb1"}
			toCreate := &tgtdv1alpha1.InitiatorGroupBinding{
				ObjectMeta: metav1.ObjectMeta{Name: key.Name},
				Spec: tgtdv1alpha1.InitiatorGroupBindingSpec{
					TargetRef:         tgtdv1alpha1.TargetReference{Name: targetKey.Name},
					InitiatorGroupRef: tgtdv1alpha1.InitiatorGroupReference{Name: ig1.Name},
				},
			}

			By("Creating InitiatorGroupBinding")
			Expect(k8sClient.Create(ctx, toCreate)).Should(Succeed())

			By("Checking target ACLs set correctly")
			Eventually(func() []string {
				actual, err := tgtadm.GetTarget(testTargetIQN)
				Expect(err).ToNot(HaveOccurred())
				if actual == nil {
					return []string{}
				}
				return actual.ACLs
			}, timeout, interval).Should(Equal([]string{"192.168.1.0/24", fmt.Sprintf("%s:node1", testInitiatorNamePrefix)}))

			By("Update InitiatorGroup")
			updateIG := &tgtdv1alpha1.InitiatorGroup{}
			err = k8sClient.Get(ctx, ig1Key, updateIG)
			Expect(err).ToNot(HaveOccurred())
			updateIG.Spec.Addresses = []string{"192.168.2.0/24", "192.168.3.0/24"}
			err = k8sClient.Update(ctx, updateIG)
			Expect(err).ToNot(HaveOccurred())

			By("Checking target ACLs set correctly")
			Eventually(func() []string {
				actual, err := tgtadm.GetTarget(testTargetIQN)
				Expect(err).ToNot(HaveOccurred())
				if actual == nil {
					return []string{}
				}
				return actual.ACLs
			}, timeout, interval).Should(Equal([]string{"192.168.2.0/24", "192.168.3.0/24", fmt.Sprintf("%s:node1", testInitiatorNamePrefix)}))
		})
	})
})

func newTarget(name, iqn string) *tgtdv1alpha1.Target {
	return &tgtdv1alpha1.Target{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: tgtdv1alpha1.TargetSpec{
			IQN:      iqn,
			NodeName: testNodeName,
		},
	}
}
