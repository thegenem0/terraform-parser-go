package main

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/depsgraph"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

func collectValidResources(rootModule *tfjson.StateModule) map[string]bool {
	validResources := make(map[string]bool)
	var collect func(mod *tfjson.StateModule)
	collect = func(mod *tfjson.StateModule) {
		for _, res := range mod.Resources {
			validResources[res.Address] = true
		}
		for _, childMod := range mod.ChildModules {
			collect(childMod)
		}
	}
	collect(rootModule)
	return validResources
}

func filterSignificantChanges(before, after *reflector.ChangeData) *reflector.ChangeData {
	if before == nil || after == nil {
		return nil // If either is nil, there's nothing to compare
	}

	significantChanges := &reflector.ChangeData{
		Values: &reflector.ComplexKVPair{
			Key:   "changes",
			Value: []*reflector.SimpleKVPair{},
		},
	}

	beforeValues := make(map[string]string)
	afterValues := make(map[string]string)

	// Convert slices to maps for easier comparison
	for _, kv := range before.Values.Value {
		if val, ok := kv.Value.(string); ok {
			beforeValues[kv.Key] = val
		}
	}
	for _, kv := range after.Values.Value {
		if val, ok := kv.Value.(string); ok {
			afterValues[kv.Key] = val
		}
	}

	// Check for changes
	for key, afterVal := range afterValues {
		if beforeVal, exists := beforeValues[key]; !exists || beforeVal != afterVal {
			significantChanges.Values.Value = append(significantChanges.Values.Value, &reflector.SimpleKVPair{Key: key, Value: afterVal})
		}
	}

	return significantChanges
}

func applyChangesFilter(changesMap map[string]*depsgraph.ResourceChanges) []*reflector.ChangeData {
	var significantChanges []*reflector.ChangeData
	for _, changes := range changesMap {
		if changes.HasChange {
			significantChanges = append(significantChanges, filterSignificantChanges(changes.Before, changes.After))
		}
	}
	return significantChanges
}

func createChangesMap(changes []*tfjson.ResourceChange) map[string]*depsgraph.ResourceChanges {
	changesMap := make(map[string]*depsgraph.ResourceChanges)
	for _, change := range changes {
		var beforeVars, afterVars reflector.ChangeData
		if change.Change.Before != nil {
			beforeVars = reflector.HandleChanges(change.Change.Before)
		}
		if change.Change.After != nil {
			afterVars = reflector.HandleChanges(change.Change.After)
		}
		resChange := &depsgraph.ResourceChanges{
			HasChange: true,
			Before:    &beforeVars,
			After:     &afterVars,
		}
		changesMap[change.Address] = resChange
	}
	return changesMap
}
