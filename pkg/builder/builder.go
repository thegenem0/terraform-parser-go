package builder

import (
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

type TreeBuilder struct {
	reflectorService reflector.IReflectorService
	tree             TreeData
}

func NewTreeBuilder(reflectorService reflector.IReflectorService) *TreeBuilder {
	return &TreeBuilder{
		reflectorService: reflectorService,
		tree: TreeData{
			Nodes: make([]PlanNodeData, 0),
		},
	}
}

func (tb *TreeBuilder) GetNodes() []PlanNodeData {
	return tb.tree.Nodes
}

func (tb *TreeBuilder) GetTree() TreeData {
	return tb.tree
}

func (tb *TreeBuilder) BuildTree(rootModule *tfjson.StateModule) {
	var createNode func(*tfjson.StateModule, string, bool) PlanNodeData

	createNode = func(mod *tfjson.StateModule, parentPath string, isRoot bool) PlanNodeData {
		nodeInfo := getNodeInfo(mod, parentPath, isRoot)
		node := PlanNodeData{
			ID:        nodeInfo.ID,
			Label:     nodeInfo.Label,
			Variables: nil,
			Children:  make([]PlanNodeData, 0),
			Changes:   nil,
		}

		for _, res := range mod.Resources {
			vars := tb.reflectorService.HandleVars(res.AttributeValues, res.Address)

			childNode := PlanNodeData{
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
	tb.addNode(topNode)
}

func (tb *TreeBuilder) addNode(node PlanNodeData) {
	tb.tree.Nodes = append(tb.tree.Nodes, node)
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
