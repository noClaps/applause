package parser

import (
	"fmt"
	"reflect"

	"github.com/noclaps/applause/internal/utils"
)

func (p *Parser) reflection() error {
	config := p.Config.Elem() // get struct value from pointer
	configType := config.Type()

	positionalsConf := []positional{}
	optionsConf := []option{}
	commandsConf := []command{}

	for i := range config.NumField() {
		field := configType.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := utils.PascalToKebabCase(field.Name)
		if name, ok := field.Tag.Lookup("name"); ok {
			fieldName = name
		}

		if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct {
			commandsConf = append(commandsConf, command{
				StructName:     field.Name,
				Name:           fieldName,
				Value:          config.Field(i),
				Help:           field.Tag.Get("help"),
				AllowEmptyArgs: true,
			})
			continue
		}
		if field.Type.Kind() == reflect.Struct || (field.Tag.Get("type") == "command" && field.Type.Kind() == reflect.Bool) {
			commandsConf = append(commandsConf, command{
				StructName: field.Name,
				Name:       fieldName,
				Value:      config.Field(i).Addr(),
				Help:       field.Tag.Get("help"),
			})
			continue
		}

		if field.Tag.Get("type") == "arg" || field.Tag.Get("type") == "" {
			positionalsConf = append(positionalsConf, positional{
				StructName: field.Name,
				Name:       fieldName,
				Type:       field.Type,
				Help:       field.Tag.Get("help"),
			})
			continue
		}

		if field.Tag.Get("type") == "option" {
			if fieldName == "help" {
				return fmt.Errorf("Error in field `%s`: Field name cannot be `Help` as this is reserved for the `--help` option.", field.Name)
			}
			if field.Tag.Get("short") == "h" {
				return fmt.Errorf("Error in field `%s`: Field short cannot be `h` as this is reserved for the `--help` option.", field.Name)
			}

			fieldValue := utils.PascalToKebabCase(field.Name)
			if v, ok := field.Tag.Lookup("value"); ok {
				fieldValue = v
			}
			if field.Type.Kind() == reflect.Bool {
				fieldValue = ""
			}

			defaultVal := config.Field(i)

			optionsConf = append(optionsConf, option{
				StructName: field.Name,
				Name:       fieldName,
				Type:       field.Type,
				Value:      fieldValue,
				Help:       field.Tag.Get("help"),
				Short:      field.Tag.Get("short"),
				Default:    defaultVal,
			})
		}
	}

	p.Positionals = positionalsConf
	p.Options = optionsConf
	p.Commands = commandsConf
	return nil
}
