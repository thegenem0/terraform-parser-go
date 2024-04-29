package depsgraph

type DepsGraph struct {
	Nodes []PlanNodeData
}

func NewDepsGraph() *DepsGraph {
	return &DepsGraph{}
}

func (dg *DepsGraph) AddNode(node PlanNodeData) {
	dg.Nodes = append(dg.Nodes, node)
}

func (dg *DepsGraph) GetNodes() []PlanNodeData {
	return dg.Nodes
}
