package builder

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

type TreeData struct {
	Nodes []PlanNodeData `json:"nodes,omitempty"`
}

type PlanNodeData struct {
	ID        string                  `json:"id,omitempty"`
	Label     string                  `json:"label,omitempty"`
	Type      TFResourceType          `json:"type,omitempty"`
	Variables *reflector.VariableData `json:"variables,omitempty"`
	Children  []PlanNodeData          `json:"children,omitempty"`
	Changes   *ResourceChanges
}

type ResourceChanges struct {
	HasChange          bool                  `json:"has_change"`
	ChangeOps          []*tfjson.Action      `json:"change_ops,omitempty"`
	SignificantChanges *reflector.ChangeData `json:"significant_changes,omitempty"`
	Before             *reflector.ChangeData `json:"before,omit"`
	After              *reflector.ChangeData `json:"after,omit"`
}

type NodeInfo struct {
	ID       string
	Label    string
	FullPath string
}

type TFResourceType string
type TFActionType string

const (
	TypeNoop    TFActionType = "no-op"
	TypeCreate  TFActionType = "create"
	TypeRead    TFActionType = "read"
	TypeUpdate  TFActionType = "update"
	TypeDelete  TFActionType = "delete"
	TypeReplace TFActionType = "replace"
)
