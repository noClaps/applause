package parser

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/noclaps/applause/internal/utils"
)

func (p *Parser) parseOptions() error {
	arguments := []string{}
	for i := 0; i < len(p.Arguments); i++ {
		arg := p.Arguments[i]

		// done parsing opts, everything after is arguments
		if arg == "--" {
			p.Arguments = p.Arguments[i+1:]
			return nil
		}

		// long option
		if len(arg) > 2 && arg[:2] == "--" {
			// --key=val
			if si := strings.Index(arg, "="); si != -1 {
				key := arg[2:si]
				val := arg[si+1:]
				optIndex := p.FindOptionByName(key)
				if optIndex == -1 {
					return fmt.Errorf("`%s` is not a recognised option.", arg[:si])
				}

				parsedVal, err := utils.ValToType(val, p.Options[optIndex].Type)
				if err != nil {
					return err
				}

				p.ParsedVals[key] = parsedVal
				continue
			}

			// --key val
			key := arg[2:]
			optIndex := p.FindOptionByName(key)
			if optIndex == -1 {
				return fmt.Errorf("`%s` is not a recognised option", arg)
			}
			if p.Options[optIndex].Type.Kind() == reflect.Bool {
				p.ParsedVals[key] = reflect.ValueOf(true)
				continue
			}

			if len(p.Arguments) <= i+1 {
				return fmt.Errorf("Value not provided for option `%s`", arg)
			}
			val := p.Arguments[i+1]
			if (len(val) > 2 && p.FindOptionByName(val[2:]) != -1) || (len(val) > 1 && p.FindOptionByShort(val[1:]) != -1) {
				return fmt.Errorf("Value not provided for option `%s`", arg)
			}

			parsedVal, err := utils.ValToType(val, p.Options[optIndex].Type)
			if err != nil {
				return err
			}

			p.ParsedVals[key] = parsedVal
			i++
			continue
		}

		// short option
		if len(arg) > 1 && arg[0] == '-' {
			optionName := arg[1:]
			optIndex := p.FindOptionByShort(optionName)
			if optIndex == -1 {
				return fmt.Errorf("`%s` is not a recognised option.", arg)
			}
			name := p.Options[optIndex].Name
			if p.Options[optIndex].Type.Kind() == reflect.Bool {
				p.ParsedVals[name] = reflect.ValueOf(true)
				continue
			}

			if len(p.Arguments) == i+1 {
				return fmt.Errorf("Value not provided for option `%s`", arg)
			}
			val := p.Arguments[i+1]
			if (len(val) > 2 && p.FindOptionByName(val[2:]) != -1) || (len(val) > 1 && p.FindOptionByShort(val[1:]) != -1) {
				return fmt.Errorf("Value not provided for option `%s`", arg)
			}

			parsedVal, err := utils.ValToType(val, p.Options[optIndex].Type)
			if err != nil {
				return err
			}

			p.ParsedVals[name] = parsedVal
			i++
			continue
		}

		arguments = append(arguments, arg)
	}

	p.Arguments = arguments
	return nil
}

func (p *Parser) parsePositionals() error {
	currentPosCounter := 0
	for i := 0; i < len(p.Arguments); i++ {
		arg := p.Arguments[i]

		if currentPosCounter == len(p.Positionals) {
			return fmt.Errorf("Extra argument: `%s`", arg)
		}

		currentPos := p.Positionals[currentPosCounter]
		name := currentPos.Name

		// Multiple arguments
		if currentPos.Type.Kind() == reflect.Slice {
			slice := reflect.MakeSlice(currentPos.Type, 0, len(p.Positionals)-currentPosCounter-1)
			posType := currentPos.Type.Elem()

			for ; len(p.Arguments)-i != len(p.Positionals)-currentPosCounter-1; i++ {
				arg = p.Arguments[i]

				val, err := utils.ValToType(arg, posType)
				if err != nil {
					return err
				}

				slice = reflect.Append(slice, val)
			}

			p.ParsedVals[name] = slice
			currentPosCounter++
			i--
			continue
		}

		// Read from stdin
		if arg == "-" {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			stdinVal := scanner.Text()

			val, err := utils.ValToType(stdinVal, currentPos.Type)
			if err != nil {
				return err
			}

			p.ParsedVals[name] = val
			currentPosCounter++
			continue
		}

		val, err := utils.ValToType(arg, currentPos.Type)
		if err != nil {
			return err
		}

		p.ParsedVals[name] = val
		currentPosCounter++
	}
	if currentPosCounter < len(p.Positionals) {
		return fmt.Errorf("Not enough arguments provided.")
	}

	return nil
}
