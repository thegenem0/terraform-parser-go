package reflector

import (
	"reflect"

	"github.com/thegenem0/terraspect_server/pkg/changes"
)

type SimpleKVPair struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
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
	Values *ComplexKVPair `json:"values"`
}

func HandleVars(variables map[string]interface{}, changeService *changes.ChangeService, modKey string) VariableData {
	var simpleValues []*SimpleKVPair
	var complexValues []*ComplexKVPair

	for key, value := range variables {
		if !isEmptyValue(value) && !isDefaultValue(value) {
			if getValueType(value) == reflect.Slice || getValueType(value) == reflect.Map {
				complexValue := HandleComplexValue(key, value)
				if complexValue != nil {
					changeService.AddResourceKey(modKey, complexValue.Key)
					complexValues = append(complexValues, complexValue)
				}
			} else {
				simpleValue := HandleSimpleValue(key, value)
				if simpleValue != nil {
					changeService.AddResourceKey(modKey, simpleValue.Key)
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

func HandleChanges(changes interface{}) ChangeData {
	if !isEmptyValue(changes) && !isDefaultValue(changes) {
		return ChangeData{
			Values: HandleComplexValue("changes", changes),
		}
	}
	return ChangeData{}
}

func HandleSimpleValue(key string, value interface{}) *SimpleKVPair {
	if isEmptyValue(value) || isDefaultValue(value) {
		return nil
	}

	return &SimpleKVPair{
		Key:   key,
		Value: value,
	}
}

func HandleComplexValue(key string, value interface{}) *ComplexKVPair {
	var simpleValues []*SimpleKVPair

	if isEmptyValue(value) || isDefaultValue(value) {
		return nil
	}

	if getValueType(value) == reflect.Slice {
		for _, v := range value.([]interface{}) {
			if getValueType(v) == reflect.Slice {
				HandleComplexValue(key, v)
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
				HandleComplexValue(k, v)
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
