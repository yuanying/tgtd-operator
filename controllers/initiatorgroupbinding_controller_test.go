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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
	ptr_util "github.com/yuanying/tgtd-operator/utils/ptr"
	"github.com/yuanying/tgtd-operator/utils/tgtadm"
)

const (
	testInitiatorNamePrefix    = "iqn.2020-06.cloud.unstable.test"
	testInitiatorAnnotationKey = "initiator-iqn"
	testNodeLabel              = "node-group"
)

var _ = Describe("InitiatorGroupBindingController", func() {

	var (
		tgtadm     = tgtadm.TgtAdmLonghornHelper{}
		igAddress1 = newInitiatorGroupWithAddress("address1", []string{"192.168.1.0/24", "192.168.2.0/24"})
		igAddress2 = newInitiatorGroupWithAddress("address2", []string{"192.168.3.0/24", "192.168.4.0/24"})
		igAddress3 = newInitiatorGroupWithAddress("address3", []string{"192.168.1.0/24", "192.168.3.0/24"})
	)

	BeforeEach(func() {
		var err error
		ctx := context.Background()
		err = k8sClient.Create(ctx, igAddress1)
		Expect(err).ToNot(HaveOccurred())
		err = k8sClient.Create(ctx, igAddress2)
		Expect(err).ToNot(HaveOccurred())
		err = k8sClient.Create(ctx, igAddress3)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		var err error
		ctx := context.Background()
		err = k8sClient.DeleteAllOf(ctx, &tgtdv1alpha1.InitiatorGroup{})
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
})

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

func newInitiatorGroupWithNodeName(name string) *tgtdv1alpha1.InitiatorGroup {
	ig := newInitiatorGroup(name)
	ig.Spec.InitiatorNameStrategy.Type = tgtdv1alpha1.NodeNameInitiatorNameStrategy
	ig.Spec.InitiatorNameStrategy.InitiatorNamePrefix = ptr_util.StringPtr(testInitiatorNamePrefix)
	return ig
}

func newInitiatorGroupWithAnnotationKey(name string) *tgtdv1alpha1.InitiatorGroup {
	ig := newInitiatorGroup(name)
	ig.Spec.InitiatorNameStrategy.Type = tgtdv1alpha1.AnnotationInitiatorNameStrategy
	ig.Spec.InitiatorNameStrategy.AnnotationKey = ptr_util.StringPtr(testInitiatorAnnotationKey)
	return ig
}
