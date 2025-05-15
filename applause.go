package applause

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/noclaps/applause/internal"
)

func Parse(argStruct any) error {
	rv := reflect.ValueOf(argStruct)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("Input value needs to be a pointer")
	}

	argType := rv.Elem().Type()
	config := make([]internal.Config, argType.NumField())

	for i := range argType.NumField() {
		field := argType.Field(i)

		fieldName := field.Tag.Get("name")
		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		}

		fieldHelp := field.Tag.Get("help")
		fieldType := field.Tag.Get("type")
		if fieldType == "" || fieldType == "arg" {
			config[i] = internal.Arg{Name: fieldName, Help: fieldHelp}
		}
		if fieldType == "option" {
			fieldShort := field.Tag.Get("short")
			config[i] = internal.Option{Name: fieldName, Help: fieldHelp, Short: fieldShort}
		}
	}

	args := os.Args[1:]
	cmdName := os.Args[0]

	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			fmt.Fprintln(os.Stderr, internal.GenerateHelp(cmdName, config))
			os.Exit(0)
		}
	}

	return nil
}
