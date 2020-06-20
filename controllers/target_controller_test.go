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
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/longhorn/go-iscsi-helper/iscsi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	// corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	tgtdv1alpha1 "github.com/yuanying/tgtd-operator/api/v1alpha1"
)

var _ = Describe("TargetController", func() {

	const (
		testRoot      = "/tmp/target"
		testImage     = "test.img"
		imageSize     = 4 * 1024 * 1024 // 4M
		testTargetIQN = "iqn-2020-06.cloud.unstable:target1"

		timeout  = time.Second * 30
		interval = time.Second * 1
	)

	var (
		imageFile = filepath.Join(testRoot, testImage)
	)

	BeforeEach(func() {
		err := exec.Command("mkdir", "-p", testRoot).Run()
		Expect(err).ToNot(HaveOccurred())

		imageFile = filepath.Join(testRoot, testImage)
		err = createTestFile(imageFile, imageSize)
		Expect(err).ToNot(HaveOccurred())

		err = exec.Command("mkfs.ext4", "-F", imageFile).Run()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := exec.Command("rm", "-rf", testRoot).Run()
		Expect(err).ToNot(HaveOccurred())
	})

	Context("Target with correct NodeName", func() {
		It("Should handle target and LUNs correctly", func() {
			spec := tgtdv1alpha1.TargetSpec{
				TargetNodeName: testNodeName,
				IQN:            testTargetIQN,
			}

			key := types.NamespacedName{
				Name: "target1",
			}

			toCreate := &tgtdv1alpha1.Target{
				ObjectMeta: metav1.ObjectMeta{
					Name: key.Name,
				},
				Spec: spec,
			}

			By("Creating the Target")
			Expect(k8sClient.Create(context.Background(), toCreate)).Should(Succeed())
			time.Sleep(time.Second * 5)

			fetched := &tgtdv1alpha1.Target{}
			Eventually(func() bool {
				k8sClient.Get(context.Background(), key, fetched)
				return fetched.Ready()
			}, timeout, interval).Should(BeTrue())

			Expect(iscsi.GetTargetTid(testTargetIQN)).ToNot(Equal(-1))
		})
	})
})

func createTestFile(file string, size int64) error {
	return exec.Command("truncate", "-s", strconv.FormatInt(size, 10), file).Run()
}
