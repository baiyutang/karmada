/*
Copyright 2023 The Karmada Authors.

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

package util

import (
	"reflect"
	"testing"

	"github.com/karmada-io/karmada/operator/pkg/constants"
	"k8s.io/apimachinery/pkg/labels"
)

func TestMergeLabels(t *testing.T) {
	tests := []struct {
		name            string
		baseLabels      map[string]string
		globalLabels    map[string]string
		componentLabels map[string]string
		systemLabels    map[string]string
		expected        map[string]string
	}{
		{
			name: "MergeAllLabelTypes_WithPriority",
			baseLabels: map[string]string{
				"base": "base-value",
			},
			globalLabels: map[string]string{
				"global": "global-value",
				"base":   "global-override",
			},
			componentLabels: map[string]string{
				"component": "component-value",
				"global":    "component-override",
			},
			systemLabels: map[string]string{
				"system": "system-value",
				"base":   "system-override",
			},
			expected: map[string]string{
				"base":      "system-override",    // system has highest priority
				"global":    "component-override", // component overrides global
				"component": "component-value",    // component specific
				"system":    "system-value",       // system specific
			},
		},
		{
			name: "MergeWithEmptyMaps",
			baseLabels: map[string]string{
				"base": "base-value",
			},
			globalLabels:    nil,
			componentLabels: nil,
			systemLabels:    nil,
			expected: map[string]string{
				"base": "base-value",
			},
		},
		{
			name:            "AllEmptyMaps",
			baseLabels:      nil,
			globalLabels:    nil,
			componentLabels: nil,
			systemLabels:    nil,
			expected:        map[string]string{},
		},
		{
			name:       "GlobalLabelsOnly",
			baseLabels: nil,
			globalLabels: map[string]string{
				"environment": "production",
				"team":        "platform",
			},
			componentLabels: nil,
			systemLabels:    nil,
			expected: map[string]string{
				"environment": "production",
				"team":        "platform",
			},
		},
		{
			name:         "ComponentLabelsOnly",
			baseLabels:   nil,
			globalLabels: nil,
			componentLabels: map[string]string{
				"component": "karmada-controller-manager",
				"version":   "v1.0.0",
			},
			systemLabels: nil,
			expected: map[string]string{
				"component": "karmada-controller-manager",
				"version":   "v1.0.0",
			},
		},
		{
			name:            "SystemLabelsOnly",
			baseLabels:      nil,
			globalLabels:    nil,
			componentLabels: nil,
			systemLabels: map[string]string{
				"app.kubernetes.io/managed-by": "karmada-operator",
				"app.kubernetes.io/name":       "test-component",
			},
			expected: map[string]string{
				"app.kubernetes.io/managed-by": "karmada-operator",
				"app.kubernetes.io/name":       "test-component",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MergeLabels(test.baseLabels, test.globalLabels, test.componentLabels, test.systemLabels)

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestGetSystemLabels(t *testing.T) {
	expected := map[string]string{
		constants.KarmadaOperatorLabelKeyName: constants.KarmadaOperator,
	}

	result := GetSystemLabels()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, but got %v", expected, result)
	}
}

func TestGetComponentSystemLabels(t *testing.T) {
	tests := []struct {
		name                string
		componentName       string
		karmadaInstanceName string
		expected            map[string]string
	}{
		{
			name:                "ValidComponentAndInstance",
			componentName:       "karmada-controller-manager",
			karmadaInstanceName: "my-karmada",
			expected: map[string]string{
				constants.KarmadaOperatorLabelKeyName: constants.KarmadaOperator,
				constants.AppNameLabel:                "karmada-controller-manager",
				constants.AppInstanceLabel:            "my-karmada",
			},
		},
		{
			name:                "EmptyStrings",
			componentName:       "",
			karmadaInstanceName: "",
			expected: map[string]string{
				constants.KarmadaOperatorLabelKeyName: constants.KarmadaOperator,
				constants.AppNameLabel:                "",
				constants.AppInstanceLabel:            "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GetComponentSystemLabels(test.componentName, test.karmadaInstanceName)

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestMergeLabelsForComponent(t *testing.T) {
	tests := []struct {
		name                string
		baseLabels          map[string]string
		globalLabels        map[string]string
		componentLabels     map[string]string
		componentName       string
		karmadaInstanceName string
		expected            map[string]string
	}{
		{
			name: "CompleteMerge",
			baseLabels: map[string]string{
				"base": "base-value",
			},
			globalLabels: map[string]string{
				"environment": "production",
			},
			componentLabels: map[string]string{
				"component": "karmada-scheduler",
			},
			componentName:       "karmada-scheduler",
			karmadaInstanceName: "prod-karmada",
			expected: map[string]string{
				"base":                                "base-value",
				"environment":                         "production",
				"component":                           "karmada-scheduler",
				constants.KarmadaOperatorLabelKeyName: constants.KarmadaOperator,
				constants.AppNameLabel:                "karmada-scheduler",
				constants.AppInstanceLabel:            "prod-karmada",
			},
		},
		{
			name: "WithLabelConflicts",
			baseLabels: map[string]string{
				"label": "base-value",
			},
			globalLabels: map[string]string{
				"label": "global-value",
			},
			componentLabels: map[string]string{
				"label": "component-value",
			},
			componentName:       "karmada-webhook",
			karmadaInstanceName: "test-karmada",
			expected: map[string]string{
				"label":                               "component-value", // component overrides global
				constants.KarmadaOperatorLabelKeyName: constants.KarmadaOperator,
				constants.AppNameLabel:                "karmada-webhook",
				constants.AppInstanceLabel:            "test-karmada",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MergeLabelsForComponent(test.baseLabels, test.globalLabels, test.componentLabels, test.componentName, test.karmadaInstanceName)

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestMergeLabelsSet(t *testing.T) {
	tests := []struct {
		name            string
		baseLabels      map[string]string
		globalLabels    map[string]string
		componentLabels map[string]string
		systemLabels    map[string]string
		expected        map[string]string
	}{
		{
			name: "MergeLabelsSet_WithPriority",
			baseLabels: map[string]string{
				"base": "base-value",
			},
			globalLabels: map[string]string{
				"global": "global-value",
			},
			componentLabels: map[string]string{
				"component": "component-value",
			},
			systemLabels: map[string]string{
				"system": "system-value",
			},
			expected: map[string]string{
				"base":      "base-value",
				"global":    "global-value",
				"component": "component-value",
				"system":    "system-value",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Convert to labels.Set
			baseSet := labels.Set(test.baseLabels)
			globalSet := labels.Set(test.globalLabels)
			componentSet := labels.Set(test.componentLabels)
			systemSet := labels.Set(test.systemLabels)

			result := MergeLabelsSet(baseSet, globalSet, componentSet, systemSet)

			// Convert result back to map for comparison
			resultMap := map[string]string(result)

			if !reflect.DeepEqual(resultMap, test.expected) {
				t.Errorf("expected %v, but got %v", test.expected, resultMap)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkMergeLabels(b *testing.B) {
	baseLabels := map[string]string{
		"base1": "value1", "base2": "value2", "base3": "value3",
	}
	globalLabels := map[string]string{
		"global1": "value1", "global2": "value2", "global3": "value3",
	}
	componentLabels := map[string]string{
		"component1": "value1", "component2": "value2", "component3": "value3",
	}
	systemLabels := map[string]string{
		"system1": "value1", "system2": "value2", "system3": "value3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MergeLabels(baseLabels, globalLabels, componentLabels, systemLabels)
	}
}

func BenchmarkGetComponentSystemLabels(b *testing.B) {
	componentName := "karmada-controller-manager"
	instanceName := "test-karmada"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetComponentSystemLabels(componentName, instanceName)
	}
}
