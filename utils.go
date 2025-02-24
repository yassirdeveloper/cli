package main

import (
	"fmt"
	"strconv"
)

const (
	NoType ValueType = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeString
)

func ParseValue(valueType ValueType, value any) (interface{}, error) {
	valueString, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("cant cast to string")
	}
	switch valueType {
	case TypeInt:
		return strconv.Atoi(valueString) // Parse string to int
	case TypeFloat:
		return strconv.ParseFloat(valueString, 64) // Parse string to float64
	case TypeBool:
		return strconv.ParseBool(valueString) // Parse string to bool
	case TypeString:
		return valueString, nil // No parsing needed for strings
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}
