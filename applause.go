package applause

import (
	"fmt"
	"reflect"

	"github.com/noclaps/applause/internal"
)

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
func Parse(args any) error {
	rv := reflect.ValueOf(args)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("Input value should be a pointer to a struct, received: %v", rv.Kind().String())
	}
	argStruct := rv.Elem().Type()
	parsedVals, err := internal.Parse(argStruct)
	if err != nil {
		return err
	}
	for k, v := range parsedVals {
		for f := range rv.Elem().NumField() {
			field := rv.Elem().Type().Field(f)
			fieldName := internal.PascalToKebabCase(field.Name)
			if fn, ok := field.Tag.Lookup("name"); ok {
				fieldName = fn
			}
			if k == fieldName {
				rv.Elem().Field(f).Set(reflect.ValueOf(v))
			}
		}
	}
	return nil
}
