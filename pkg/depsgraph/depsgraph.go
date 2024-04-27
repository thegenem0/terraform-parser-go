package depsgraph

type DepsGraph struct {
	Nodes []PlanNodeData
	Edges []PlanEdgeData
}

func NewDepsGraph() *DepsGraph {
	return &DepsGraph{}
}

func (dg *DepsGraph) AddEdge(edge PlanEdgeData) {
	dg.Edges = append(dg.Edges, edge)
}

func (dg *DepsGraph) AddNode(node PlanNodeData) {
	dg.Nodes = append(dg.Nodes, node)
}

func (dg *DepsGraph) GetNodes() []PlanNodeData {
	return dg.Nodes
}

func (dg *DepsGraph) GetEdges() []PlanEdgeData {
	return dg.Edges
}
