package changes

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/thegenem0/terraspect_server/pkg/reflector"
)

type IChangeService interface {
	GetResourceKeys() []ResourceKey
	GetChanges() []Change
	BuildChanges(changeData []*tfjson.ResourceChange)
	IsValidKey(modKey string, key string) bool
}

type ChangeService struct {
	reflectorService  reflector.IReflectorService
	validResourceKeys []ResourceKey
	changes           []Change
}

func NewChangeService(reflectorService reflector.IReflectorService) *ChangeService {
	return &ChangeService{
		reflectorService:  reflectorService,
		validResourceKeys: make([]ResourceKey, 0),
		changes:           make([]Change, 0),
	}
}

func (cs *ChangeService) GetResourceKeys() []ResourceKey {
	return cs.validResourceKeys
}

func (cs *ChangeService) IsValidKey(modKey string, key string) bool {
	for _, k := range cs.validResourceKeys {
		if k.ModKey == modKey {
			for _, v := range k.Keys {
				if v == key {
					return true
				}
			}
		}
	}
	return false
}

func (cs *ChangeService) GetChanges() []Change {
	return cs.changes
}

func (cs *ChangeService) BuildChanges(changeData []*tfjson.ResourceChange) {
	for _, change := range changeData {
		if change.Change.Actions.NoOp() {
			continue
		} else {
			cs.addChangeResource(
				change.Change.Actions,
				change.Address,
				change.PreviousAddress,
				change.Change.Before,
				change.Change.After,
			)
		}
	}
}

func (cs *ChangeService) addChangeResource(actions tfjson.Actions, address string, previousAddress string, before interface{}, after interface{}) {
	change := ChangeItem{
		Actions:         actions,
		Address:         address,
		PreviousAddress: previousAddress,
		BeforeValue:     before,
		AfterValue:      after,
	}

	cs.changes = append(cs.changes, Change{
		ModKey:  address,
		Changes: []ChangeItem{change},
	})
}
