package changes

import tfjson "github.com/hashicorp/terraform-json"

type ResourceKey struct {
	ModKey string
	Keys   []string
}

type ChangeItem struct {
	Actions         tfjson.Actions
	Address         string
	PreviousAddress string
	BeforeValue     interface{}
	AfterValue      interface{}
}

type Change struct {
	ModKey  string
	Changes []ChangeItem
}
