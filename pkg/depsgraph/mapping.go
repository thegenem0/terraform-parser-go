package depsgraph

import "github.com/hashicorp/terraform-config-inspect/tfconfig"

type Resource struct {
	Type TFResourceType
	Name string
	Line *int

	Children map[string]*Resource

	ChangeAction TFActionType

	Required     *bool
	Sensitive    bool
	Provider     string
	ResourceType string

	Source  string
	Version string
}

type ModuleCall struct {
	tfconfig.ModuleCall
}
