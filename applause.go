package applause

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/noclaps/applause/internal"
)

func Parse(argStruct any) any {
	config := make(map[string]internal.Config)

	argType := reflect.TypeOf(argStruct)
	for i := range argType.NumField() {
		field := argType.Field(i)

		fieldName := field.Tag.Get("name")
		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		}

		fieldHelp := field.Tag.Get("help")
		fieldType := field.Tag.Get("type")
		if fieldType == "" || fieldType == "arg" {
			config[fieldName] = internal.Arg{Help: fieldHelp}
		}
		if fieldType == "option" {
			fieldShort := field.Tag.Get("short")

			config[fieldName] = internal.Option{Help: fieldHelp, Short: fieldShort}
		}
	}

	args := os.Args[1:]
	cmdName := os.Args[0]

	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			fmt.Fprintln(os.Stderr, internal.GenerateHelp(cmdName, config))
			os.Exit(0)
		}

		if arg[0:2] == "--" {
			eqIndex := strings.Index(arg[2:], "=")
			argName := arg[2:]
			if eqIndex != -1 {
				argName = arg[2:eqIndex]
			}
			if _, ok := config[argName]; !ok {
				fmt.Fprintln(os.Stderr, "Unknown option:", arg)
			}
		}
	}

	fmt.Printf("%#v\n", config)

	return argStruct
}
