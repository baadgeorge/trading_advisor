package telegram

import (
	"reflect"
	"strconv"
	"strings"
)

func ConvId(id string) (uint32, error) {
	convId, err := strconv.ParseUint(id, 10, 32)
	return uint32(convId), err
}

func StrategyParamParser(param string) []string {
	return strings.Fields(param)
}

type field struct {
	field_name string
	field_type string
}

func GetStrategyFieldsFromStruct(strat interface{}) []field {
	var Fields []field
	val := reflect.ValueOf(strat)
	structType := val.Type()
	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Tag.Get("reflect") != "-" {
			Fields = append(Fields, field{
				field_name: structType.Field(i).Name,
				field_type: structType.Field(i).Type.String(),
			})
		}
	}
	return Fields
}
