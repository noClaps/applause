package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func ValToType(input string, returnType reflect.Type) (reflect.Value, error) {
	switch returnType.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(input)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Error parsing %s: %v", returnType, err)
		}
		return reflect.ValueOf(b), nil
	case reflect.Float32, reflect.Float64:
		bitSize := 64
		if returnType.Kind() == reflect.Float32 {
			bitSize = 32
		}
		f, err := strconv.ParseFloat(input, bitSize)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Error parsing %s: %v", returnType, err)
		}
		return reflect.ValueOf(f).Convert(returnType), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize := 64
		switch returnType.Kind() {
		case reflect.Int8:
			bitSize = 8
		case reflect.Int16:
			bitSize = 16
		case reflect.Int32:
			bitSize = 32
		}
		i, err := strconv.ParseInt(input, 10, bitSize)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Error parsing %s: %v", returnType, err)
		}
		return reflect.ValueOf(i).Convert(returnType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bitSize := 64
		switch returnType.Kind() {
		case reflect.Uint8:
			bitSize = 8
		case reflect.Uint16:
			bitSize = 16
		case reflect.Uint32:
			bitSize = 32
		}
		u, err := strconv.ParseUint(input, 10, bitSize)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Error parsing %s: %v", returnType, err)
		}
		return reflect.ValueOf(u).Convert(returnType), nil
	case reflect.Complex64, reflect.Complex128:
		bitSize := 128
		if returnType.Kind() == reflect.Float64 {
			bitSize = 64
		}
		c, err := strconv.ParseComplex(input, bitSize)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Error parsing %s: %v", returnType, err)
		}
		return reflect.ValueOf(c).Convert(returnType), nil
	case reflect.String:
		return reflect.ValueOf(input), nil
	}
	return reflect.Value{}, fmt.Errorf("Type `%s` is unsupported, please use a supported type.", returnType)
}
