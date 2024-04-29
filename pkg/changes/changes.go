package changes

type IChangeService interface {
	GetResourceKeys() []ResourceKey
	AddResourceKey(modKey string, key string)
	IsValidKey(modKey string, key string) bool
}

type ChangeService struct {
	validResourceKeys []ResourceKey
}

type ResourceKey struct {
	ModKey string
	Keys   []string
}

func NewChangeService() *ChangeService {
	return &ChangeService{
		validResourceKeys: make([]ResourceKey, 0),
	}
}

func (cs *ChangeService) GetResourceKeys() []ResourceKey {
	return cs.validResourceKeys
}

func (cs *ChangeService) AddResourceKey(modKey string, key string) {
	for i, k := range cs.validResourceKeys {
		if k.ModKey == modKey {
			cs.validResourceKeys[i].Keys = append(cs.validResourceKeys[i].Keys, key)
			return
		}
	}
	cs.validResourceKeys = append(cs.validResourceKeys, ResourceKey{ModKey: modKey, Keys: []string{key}})
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
