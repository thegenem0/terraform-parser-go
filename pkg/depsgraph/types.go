package depsgraph

type PlanNodeData struct {
	ID        string         `json:"id,omitempty"`
	Label     string         `json:"label,omitempty"`
	Type      TFResourceType `json:"type,omitempty"`
	Variables Variable       `json:"variables,omitempty"`
	Children  []PlanNodeData `json:"children,omitempty"`
}

type Variable []KeyValue

type KeyValue struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type PlanEdgeData struct {
	ID     string
	Source string
	Target string
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
