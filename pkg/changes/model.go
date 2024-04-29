package changes

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

type ResourceKey struct {
	ModKey string
	Keys   []string
}

type ChangeItem struct {
	Actions         tfjson.Actions
	Address         string
	PreviousAddress string
	Changes         reflector.ChangeData
}

type Change struct {
	ModKey  string
	Changes []ChangeItem
}
