package parser

import (
	"fmt"
	"os"
	"reflect"
	"slices"
)

type Parser struct {
	Name        string        // command name
	Arguments   []string      // OS arguments
	Config      reflect.Value // pointer to config struct
	Positionals []positional
	Options     []option
	Commands    []command
	Help        string
	Usage       string
	ParsedVals  map[string]reflect.Value
}

// config should be a pointer to a struct
func NewParser(cmdName string, args []string, config reflect.Value) *Parser {
	p := Parser{
		Name:       cmdName,
		Arguments:  args,
		ParsedVals: make(map[string]reflect.Value),
		Config:     config,
	}

	p.reflection()

	p.generateUsage()
	p.generateHelp()

	return &p
}

func (p *Parser) Parse() error {
	if len(p.Commands) > 0 {
		if len(p.Arguments) == 0 || p.Arguments[0] == "-h" || p.Arguments[0] == "--help" {
			fmt.Println(p.Help)
			os.Exit(0)
		}
		cIndex := p.FindComandByName(p.Arguments[0])
		if cIndex != -1 {
			command := p.Commands[cIndex]
			nestedCmdName := fmt.Sprintf("%s %s", p.Name, command.Name)
			nestedP := NewParser(nestedCmdName, p.Arguments[1:], command.Value)
			return nestedP.Parse()
		}
	}

	if len(p.Positionals) > 0 && len(p.Arguments) == 0 || slices.ContainsFunc(p.Arguments, func(arg string) bool {
		return arg == "--help" || arg == "-h"
	}) {
		fmt.Println(p.Help)
		os.Exit(0)
	}

	if err := p.parseOptions(); err != nil {
		return err
	}
	if err := p.parsePositionals(); err != nil {
		return err
	}

	for k, v := range p.ParsedVals {
		if posIndex := p.FindPositionalByName(k); posIndex != -1 {
			positional := p.Positionals[posIndex]
			p.Config.Elem().FieldByName(positional.StructName).Set(v)
		}
		if optIndex := p.FindOptionByName(k); optIndex != -1 {
			option := p.Options[optIndex]
			p.Config.Elem().FieldByName(option.StructName).Set(v)
		}
	}

	return nil
}
