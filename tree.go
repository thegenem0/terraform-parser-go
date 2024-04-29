package main

import (
	"reflect"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/changes"
	"github.com/thegenem0/terraspect_server/pkg/depsgraph"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

type NodeInfo struct {
	ID       string
	Label    string
	FullPath string
}

func BuildTree(graph *depsgraph.DepsGraph, rootModule *tfjson.StateModule, changeService *changes.ChangeService, changes []*tfjson.ResourceChange) {
	var createNode func(*tfjson.StateModule, string, bool) depsgraph.PlanNodeData
	createNode = func(mod *tfjson.StateModule, parentPath string, isRoot bool) depsgraph.PlanNodeData {
		nodeInfo := getNodeInfo(mod, parentPath, isRoot)
		node := depsgraph.PlanNodeData{
			ID:        nodeInfo.ID,
			Label:     nodeInfo.Label,
			Variables: nil,
			Children:  make([]depsgraph.PlanNodeData, 0),
			Changes:   nil,
		}

		for _, res := range mod.Resources {
			vars := reflector.HandleVars(res.AttributeValues, changeService, res.Address)

			childNode := depsgraph.PlanNodeData{
				ID:        res.Address,
				Label:     res.Name,
				Variables: &vars,
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

func getValueType(value interface{}) reflect.Value {
	return reflect.ValueOf(value)
}
