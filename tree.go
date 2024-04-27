package main

import (
	"reflect"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/depsgraph"
)

type NodeInfo struct {
	ID       string
	Label    string
	FullPath string
}

func BuildTree(graph *depsgraph.DepsGraph, rootModule *tfjson.StateModule) {

	var createNode func(*tfjson.StateModule, string, bool) depsgraph.PlanNodeData

	createNode = func(mod *tfjson.StateModule, parentPath string, isRoot bool) depsgraph.PlanNodeData {
		nodeInfo := getNodeInfo(mod, parentPath, isRoot)

		node := depsgraph.PlanNodeData{
			ID:        nodeInfo.ID,
			Label:     nodeInfo.Label,
			Variables: nil,
			Children:  make([]depsgraph.PlanNodeData, 0),
		}

		for _, res := range mod.Resources {
			childNode := depsgraph.PlanNodeData{
				ID:        res.Address,
				Label:     res.Name,
				Variables: convertToVariable(res.AttributeValues),
			}
			node.Children = append(node.Children, childNode)
		}

		for _, childMod := range mod.ChildModules {
			childPath := nodeInfo.FullPath
			childNode := createNode(childMod, childPath, false)
			node.Children = append(node.Children, childNode)
		}

		return node
	}

	topNode := createNode(rootModule, "", true)

	graph.AddNode(topNode)
}

func getNodeInfo(mod *tfjson.StateModule, parentPath string, isRoot bool) NodeInfo {
	var id, label, fullPath string

	if isRoot {
		id = "root"
		label = "Root Node"
		fullPath = parentPath
	} else {
		fullPath = parentPath
		if parentPath != "" {
			fullPath += "."
		}
		fullPath += parseModulePath(mod.Address)

		id = mod.Address
		label = fullPath
	}

	return NodeInfo{
		ID:       id,
		Label:    label,
		FullPath: fullPath,
	}
}

// function to determine level of nesting in the module path
func parseModulePath(path string) string {
	parts := strings.Split(path, ".")
	var components []string

	for _, part := range parts {
		if idx := strings.Index(part, "["); idx != -1 {
			part = part[:idx]
		}

		if part != "module" {
			components = append(components, part)
		}
	}

	return strings.Join(components, ".")
}

// Generate Key-Value pairs from the resource attributes.
func convertToVariable(varMap map[string]interface{}) depsgraph.Variable {
	var vars depsgraph.Variable
	for key, value := range varMap {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			convertedValue := convertToVariable(nestedMap)
			if len(convertedValue) == 0 {
				continue
			}
			value = convertedValue
		}
		if isEmptyValue(value) || isDefaultValue(value) {
			continue
		}
		vars = append(vars, depsgraph.KeyValue{
			Key:   key,
			Value: value,
		})
	}
	return vars
}

// Filter out empty values from the variable map.
func isEmptyValue(v interface{}) bool {
	switch value := v.(type) {
	case string:
		return value == ""
	case map[string]interface{}:
		return len(value) == 0
	case []interface{}:
		return len(value) == 0
	case nil:
		return true
	default:
		return false
	}
}

// Filter out default values from the variable map.
func isDefaultValue(v interface{}) bool {
	switch value := v.(type) {
	case bool:
		return !value
	case string:
		return value == ""
	case int, int32, int64, float32, float64:
		return reflect.ValueOf(value).Float() == 0
	case []interface{}:
		return len(value) == 0
	case map[string]interface{}:
		return len(value) == 0
	case nil:
		return true
	default:
		return false
	}
}
