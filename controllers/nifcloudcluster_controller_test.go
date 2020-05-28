/*
Copyright 2020 FUJITSU CLOUD TECHNOLOGIES LIMITED. All Rights Reserved.

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
	infrav1alpha2 "github.com/nifcloud-labs/cluster-api-provider-nifcloud/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var _ = Describe("NifcloudClusterReconciler", func() {
	Context("Reconcile an NifcloudCluster", func() {
		It("should not error and not requeue the request", func() {
			ctx := context.Background()
			reconciler := &NifcloudClusterReconciler{
				Client: k8sClient,
				Log:    log.Log,
			}

			cluster := &infrav1alpha2.NifcloudCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "hoge", Namespace: "default"},
			}

			// Create the NifcloudCluster and expect the Reconcile to be created
			Expect(k8sClient.Create(ctx, cluster)).To(Succeed())

			result, err := reconciler.Reconcile(ctrl.Request{
				NamespacedName: client.ObjectKey{
					Namespace: cluster.Namespace,
					Name:      cluster.Name,
				},
			})
			Expect(err).To(BeNil())
			Expect(result.RequeueAfter).To(BeZero())
		})
	})
})
