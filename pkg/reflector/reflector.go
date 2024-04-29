package reflector

import (
	"fmt"
	"reflect"
)

type SimpleKVPair struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	PrevValue interface{} `json:"prev_value,omitempty"`
}

type ComplexKVPair struct {
	Key   string          `json:"key"`
	Value []*SimpleKVPair `json:"value"`
}

type VariableData struct {
	SimpleValues  []*SimpleKVPair  `json:"simple_values"`
	ComplexValues []*ComplexKVPair `json:"complex_values"`
}

type ChangeData struct {
	Values *ComplexKVPair `json:"values,omitempty"`
}

type IReflectorService interface {
	HandleVars(variables map[string]interface{}, modKey string) VariableData
	HandleChanges(beforeChanges interface{}, afterChanges interface{}) ChangeData
}

type ReflectorService struct {
}

func NewReflectorService() *ReflectorService {
	return &ReflectorService{}
}

func (rs *ReflectorService) HandleVars(variables map[string]interface{}, modKey string) VariableData {
	var simpleValues []*SimpleKVPair
	var complexValues []*ComplexKVPair

	for key, value := range variables {
		if !isEmptyValue(value) && !isDefaultValue(value) {
			if getValueType(value) == reflect.Slice || getValueType(value) == reflect.Map {
				complexValue := handleComplexValue(key, value)
				if complexValue != nil {
					complexValues = append(complexValues, complexValue)
				}
			} else {
				simpleValue := handleSimpleValue(key, value)
				if simpleValue != nil {
					simpleValues = append(simpleValues, simpleValue)
				}
			}
		}
	}

	return VariableData{
		SimpleValues:  simpleValues,
		ComplexValues: complexValues,
	}
}

func (rs *ReflectorService) HandleChanges(beforeChanges interface{}, afterChanges interface{}) ChangeData {
	var beforeData, afterData ChangeData

	if !isEmptyValue(beforeChanges) && !isDefaultValue(beforeChanges) {
		beforeData = ChangeData{
			Values: handleComplexValue("before", beforeChanges),
		}
	}

	if !isEmptyValue(afterChanges) && !isDefaultValue(afterChanges) {
		afterData = ChangeData{
			Values: handleComplexValue("after", afterChanges),
		}
	}

	if beforeData.Values == nil && afterData.Values == nil {
		return ChangeData{}
	}

	changedVars := filterChangedVars(&beforeData, &afterData)
	for _, common := range changedVars {
		fmt.Println(common.Key, common.Value)
	}

	if len(changedVars) == 0 {
		return ChangeData{}
	}

	return ChangeData{
		Values: &ComplexKVPair{
			Key:   "changes",
			Value: changedVars,
		},
	}
}

func handleSimpleValue(key string, value interface{}) *SimpleKVPair {
	if isEmptyValue(value) || isDefaultValue(value) {
		return nil
	}

	return &SimpleKVPair{
		Key:   key,
		Value: value,
	}
}

func handleComplexValue(key string, value interface{}) *ComplexKVPair {
	var simpleValues []*SimpleKVPair

	if isEmptyValue(value) || isDefaultValue(value) {
		return nil
	}

	if getValueType(value) == reflect.Slice {
		for _, v := range value.([]interface{}) {
			if getValueType(v) == reflect.Slice {
				handleComplexValue(key, v)
			} else {
				simpleValues = append(simpleValues, &SimpleKVPair{
					Key:   "option",
					Value: v,
				})
			}
		}
	} else if getValueType(value) == reflect.Map {
		for k, v := range value.(map[string]interface{}) {
			if getValueType(v) == reflect.Map {
				handleComplexValue(k, v)
			} else {
				simpleValues = append(simpleValues, &SimpleKVPair{
					Key:   k,
					Value: v,
				})
			}
		}
	} else {
		simpleValues = append(simpleValues, &SimpleKVPair{
			Key:   key,
			Value: value,
		})
	}
	return &ComplexKVPair{
		Key:   key,
		Value: simpleValues,
	}
}

func getValueType(value interface{}) reflect.Kind {
	return reflect.ValueOf(value).Kind()
}

// Filter out empty values from the variable map.
func isEmptyValue(v interface{}) bool {
	switch value := v.(type) {
	case string:
		return value == ""
	case map[string]interface{}:
		return len(value) == 0
	case []interface{}:
		return len(value) == 0
	case nil:
		return true
	default:
		return false
	}
}

// Filter out default values from the variable map.
func isDefaultValue(v interface{}) bool {
	switch value := v.(type) {
	case bool:
		return !value
	case string:
		return value == ""
	case int, int32, int64, float32, float64:
		return reflect.ValueOf(value).Float() == 0
	case []interface{}:
		return len(value) == 0
	case map[string]interface{}:
		return len(value) == 0
	case nil:
		return true
	default:
		return false
	}
}

func sliceToMap(pairs []*SimpleKVPair) map[string]interface{} {
	m := make(map[string]interface{})
	for _, pair := range pairs {
		m[pair.Key] = pair.Value
	}
	return m
}

func filterChangedVars(before, after *ChangeData) []*SimpleKVPair {
	beforeMap := sliceToMap(before.Values.Value)
	var changedVars []*SimpleKVPair

	for _, afterItem := range after.Values.Value {
		beforeItem, exists := beforeMap[afterItem.Key]
		if exists && !reflect.DeepEqual(beforeItem, afterItem.Value) {
			changedVars = append(changedVars, &SimpleKVPair{
				Key:       afterItem.Key,
				Value:     afterItem.Value,
				PrevValue: beforeItem,
			})
		}
	}
	return changedVars
}
