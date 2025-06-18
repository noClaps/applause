package applause

import (
	"fmt"
	"reflect"

	"github.com/noclaps/applause/internal"
)

// The help string for the command. This will only contain a value if
// `applause.Parse()` has been called first, otherwise it will be an empty
// string.
var Help string = ""

// The usage string for the command. This will only contain a value if
// `applause.Parse()` has been called first, otherwise it will be an empty
// string.
var Usage string = ""

// The input is a pointer to the args struct. Each field in the args struct
// should have some tags:
//
// - `type`: The type can be "arg" or "option". If omitted, the default is
// "arg". If any other type is provided, the field is ignored.
//
// - `name`: The name of the argument or option. If omitted, the default is the
// field name in kebab-case. If you'd like to have an option have a different
// name, you can write it as `name:"option-name"` in the tags.
//
// - `help`: The help text for the argument or option, will be displayed in the
// command help when the command is called with `--help` or `-h`.
//
// - `short`: Only applicable when `type` is "option". The short form of the
// option. For instance, if you have a field with the tag
// `name:"option" short:"o"`, you can call the command with `--option` or `-o`.
//
// - `value`: Only applicable when `type` is "option" and the field type is not
// "bool". The name of the option value to be displayed in the help text. For
// instance, `name:"option" value:"val"` will be displayed as `--option <val>`
// in the help text.
//
// All fields that you'd like to be parsed should be exported in the struct.
func Parse(args any) error {
	rv := reflect.ValueOf(args)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("Input value should be a pointer to a struct, received: %v", rv.Kind().String())
	}
	argStruct := rv.Elem().Type()
	argsConfig, optionsConfig, err := internal.HandleReflection(argStruct, rv.Elem())
	if err != nil {
		return fmt.Errorf("Error during reflection: %v", err)
	}
	Help = internal.GenerateHelp(argsConfig, optionsConfig)
	Usage = internal.GenerateUsage(argsConfig, optionsConfig)
	parsedVals, err := internal.Parse(argsConfig, optionsConfig)
	if err != nil {
		return fmt.Errorf("Error during parsing: %v", err)
	}

	for k, v := range parsedVals {
		for f := range rv.Elem().NumField() {
			field := rv.Elem().Type().Field(f)
			fieldName := internal.PascalToKebabCase(field.Name)
			if fn, ok := field.Tag.Lookup("name"); ok {
				fieldName = fn
			}
			if k == fieldName {
				fieldType := rv.Elem().Field(f).Type().String()
				if fieldType[0:2] != "[]" {
					rv.Elem().Field(f).Set(reflect.ValueOf(v))
					continue
				}

				switch fieldType[2:] {
				case "bool":
					v := v.([]any)
					slice := make([]bool, len(v))
					for i := range v {
						val, ok := v[i].(bool)
						if !ok {
							return fmt.Errorf("Error parsing bool slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "float32":
					v := v.([]any)
					slice := make([]float32, len(v))
					for i := range v {
						val, ok := v[i].(float32)
						if !ok {
							return fmt.Errorf("Error parsing float32 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "float64":
					v := v.([]any)
					slice := make([]float64, len(v))
					for i := range v {
						val, ok := v[i].(float64)
						if !ok {
							return fmt.Errorf("Error parsing float64 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "int":
					v := v.([]any)
					slice := make([]int, len(v))
					for i := range v {
						val, ok := v[i].(int)
						if !ok {
							return fmt.Errorf("Error parsing int slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "int8":
					v := v.([]any)
					slice := make([]int8, len(v))
					for i := range v {
						val, ok := v[i].(int8)
						if !ok {
							return fmt.Errorf("Error parsing int8 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "int16":
					v := v.([]any)
					slice := make([]int16, len(v))
					for i := range v {
						val, ok := v[i].(int16)
						if !ok {
							return fmt.Errorf("Error parsing int16 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "int32":
					v := v.([]any)
					slice := make([]int32, len(v))
					for i := range v {
						val, ok := v[i].(int32)
						if !ok {
							return fmt.Errorf("Error parsing int32 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "int64":
					v := v.([]any)
					slice := make([]int64, len(v))
					for i := range v {
						val, ok := v[i].(int64)
						if !ok {
							return fmt.Errorf("Error parsing int64 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "uint":
					v := v.([]any)
					slice := make([]uint, len(v))
					for i := range v {
						val, ok := v[i].(uint)
						if !ok {
							return fmt.Errorf("Error parsing uint slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "uint8":
					v := v.([]any)
					slice := make([]uint8, len(v))
					for i := range v {
						val, ok := v[i].(uint8)
						if !ok {
							return fmt.Errorf("Error parsing uint8 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "uint16":
					v := v.([]any)
					slice := make([]uint16, len(v))
					for i := range v {
						val, ok := v[i].(uint16)
						if !ok {
							return fmt.Errorf("Error parsing uint16 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "uint32":
					v := v.([]any)
					slice := make([]uint32, len(v))
					for i := range v {
						val, ok := v[i].(uint32)
						if !ok {
							return fmt.Errorf("Error parsing uint32 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "uint64":
					v := v.([]any)
					slice := make([]uint64, len(v))
					for i := range v {
						val, ok := v[i].(uint64)
						if !ok {
							return fmt.Errorf("Error parsing uint64 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "complex64":
					v := v.([]any)
					slice := make([]complex64, len(v))
					for i := range v {
						val, ok := v[i].(complex64)
						if !ok {
							return fmt.Errorf("Error parsing complex64 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "complex128":
					v := v.([]any)
					slice := make([]complex128, len(v))
					for i := range v {
						val, ok := v[i].(complex128)
						if !ok {
							return fmt.Errorf("Error parsing complex128 slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				case "string":
					v := v.([]any)
					slice := make([]string, len(v))
					for i := range v {
						val, ok := v[i].(string)
						if !ok {
							return fmt.Errorf("Error parsing string slice from `%v`", v)
						}
						slice[i] = val
					}
					rv.Elem().Field(f).Set(reflect.ValueOf(slice))
				default:
					return fmt.Errorf("Type `%s` is unsupported, please use a supported type.", fieldType)
				}

			}
		}
	}
	return nil
}
