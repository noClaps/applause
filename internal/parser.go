package internal

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func Parse(argsConfig []arg, optionsConfig []option) (map[string]any, error) {
	parsedVals := make(map[string]any)

	cmdArgs := os.Args[1:]
	if len(argsConfig) > 0 && len(cmdArgs) == 0 || slices.ContainsFunc(cmdArgs, func(arg string) bool {
		return arg == "--help" || arg == "-h"
	}) {
		fmt.Println(GenerateHelp(argsConfig, optionsConfig))
		os.Exit(0)
	}

	arguments := []string{}
	onlyParseArgs := false
	for i := 0; i < len(cmdArgs); i++ {
		arg := cmdArgs[i]

		if arg == "--" {
			onlyParseArgs = true
			continue
		}

		// Long option
		if len(arg) >= 2 && arg[:2] == "--" && !onlyParseArgs {
			if si := strings.Index(arg, "="); si != -1 {
				key := arg[2:si]
				val := arg[si+1:]
				optIndex := slices.IndexFunc(optionsConfig, func(o option) bool {
					return o.Name == key
				})
				if optIndex == -1 {
					return nil, fmt.Errorf("`%s` is not a recognised option.", arg[:si])
				}
				parsedVal, err := valFromString(val, optionsConfig[optIndex].Type)
				if err != nil {
					return nil, fmt.Errorf("Error parsing val from `%v`: %v", val, err)
				}
				parsedVals[key] = parsedVal
				continue
			}

			key := arg[2:]
			optIndex := slices.IndexFunc(optionsConfig, func(o option) bool {
				return o.Name == key
			})
			if optIndex == -1 {
				return nil, fmt.Errorf("`%s` is not a recognised option", arg)
			}
			if optionsConfig[optIndex].Type == "bool" {
				parsedVals[key] = true
				continue
			}
			if len(cmdArgs) <= i+1 {
				return nil, fmt.Errorf("Value not provided for option `%s`", arg)
			}
			val := cmdArgs[i+1]
			parsedVal, err := valFromString(val, optionsConfig[optIndex].Type)
			if err != nil {
				return nil, fmt.Errorf("Error parsing val from `%v`: %v", val, err)
			}
			parsedVals[key] = parsedVal
			i++
			continue
		}

		if len(arg) > 1 && arg[0] == '-' && !onlyParseArgs {
			optionName := arg[1:]
			optIndex := slices.IndexFunc(optionsConfig, func(o option) bool {
				return o.Short == optionName
			})
			if optIndex == -1 {
				return nil, fmt.Errorf("`%s` is not a recognised option.", arg)
			}
			name := optionsConfig[optIndex].Name
			if optionsConfig[optIndex].Type == "bool" {
				parsedVals[name] = true
				continue
			}
			val, err := valFromString(cmdArgs[i+1], optionsConfig[optIndex].Type)
			if err != nil {
				return nil, fmt.Errorf("Error parsing val from `%v`: %v", cmdArgs[i+1], err)
			}
			parsedVals[name] = val
			i++
			continue
		}

		arguments = append(arguments, arg)
	}

	currentArgCounter := 0
	for i := 0; i < len(arguments); i++ {
		arg := arguments[i]

		if currentArgCounter == len(argsConfig) {
			return nil, fmt.Errorf("Extra argument: `%s`", arg)
		}

		currentArg := argsConfig[currentArgCounter]
		name := currentArg.Name

		// Multiple arguments
		if currentArg.Type[0:2] == "[]" {
			argType := currentArg.Type[2:]
			parsedArr := []any{}
			for ; len(arguments)-i != len(argsConfig)-currentArgCounter-1; i++ {
				arg = arguments[i]
				val, err := valFromString(arg, argType)
				if err != nil {
					return nil, fmt.Errorf("Error parsing val from `%v`: %v", arg, err)
				}
				parsedArr = append(parsedArr, val)
			}
			parsedVals[name] = parsedArr
			currentArgCounter++
			i--
			continue
		}

		if arg == "-" {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			stdinVal := scanner.Text()
			val, err := valFromString(stdinVal, currentArg.Type)
			if err != nil {
				return nil, fmt.Errorf("Error parsing val from `%v`: %v", stdinVal, err)
			}
			parsedVals[name] = val
			currentArgCounter++
			continue
		}

		val, err := valFromString(arg, currentArg.Type)
		if err != nil {
			return nil, fmt.Errorf("Error parsing val from `%v`: %v", arg, err)
		}
		parsedVals[name] = val
		cmdArgs = slices.Delete(cmdArgs, 0, 1)
		currentArgCounter++
	}
	if currentArgCounter < len(argsConfig) {
		return nil, fmt.Errorf("Not enough arguments provided.")
	}

	return parsedVals, nil
}

// not a huge fan of this solution since it assumes that these will be the only
// data types that can be input from the CLI, but that may not necessarily be
// the case indefinitely.
func valFromString(input string, returnType string) (any, error) {
	switch returnType {
	case "bool":
		b, err := strconv.ParseBool(input)
		if err != nil {
			return nil, fmt.Errorf("Error parsing bool: %v", err)
		}
		return bool(b), nil
	case "float32":
		f, err := strconv.ParseFloat(input, 32)
		if err != nil {
			return nil, fmt.Errorf("Error parsing float32: %v", err)
		}
		return float32(f), nil
	case "float64":
		f, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing float64: %v", err)
		}
		return float64(f), nil
	case "int":
		i, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing int: %v", err)
		}
		return int(i), nil
	case "int8":
		i, err := strconv.ParseInt(input, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("Error parsing int8: %v", err)
		}
		return int8(i), nil
	case "int16":
		i, err := strconv.ParseInt(input, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("Error parsing int16: %v", err)
		}
		return int16(i), nil
	case "int32":
		i, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Error parsing int32: %v", err)
		}
		return int32(i), nil
	case "int64":
		i, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing int64: %v", err)
		}
		return int64(i), nil
	case "uint":
		u, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing uint: %v", err)
		}
		return uint(u), nil
	case "uint8":
		u, err := strconv.ParseUint(input, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("Error parsing uint8: %v", err)
		}
		return uint8(u), nil
	case "uint16":
		u, err := strconv.ParseUint(input, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("Error parsing uint16: %v", err)
		}
		return uint16(u), nil
	case "uint32":
		u, err := strconv.ParseUint(input, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Error parsing uint32: %v", err)
		}
		return uint32(u), nil
	case "uint64":
		u, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing uint64: %v", err)
		}
		return uint64(u), nil
	case "complex64":
		c, err := strconv.ParseComplex(input, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing complex64: %v", err)
		}
		return complex64(c), nil
	case "complex128":
		c, err := strconv.ParseComplex(input, 128)
		if err != nil {
			return nil, fmt.Errorf("Error parsing complex128: %v", err)
		}
		return complex128(c), nil
	case "string":
		return input, nil
	}

	return nil, fmt.Errorf("Type `%s` is unsupported, please use a supported type.", returnType)
}
