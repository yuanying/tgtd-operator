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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
	ptr_util "github.com/yuanying/tgtd-operator/utils/ptr"
)

const (
	testInitiatorNamePrefix = "iqn.2020-06.cloud.unstable.test"
	testNodeLabel           = "node-group"
)

var _ = Describe("InitiatorGroupController", func() {

	var (
		initialNodes = []corev1.Node{
			*newNode("node1"),
			*newNodeWithLabel("node2", testNodeLabel, "group1"),
		}
	)

	Context("InitiatorGroup with Address", func() {

		It("should update .Status.Addresses", func() {
			var err error
			ctx := context.Background()
			key := types.NamespacedName{Name: "ig1"}

			By("Creating the InitiatorGroup")

			addresses := []string{"192.168.1.0/24", "192.168.2.0/24"}
			ig := newInitiatorGroupWithAddress(key.Name, addresses)
			err = k8sClient.Create(ctx, ig)
			Expect(err).ToNot(HaveOccurred())

			By("Checking created InitiatorGroup")
			fetched := &tgtdv1alpha1.InitiatorGroup{}
			Eventually(func() []string {
				k8sClient.Get(context.Background(), key, fetched)
				return fetched.Status.Addresses
			}, timeout, interval).Should(Equal(addresses))

			By("Updating the InitiatorGroup")
			addresses = []string{"192.168.3.0/24", "192.168.4.0/24"}
			update := newInitiatorGroupWithAddress(key.Name, addresses)
			err = k8sClient.Get(context.Background(), key, update)
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Update(context.Background(), update)).Should(Succeed())

			By("Checking updated InitiatorGroup")
			fetchedUpdated := &tgtdv1alpha1.InitiatorGroup{}
			Eventually(func() []string {
				k8sClient.Get(context.Background(), key, fetchedUpdated)
				return fetchedUpdated.Status.Addresses
			}, timeout, interval).Should(Equal(addresses))

			By("Deleting the InitialGroup")
			Eventually(func() error {
				f := &tgtdv1alpha1.InitiatorGroup{}
				k8sClient.Get(context.Background(), key, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())

			Eventually(func() error {
				f := &tgtdv1alpha1.InitiatorGroup{}
				return k8sClient.Get(context.Background(), key, f)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})

	Context("InitiatorGroup with strategy", func() {

		BeforeEach(func() {
			// var err error
			// ctx := context.Background()
			createNodes(initialNodes)
		})

		When("it doesn't have nodeSelector", func() {

			It("should register all nodes as initiator", func() {

			})
		})

		When("it has nodeSelector", func() {

			It("should register only labeled nodes as initiator", func() {

			})
		})
	})

	AfterEach(func() {
		var err error
		ctx := context.Background()
		err = k8sClient.DeleteAllOf(ctx, &tgtdv1alpha1.InitiatorGroup{})
		Expect(err).ToNot(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &tgtdv1alpha1.InitiatorGroup{})
		Expect(err).ToNot(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &corev1.Node{})
		Expect(err).ToNot(HaveOccurred())
	})

})

func createNodes(nodes []corev1.Node) {
	var err error
	ctx := context.Background()
	for i := range nodes {
		n := nodes[i]
		err = k8sClient.Create(ctx, &n)
		Expect(err).ToNot(HaveOccurred())
	}
}

func newNode(name string) *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func newNodeWithLabel(name, key, value string) *corev1.Node {
	node := newNode(name)
	node.Labels = make(map[string]string)
	node.Labels[key] = value
	return node
}

func newInitiatorGroup(name string) *tgtdv1alpha1.InitiatorGroup {
	return &tgtdv1alpha1.InitiatorGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func newInitiatorGroupWithAddress(name string, addresses []string) *tgtdv1alpha1.InitiatorGroup {
	ig := newInitiatorGroup(name)
	ig.Spec.Addresses = addresses
	return ig
}

func newInitiatorGroupSelectorWithNodeName(name string) *tgtdv1alpha1.InitiatorGroup {
	ig := newInitiatorGroup(name)
	ig.Spec.InitiatorNameStrategy.Type = tgtdv1alpha1.NodeNameInitiatorNameStrategy
	ig.Spec.InitiatorNameStrategy.InitiatorNamePrefix = ptr_util.StringPtr(testInitiatorNamePrefix)
	return ig
}
