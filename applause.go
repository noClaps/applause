package applause

import (
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/noclaps/applause/internal/parser"
)

// The help string for the command. This will only contain a value if
// [Parse] has been called first, otherwise it will be an empty
// string.
var Help string = ""

// The usage string for the command. This will only contain a value if
// [Parse] has been called first, otherwise it will be an empty
// string.
var Usage string = ""

/*
The input is a pointer to the args struct. Each field in the args struct
should have some tags:

  - `type`: The type can be "arg", "option" or "command". If omitted, the
    default is "arg". If any other type is provided, the field is ignored.

  - `name`: The name of the argument, option or command. If omitted, the
    default is the field name in kebab-case. If you'd like to have an option
    have a different name, you can write it as `name:"option-name"` in the
    tags.

  - `help`: The help text for the argument, option or command, will be
    displayed in the command help when the command is called with `--help`
    or `-h`.

  - `short`: Only applicable when `type` is "option". The short form of the
    option. For instance, if you have a field with the tag
    `name:"option" short:"o"`, you can call the command with `--option` or
    `-o`.

  - `value`: Only applicable when `type` is "option" and the field type is
    not "bool". The name of the option value to be displayed in the help
    text. For instance, `name:"option" value:"val"` will be displayed as
    `--option <val>` in the help text.

  - `completion`: Only applicable when `type` is "arg", "option" or omitted.
    You can define a completion one of three ways:

    If you do `completion:"files"`, it will autocomplete to files. This is
    useful to complete file paths.

    If you do `completion:"some values here"`, it will autocomplete to
    "some", "values" and "here". This is useful when you have a fixed set of
    values you want to autocomplete.

    If you do `completion:"$(echo 'some command here')"`, it will
    autocomplete to the output of that command at runtime. This is useful
    for completions that need to be dynamic.

All fields that you'd like to be parsed should be exported in the struct.
*/
func Parse(args any) error {
	rv := reflect.ValueOf(args)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("Input value should be a pointer to a struct, received: %v", rv.Kind().String())
	}

	cmdName := path.Base(os.Args[0])
	parser := parser.NewParser(cmdName, os.Args[1:], rv)
	Help = parser.Help
	Usage = parser.Usage

	if err := parser.Parse(); err != nil {
		return err
	}

	return nil
}
