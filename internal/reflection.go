package internal

import (
	"fmt"
	"reflect"
	"strings"
)

func HandleReflection(argStruct reflect.Type) ([]arg, []option, error) {
	argsConf := []arg{}
	optionsConf := []option{}
	for i := range argStruct.NumField() {
		field := argStruct.Field(i)

		fieldName := pascalToKebabCase(field.Name)
		if name := field.Tag.Get("name"); name != "" {
			fieldName = name
		}
		if field.Tag.Get("type") == "" || field.Tag.Get("type") == "arg" {
			argsConf = append(argsConf, arg{
				Name: fieldName,
				Type: field.Type.Kind().String(),
				Help: field.Tag.Get("help"),
			})
		}
		if field.Tag.Get("type") == "option" {
			if field.Name == "help" {
				return nil, nil, fmt.Errorf("Error in field `%s`: Field name cannot be `help` as this is reserved for the `--help` option.", field.Name)
			}
			if field.Tag.Get("short") == "h" {
				return nil, nil, fmt.Errorf("Error in field `%s`: Field short cannot be `h` as this is reserved for the `--help` option.", field.Name)
			}
			fieldValue := pascalToKebabCase(field.Name)
			if v, ok := field.Tag.Lookup("value"); ok {
				fieldValue = v
			}
			if field.Type.Kind().String() == "bool" {
				fieldValue = ""
			}

			optionsConf = append(optionsConf, option{
				Name:  fieldName,
				Type:  field.Type.Kind().String(),
				Value: fieldValue,
				Help:  field.Tag.Get("help"),
				Short: field.Tag.Get("short"),
			})
		}
	}

	return argsConf, optionsConf, nil
}

func pascalToKebabCase(pascal string) string {
	kebab := strings.ToLower(pascal[0:1])
	for _, c := range pascal[1:] {
		if strings.ToUpper(string(c)) == string(c) {
			kebab += "-"
		}
		kebab += strings.ToLower(string(c))
	}
	return kebab
}
