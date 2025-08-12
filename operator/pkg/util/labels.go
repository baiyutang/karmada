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
	"k8s.io/apimachinery/pkg/labels"

	"github.com/karmada-io/karmada/operator/pkg/constants"
)

// MergeLabels merges multiple label maps with priority order.
// Labels from higher priority maps will override labels from lower priority maps.
// The order of precedence is:
// 1. systemLabels (highest priority - these are required system labels)
// 2. componentLabels (component-specific labels)
// 3. globalLabels (global labels from Karmada spec)
// 4. baseLabels (lowest priority - labels already on the resource)
func MergeLabels(baseLabels, globalLabels, componentLabels, systemLabels map[string]string) map[string]string {
	// Start with base labels (lowest priority)
	result := make(map[string]string)
	for k, v := range baseLabels {
		result[k] = v
	}

	// Merge global labels (override base labels)
	for k, v := range globalLabels {
		result[k] = v
	}

	// Merge component labels (override global and base labels)
	for k, v := range componentLabels {
		result[k] = v
	}

	// Merge system labels (highest priority - override all others)
	for k, v := range systemLabels {
		result[k] = v
	}

	return result
}

// GetSystemLabels returns the standard system labels that should be applied to all Karmada resources.
// These labels are used to identify resources managed by the Karmada operator.
func GetSystemLabels() map[string]string {
	return map[string]string{
		constants.KarmadaOperatorLabelKeyName: constants.KarmadaOperator,
	}
}

// GetComponentSystemLabels returns system labels specific to a component.
// It includes the base system labels plus component-specific identification labels.
func GetComponentSystemLabels(componentName, karmadaInstanceName string) map[string]string {
	systemLabels := GetSystemLabels()
	systemLabels[constants.AppNameLabel] = componentName
	systemLabels[constants.AppInstanceLabel] = karmadaInstanceName
	return systemLabels
}

// MergeLabelsForComponent is a convenience function that merges labels for a specific component.
// It combines global labels, component labels, and system labels following the correct precedence.
func MergeLabelsForComponent(baseLabels, globalLabels, componentLabels map[string]string, componentName, karmadaInstanceName string) map[string]string {
	systemLabels := GetComponentSystemLabels(componentName, karmadaInstanceName)
	return MergeLabels(baseLabels, globalLabels, componentLabels, systemLabels)
}

// MergeLabelsSet merges multiple labels.Set with priority order, returning a labels.Set.
// This is a wrapper around MergeLabels that works with the Kubernetes labels.Set type.
func MergeLabelsSet(baseLabels, globalLabels, componentLabels, systemLabels labels.Set) labels.Set {
	merged := MergeLabels(
		map[string]string(baseLabels),
		map[string]string(globalLabels),
		map[string]string(componentLabels),
		map[string]string(systemLabels),
	)
	return labels.Set(merged)
}
