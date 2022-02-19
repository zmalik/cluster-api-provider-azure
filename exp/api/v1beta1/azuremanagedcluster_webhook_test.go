/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAzureManagedCluster_ValidateUpdate(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		name    string
		oldAMC  *AzureManagedCluster
		amc     *AzureManagedCluster
		wantErr bool
	}{
		{
			name: "AzureManagedCluster annotations passed as custom headers are immutable",
			oldAMC: &AzureManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"infrastructure.cluster.x-k8s.io/custom-header-SomeFeature": "true",
					},
				},
				Spec: AzureManagedClusterSpec{},
			},
			amc: &AzureManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"infrastructure.cluster.x-k8s.io/custom-header-SomeFeature": "false",
					},
				},
				Spec: AzureManagedClusterSpec{},
			},
			wantErr: true,
		},
		{
			name: "AzureManagedCluster annotations passed as custom headers are immutable 2",
			oldAMC: &AzureManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"infrastructure.cluster.x-k8s.io/custom-header-SomeFeature": "true",
					},
				},
				Spec: AzureManagedClusterSpec{},
			},
			amc: &AzureManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{},
				},
				Spec: AzureManagedClusterSpec{},
			},
			wantErr: true,
		},
		{
			name: "AzureManagedCluster annotations passed as custom headers are immutable 3",
			oldAMC: &AzureManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"infrastructure.cluster.x-k8s.io/custom-header-SomeFeature": "true",
					},
				},
				Spec: AzureManagedClusterSpec{},
			},
			amc: &AzureManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"infrastructure.cluster.x-k8s.io/custom-header-SomeFeature":    "true",
						"infrastructure.cluster.x-k8s.io/custom-header-AnotherFeature": "true",
					},
				},
				Spec: AzureManagedClusterSpec{},
			},
			wantErr: true,
		},
		{
			name: "AzureManagedCluster annotations not passed as custom headers are mutable",
			oldAMC: &AzureManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"annotation-a": "true",
						"infrastructure.cluster.x-k8s.io/custom-header-SomeFeature": "true",
					},
				},
				Spec: AzureManagedClusterSpec{},
			},
			amc: &AzureManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"infrastructure.cluster.x-k8s.io/custom-header-SomeFeature": "true",
						"annotation-b": "true",
					},
				},
				Spec: AzureManagedClusterSpec{},
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.amc.ValidateUpdate(tc.oldAMC)
			if tc.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
		})
	}
}
