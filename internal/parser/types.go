package parser

import (
	"reflect"
	"slices"
)

type positional struct {
	StructName string       // original name in struct
	Name       string       // positional name
	Help       string       // positional help
	Type       reflect.Type // positional type
}

type option struct {
	StructName string        // original name in struct
	Name       string        // option name
	Help       string        // option help
	Type       reflect.Type  // option type
	Short      string        // option short form
	Default    reflect.Value // option default value
	Value      string        // option argument name
}

// Returns the index of the named positional, otherwise -1 if the positional doesn't exist
func (p Parser) FindPositionalByName(name string) int {
	return slices.IndexFunc(p.Positionals, func(p positional) bool {
		return p.Name == name
	})
}

// Returns the index of the named option, otherwise -1 if the option doesn't exist
func (p Parser) FindOptionByName(name string) int {
	return slices.IndexFunc(p.Options, func(o option) bool {
		return o.Name == name
	})
}

// Returns the index of the named short option, otherwise -1 if the option doesn't exist
func (p Parser) FindOptionByShort(short string) int {
	return slices.IndexFunc(p.Options, func(o option) bool {
		return o.Short == short
	})
}
