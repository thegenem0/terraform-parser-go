package reflector

import (
	"reflect"
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

type IReflectorService interface {
	HandleVars(variables map[string]interface{}, modKey string) VariableData
	HandleChanges(changes interface{}) ChangeData
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

func (rs *ReflectorService) HandleChanges(changes interface{}) ChangeData {
	if !isEmptyValue(changes) && !isDefaultValue(changes) {
		return ChangeData{
			Values: handleComplexValue("changes", changes),
		}
	}
	return ChangeData{}
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
