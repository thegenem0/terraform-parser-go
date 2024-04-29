package depsgraph

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

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
	Before             *reflector.ChangeData `json:"omit"`
	After              *reflector.ChangeData `json:"omit"`
}

type TFResourceType string
type TFActionType string

const (
	TypeFile     TFResourceType = "file"
	TypeLocal    TFResourceType = "locals"
	TypeVariable TFResourceType = "variable"
	TypeOutput   TFResourceType = "output"
	TypeResource TFResourceType = "resource"
	TypeData     TFResourceType = "data"
	TypeModule   TFResourceType = "module"
	TypeDefault  string         = "unknown file"
)

const (
	TypeNoop    TFActionType = "no-op"
	TypeCreate  TFActionType = "create"
	TypeRead    TFActionType = "read"
	TypeUpdate  TFActionType = "update"
	TypeDelete  TFActionType = "delete"
	TypeReplace TFActionType = "replace"
)
