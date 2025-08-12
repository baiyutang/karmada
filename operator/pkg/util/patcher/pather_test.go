/*
Copyright 2024 The Karmada Authors.

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

package patcher

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/karmada-io/karmada/operator/pkg/apis/operator/v1alpha1"
	"github.com/karmada-io/karmada/operator/pkg/constants"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func TestPatchForDeployment(t *testing.T) {
	tests := []struct {
		name       string
		patcher    *Patcher
		deployment *appsv1.Deployment
		want       *appsv1.Deployment
	}{
		{
			name: "PatchForDeployment_WithComponentLabelsAndAnnotations_Patched",
			patcher: &Patcher{
				componentLabels: map[string]string{
					"label1": "value1-patched",
				},
				annotations: map[string]string{
					"annotation1": "annot1-patched",
				},
			},
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test",
					Labels: map[string]string{
						"label1": "value1",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								"annotation1": "annot1",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test",
					Labels: map[string]string{
						"label1": "value1-patched",
					},
					Annotations: map[string]string{
						"annotation1": "annot1-patched",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"label1": "value1-patched",
							},
							Annotations: map[string]string{
								"annotation1": "annot1-patched",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PatchForDeployment_WithResourcesExtraArgsAndFeatureGates_Patched",
			patcher: &Patcher{
				extraArgs: map[string]string{
					"some-arg": "some-value",
				},
				featureGates: map[string]bool{
					"SomeGate": true,
				},
				resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
			},
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
									Command: []string{
										"/bin/bash",
										"--feature-gates=OldGate=false",
									},
								},
							},
						},
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-deployment",
					Namespace:   "test",
					Labels:      map[string]string{},
					Annotations: map[string]string{},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{},
							Annotations: map[string]string{},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
									Command: []string{
										"/bin/bash",
										"--feature-gates=OldGate=false,SomeGate=true",
										"--some-arg=some-value",
									},
									Resources: corev1.ResourceRequirements{
										Limits: corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("500m"),
											corev1.ResourceMemory: resource.MustParse("128Mi"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PatchForDeployment_WithExtraVolumesAndVolumeMounts_Patched",
			patcher: &Patcher{
				extraVolumes: []corev1.Volume{
					{
						Name: "extra-volume",
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{},
						},
					},
				},
				extraVolumeMounts: []corev1.VolumeMount{
					{
						Name:      "extra-volume",
						MountPath: "/extra/path",
					},
				},
			},
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-deployment",
					Namespace:   "test",
					Labels:      map[string]string{},
					Annotations: map[string]string{},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{},
							Annotations: map[string]string{},
						},
						Spec: corev1.PodSpec{
							Volumes: []corev1.Volume{
								{
									Name: "extra-volume",
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{},
									},
								},
							},
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "extra-volume",
											MountPath: "/extra/path",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set component info to ensure system labels are added
			test.patcher.WithComponentInfo("test-component", "test-instance")
			test.patcher.ForDeployment(test.deployment)

			// Check that all expected labels are present
			for key, value := range test.want.Labels {
				if deploymentValue, exists := test.deployment.Labels[key]; !exists || deploymentValue != value {
					t.Errorf("expected label %s=%s, but got %s or doesn't exist", key, value, deploymentValue)
				}
			}

			// Check that all expected pod template labels are present
			for key, value := range test.want.Spec.Template.Labels {
				if podValue, exists := test.deployment.Spec.Template.Labels[key]; !exists || podValue != value {
					t.Errorf("expected pod template label %s=%s, but got %s or doesn't exist", key, value, podValue)
				}
			}

			// Check that system labels are present
			if _, exists := test.deployment.Labels["app.kubernetes.io/managed-by"]; !exists {
				t.Error("expected system label 'app.kubernetes.io/managed-by' to be present")
			}

			if _, exists := test.deployment.Spec.Template.Labels["app.kubernetes.io/managed-by"]; !exists {
				t.Error("expected pod template system label 'app.kubernetes.io/managed-by' to be present")
			}
		})
	}
}

func TestPatchForStatefulSet(t *testing.T) {
	tests := []struct {
		name        string
		patcher     *Patcher
		statefulSet *appsv1.StatefulSet
		want        *appsv1.StatefulSet
	}{
		{
			name: "PatchForStatefulSet_WithComponentLabelsAndAnnotations_Patched",
			patcher: &Patcher{
				componentLabels: map[string]string{
					"label1": "value1-patched",
				},
				annotations: map[string]string{
					"annotation1": "annot1-patched",
				},
			},
			statefulSet: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-statefulset",
					Namespace: "test",
					Labels: map[string]string{
						"label1": "value1",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
			want: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-statefulset",
					Namespace: "test",
					Labels: map[string]string{
						"label1": "value1-patched",
					},
					Annotations: map[string]string{
						"annotation1": "annot1-patched",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"label1": "value1-patched",
							},
							Annotations: map[string]string{
								"annotation1": "annot1-patched",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PatchForStatefulSet_WithVolumes_Patched",
			patcher: &Patcher{
				volume: &v1alpha1.VolumeData{
					EmptyDir: &corev1.EmptyDirVolumeSource{
						Medium:    corev1.StorageMediumMemory,
						SizeLimit: &resource.Quantity{},
					},
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/tmp",
						Type: ptr.To(corev1.HostPathDirectory),
					},
					VolumeClaim: &corev1.PersistentVolumeClaimTemplate{
						Spec: corev1.PersistentVolumeClaimSpec{
							AccessModes: []corev1.PersistentVolumeAccessMode{
								corev1.ReadWriteOnce,
							},
							Resources: corev1.VolumeResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceStorage: resource.MustParse("1024m"),
								},
							},
						},
					},
				},
			},
			statefulSet: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-statefulset",
					Namespace: "test",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
							Volumes: []corev1.Volume{},
						},
					},
				},
			},
			want: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-statefulset",
					Namespace:   "test",
					Labels:      map[string]string{},
					Annotations: map[string]string{},
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{},
							Annotations: map[string]string{},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: constants.EtcdDataVolumeName,
									VolumeSource: corev1.VolumeSource{
										EmptyDir: &corev1.EmptyDirVolumeSource{},
									},
								},
								{
									Name: constants.EtcdDataVolumeName,
									VolumeSource: corev1.VolumeSource{
										HostPath: &corev1.HostPathVolumeSource{
											Path: "/tmp",
											Type: ptr.To(corev1.HostPathDirectory),
										},
									},
								},
							},
						},
					},
					VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: constants.EtcdDataVolumeName,
							},
							Spec: corev1.PersistentVolumeClaimSpec{
								AccessModes: []corev1.PersistentVolumeAccessMode{
									corev1.ReadWriteOnce,
								},
								Resources: corev1.VolumeResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceStorage: resource.MustParse("1024m"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PatchForStatefulSet_WithResourcesAndExtraArgs_Patched",
			patcher: &Patcher{
				extraArgs: map[string]string{
					"some-arg": "some-value",
				},
				resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
			},
			statefulSet: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-statefulset",
					Namespace: "test",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
			want: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-statefulset",
					Namespace:   "test",
					Labels:      map[string]string{},
					Annotations: map[string]string{},
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      map[string]string{},
							Annotations: map[string]string{},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
									Command: []string{
										"--some-arg=some-value",
									},
									Resources: corev1.ResourceRequirements{
										Limits: corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("500m"),
											corev1.ResourceMemory: resource.MustParse("128Mi"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set component info to ensure system labels are added
			test.patcher.WithComponentInfo("test-component", "test-instance")
			test.patcher.ForStatefulSet(test.statefulSet)

			// Check that all expected labels are present
			for key, value := range test.want.Labels {
				if statefulSetValue, exists := test.statefulSet.Labels[key]; !exists || statefulSetValue != value {
					t.Errorf("expected label %s=%s, but got %s or doesn't exist", key, value, statefulSetValue)
				}
			}

			// Check that all expected pod template labels are present
			for key, value := range test.want.Spec.Template.Labels {
				if podValue, exists := test.statefulSet.Spec.Template.Labels[key]; !exists || podValue != value {
					t.Errorf("expected pod template label %s=%s, but got %s or doesn't exist", key, value, podValue)
				}
			}

			// Check that system labels are present
			if _, exists := test.statefulSet.Labels["app.kubernetes.io/managed-by"]; !exists {
				t.Error("expected system label 'app.kubernetes.io/managed-by' to be present")
			}

			if _, exists := test.statefulSet.Spec.Template.Labels["app.kubernetes.io/managed-by"]; !exists {
				t.Error("expected pod template system label 'app.kubernetes.io/managed-by' to be present")
			}
		})
	}
}

// TestPatcherWithGlobalLabels tests the global labels functionality
func TestPatcherWithGlobalLabels(t *testing.T) {
	tests := []struct {
		name       string
		patcher    *Patcher
		deployment *appsv1.Deployment
		want       *appsv1.Deployment
	}{
		{
			name: "WithGlobalLabels_AppliedToDeployment",
			patcher: &Patcher{
				globalLabels: map[string]string{
					"environment": "production",
					"team":        "platform",
				},
				componentLabels: map[string]string{
					"component": "karmada-controller-manager",
				},
			},
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test",
					Labels: map[string]string{
						"existing": "label",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"existing": "label",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test",
					Labels: map[string]string{
						"existing":    "label",
						"environment": "production",
						"team":        "platform",
						"component":   "karmada-controller-manager",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"existing":    "label",
								"environment": "production",
								"team":        "platform",
								"component":   "karmada-controller-manager",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "WithComponentInfo_SystemLabelsApplied",
			patcher: &Patcher{
				globalLabels: map[string]string{
					"environment": "staging",
				},
				componentLabels: map[string]string{
					"component": "karmada-scheduler",
				},
			},
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test",
					Labels: map[string]string{
						"environment":                  "staging",
						"component":                    "karmada-scheduler",
						"app.kubernetes.io/managed-by": "karmada-operator",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"environment":                  "staging",
								"component":                    "karmada-scheduler",
								"app.kubernetes.io/managed-by": "karmada-operator",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test-container",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set component info for the second test case
			if test.name == "WithComponentInfo_SystemLabelsApplied" {
				test.patcher.WithComponentInfo("karmada-scheduler", "test-karmada")
			}

			test.patcher.ForDeployment(test.deployment)

			// Check that all expected labels are present
			for key, value := range test.want.Labels {
				if deploymentValue, exists := test.deployment.Labels[key]; !exists || deploymentValue != value {
					t.Errorf("expected label %s=%s, but got %s or doesn't exist", key, value, deploymentValue)
				}
			}

			// Check that all expected pod template labels are present
			for key, value := range test.want.Spec.Template.Labels {
				if podValue, exists := test.deployment.Spec.Template.Labels[key]; !exists || podValue != value {
					t.Errorf("expected pod template label %s=%s, but got %s or doesn't exist", key, value, podValue)
				}
			}
		})
	}
}

// TestPatcherForService tests the Service patching functionality
func TestPatcherForService(t *testing.T) {
	tests := []struct {
		name    string
		patcher *Patcher
		service *corev1.Service
		want    *corev1.Service
	}{
		{
			name: "ForService_WithGlobalAndComponentLabels",
			patcher: &Patcher{
				globalLabels: map[string]string{
					"environment": "production",
					"region":      "us-west-2",
				},
				componentLabels: map[string]string{
					"component": "karmada-apiserver",
				},
				annotations: map[string]string{
					"service.alpha.kubernetes.io/aws-load-balancer-type": "nlb",
				},
			},
			service: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
					Namespace: "test",
					Labels: map[string]string{
						"existing": "label",
					},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port: 80,
						},
					},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
					Namespace: "test",
					Labels: map[string]string{
						"existing":    "label",
						"environment": "production",
						"region":      "us-west-2",
						"component":   "karmada-apiserver",
					},
					Annotations: map[string]string{
						"service.alpha.kubernetes.io/aws-load-balancer-type": "nlb",
					},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port: 80,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.patcher.ForService(test.service)

			// Check that all expected labels are present
			for key, value := range test.want.Labels {
				if serviceValue, exists := test.service.Labels[key]; !exists || serviceValue != value {
					t.Errorf("expected label %s=%s, but got %s or doesn't exist", key, value, serviceValue)
				}
			}

			// Check that all expected annotations are present
			for key, value := range test.want.Annotations {
				if serviceValue, exists := test.service.Annotations[key]; !exists || serviceValue != value {
					t.Errorf("expected annotation %s=%s, but got %s or doesn't exist", key, value, serviceValue)
				}
			}
		})
	}
}

// TestPatcherForSecret tests the Secret patching functionality
func TestPatcherForSecret(t *testing.T) {
	test := struct {
		name           string
		patcher        *Patcher
		secret         *corev1.Secret
		expectedLabels map[string]string
	}{
		name: "secret with global and component labels",
		patcher: NewPatcher().
			WithGlobalLabels(map[string]string{"env": "prod", "team": "platform"}).
			WithComponentLabels(labels.Set{"component": "secret", "tier": "data"}).
			WithComponentInfo("test-secret", "test-instance"),
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-secret",
				Namespace: "default",
				Labels:    map[string]string{"existing": "label"},
			},
		},
		expectedLabels: map[string]string{
			"existing":                     "label",
			"env":                          "prod",
			"team":                         "platform",
			"component":                    "secret",
			"tier":                         "data",
			"app.kubernetes.io/managed-by": "karmada-operator",
			"app.kubernetes.io/name":       "test-secret",
			"app.kubernetes.io/instance":   "test-instance",
		},
	}

	test.patcher.ForSecret(test.secret)

	// Verify labels are correctly merged
	for key, expectedValue := range test.expectedLabels {
		if actualValue, exists := test.secret.Labels[key]; !exists || actualValue != expectedValue {
			t.Errorf("Expected label %s=%s, but got %s", key, expectedValue, actualValue)
		}
	}
}

func TestPatcherForServiceAccount(t *testing.T) {
	test := struct {
		name           string
		patcher        *Patcher
		serviceAccount *corev1.ServiceAccount
		expectedLabels map[string]string
	}{
		name: "serviceaccount with global and component labels",
		patcher: NewPatcher().
			WithGlobalLabels(map[string]string{"env": "prod", "team": "platform"}).
			WithComponentLabels(labels.Set{"component": "serviceaccount", "tier": "rbac"}).
			WithComponentInfo("test-sa", "test-instance"),
		serviceAccount: &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-sa",
				Namespace: "default",
				Labels:    map[string]string{"existing": "label"},
			},
		},
		expectedLabels: map[string]string{
			"existing":                     "label",
			"env":                          "prod",
			"team":                         "platform",
			"component":                    "serviceaccount",
			"tier":                         "rbac",
			"app.kubernetes.io/managed-by": "karmada-operator",
			"app.kubernetes.io/name":       "test-sa",
			"app.kubernetes.io/instance":   "test-instance",
		},
	}

	test.patcher.ForServiceAccount(test.serviceAccount)

	// Verify labels are correctly merged
	for key, expectedValue := range test.expectedLabels {
		if actualValue, exists := test.serviceAccount.Labels[key]; !exists || actualValue != expectedValue {
			t.Errorf("Expected label %s=%s, but got %s", key, expectedValue, actualValue)
		}
	}
}

func TestPatcherForClusterRole(t *testing.T) {
	test := struct {
		name           string
		patcher        *Patcher
		clusterRole    *rbacv1.ClusterRole
		expectedLabels map[string]string
	}{
		name: "clusterrole with global and component labels",
		patcher: NewPatcher().
			WithGlobalLabels(map[string]string{"env": "prod", "team": "platform"}).
			WithComponentLabels(labels.Set{"component": "clusterrole", "tier": "rbac"}).
			WithComponentInfo("test-cr", "test-instance"),
		clusterRole: &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "test-cr",
				Labels: map[string]string{"existing": "label"},
			},
		},
		expectedLabels: map[string]string{
			"existing":                     "label",
			"env":                          "prod",
			"team":                         "platform",
			"component":                    "clusterrole",
			"tier":                         "rbac",
			"app.kubernetes.io/managed-by": "karmada-operator",
			"app.kubernetes.io/name":       "test-cr",
			"app.kubernetes.io/instance":   "test-instance",
		},
	}

	test.patcher.ForClusterRole(test.clusterRole)

	// Verify labels are correctly merged
	for key, expectedValue := range test.expectedLabels {
		if actualValue, exists := test.clusterRole.Labels[key]; !exists || actualValue != expectedValue {
			t.Errorf("Expected label %s=%s, but got %s", key, expectedValue, actualValue)
		}
	}
}

func TestPatcherForClusterRoleBinding(t *testing.T) {
	test := struct {
		name               string
		patcher            *Patcher
		clusterRoleBinding *rbacv1.ClusterRoleBinding
		expectedLabels     map[string]string
	}{
		name: "clusterrolebinding with global and component labels",
		patcher: NewPatcher().
			WithGlobalLabels(map[string]string{"env": "prod", "team": "platform"}).
			WithComponentLabels(labels.Set{"component": "clusterrolebinding", "tier": "rbac"}).
			WithComponentInfo("test-crb", "test-instance"),
		clusterRoleBinding: &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "test-crb",
				Labels: map[string]string{"existing": "label"},
			},
		},
		expectedLabels: map[string]string{
			"existing":                     "label",
			"env":                          "prod",
			"team":                         "platform",
			"component":                    "clusterrolebinding",
			"tier":                         "rbac",
			"app.kubernetes.io/managed-by": "karmada-operator",
			"app.kubernetes.io/name":       "test-crb",
			"app.kubernetes.io/instance":   "test-instance",
		},
	}

	test.patcher.ForClusterRoleBinding(test.clusterRoleBinding)

	// Verify labels are correctly merged
	for key, expectedValue := range test.expectedLabels {
		if actualValue, exists := test.clusterRoleBinding.Labels[key]; !exists || actualValue != expectedValue {
			t.Errorf("Expected label %s=%s, but got %s", key, expectedValue, actualValue)
		}
	}
}

func TestPatcherForRole(t *testing.T) {
	test := struct {
		name           string
		patcher        *Patcher
		role           *rbacv1.Role
		expectedLabels map[string]string
	}{
		name: "role with global and component labels",
		patcher: NewPatcher().
			WithGlobalLabels(map[string]string{"env": "prod", "team": "platform"}).
			WithComponentLabels(labels.Set{"component": "role", "tier": "rbac"}).
			WithComponentInfo("test-role", "test-instance"),
		role: &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-role",
				Namespace: "default",
				Labels:    map[string]string{"existing": "label"},
			},
		},
		expectedLabels: map[string]string{
			"existing":                     "label",
			"env":                          "prod",
			"team":                         "platform",
			"component":                    "role",
			"tier":                         "rbac",
			"app.kubernetes.io/managed-by": "karmada-operator",
			"app.kubernetes.io/name":       "test-role",
			"app.kubernetes.io/instance":   "test-instance",
		},
	}

	test.patcher.ForRole(test.role)

	// Verify labels are correctly merged
	for key, expectedValue := range test.expectedLabels {
		if actualValue, exists := test.role.Labels[key]; !exists || actualValue != expectedValue {
			t.Errorf("Expected label %s=%s, but got %s", key, expectedValue, actualValue)
		}
	}
}

func TestPatcherForRoleBinding(t *testing.T) {
	test := struct {
		name           string
		patcher        *Patcher
		roleBinding    *rbacv1.RoleBinding
		expectedLabels map[string]string
	}{
		name: "rolebinding with global and component labels",
		patcher: NewPatcher().
			WithGlobalLabels(map[string]string{"env": "prod", "team": "platform"}).
			WithComponentLabels(labels.Set{"component": "rolebinding", "tier": "rbac"}).
			WithComponentInfo("test-rb", "test-instance"),
		roleBinding: &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-rb",
				Namespace: "default",
				Labels:    map[string]string{"existing": "label"},
			},
		},
		expectedLabels: map[string]string{
			"existing":                     "label",
			"env":                          "prod",
			"team":                         "platform",
			"component":                    "rolebinding",
			"tier":                         "rbac",
			"app.kubernetes.io/managed-by": "karmada-operator",
			"app.kubernetes.io/name":       "test-rb",
			"app.kubernetes.io/instance":   "test-instance",
		},
	}

	test.patcher.ForRoleBinding(test.roleBinding)

	// Verify labels are correctly merged
	for key, expectedValue := range test.expectedLabels {
		if actualValue, exists := test.roleBinding.Labels[key]; !exists || actualValue != expectedValue {
			t.Errorf("Expected label %s=%s, but got %s", key, expectedValue, actualValue)
		}
	}
}
