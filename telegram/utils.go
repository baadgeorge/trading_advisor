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
	var FieldsType string
	val := reflect.ValueOf(strat)
	structType := val.Type()
	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Tag.Get("reflect") != "-" {
			switch structType.Field(i).Type.String() {
			case "float64", "float32":
				FieldsType = "вещественное"
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64":
				FieldsType = "целое"
			}
			Fields = append(Fields, field{
				field_name: structType.Field(i).Name,
				field_type: FieldsType,
			})
		}
	}
	return Fields
}
