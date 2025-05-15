package internal

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

func Parse(argStruct reflect.Type) (map[string]any, error) {
	argsConfig, optionsConfig, err := HandleReflection(argStruct)
	if err != nil {
		return nil, err
	}
	parsedVals := make(map[string]any)

	cmdArgs := os.Args[1:]
	if len(argsConfig) > 0 && len(cmdArgs) == 0 {
		fmt.Println(generateHelp(argsConfig, optionsConfig))
		os.Exit(0)
	}

	currentArgCounter := 0
	for len(cmdArgs) > 0 {
		arg := cmdArgs[0]

		if arg == "--help" || arg == "-h" {
			fmt.Fprintln(os.Stderr, generateHelp(argsConfig, optionsConfig))
			os.Exit(0)
		}

		// Long option
		if len(arg) >= 2 && arg[:2] == "--" {
			if i := strings.Index(arg, "="); i != -1 {
				key := arg[2:i]
				val := arg[i+1:]
				optIndex := slices.IndexFunc(optionsConfig, func(o option) bool {
					return o.Name == key
				})
				if optIndex == -1 {
					logError(fmt.Errorf("`%s` is not a recognised option.", arg[:i]))
					printUsage(argsConfig, optionsConfig)
					os.Exit(1)
				}
				parsedVal, err := valFromString(val, optionsConfig[optIndex].Type)
				if err != nil {
					return nil, err
				}
				parsedVals[key] = parsedVal
				cmdArgs = slices.Delete(cmdArgs, 0, 1)
				continue
			}

			key := arg[2:]
			optIndex := slices.IndexFunc(optionsConfig, func(o option) bool {
				return o.Name == key
			})
			if optIndex == -1 {
				logError(fmt.Errorf("`%s` is not a recognised option", arg))
				printUsage(argsConfig, optionsConfig)
				os.Exit(1)
			}
			if optionsConfig[optIndex].Type == "bool" {
				parsedVals[key] = true
				cmdArgs = slices.Delete(cmdArgs, 0, 1)
				continue
			}
			if len(cmdArgs) <= 1 {
				logError(fmt.Errorf("Value not provided for option `%s`", arg))
				os.Exit(1)
			}
			val := cmdArgs[1]
			parsedVal, err := valFromString(val, optionsConfig[optIndex].Type)
			if err != nil {
				return nil, err
			}
			parsedVals[key] = parsedVal
			cmdArgs = slices.Delete(cmdArgs, 0, 2)
			continue
		}

		// Short option
		if arg[0] == '-' {
			optionName := arg[1:]
			optIndex := slices.IndexFunc(optionsConfig, func(o option) bool {
				return o.Short == optionName
			})
			if optIndex == -1 {
				fmt.Fprintf(os.Stderr, "\033[31mERROR:\033[0m `%s` is not a recognised option.\n", arg)
				printUsage(argsConfig, optionsConfig)
				os.Exit(1)
			}
			name := optionsConfig[optIndex].Name
			if optionsConfig[optIndex].Type == "bool" {
				parsedVals[name] = true
				cmdArgs = slices.Delete(cmdArgs, 0, 1)
				continue
			}
			val, err := valFromString(cmdArgs[1], optionsConfig[optIndex].Type)
			if err != nil {
				return nil, err
			}
			parsedVals[name] = val
			cmdArgs = slices.Delete(cmdArgs, 0, 2)
			continue
		}

		if currentArgCounter >= len(argsConfig) {
			logError(fmt.Errorf("Extra argument: `%s`", arg))
			printUsage(argsConfig, optionsConfig)
			os.Exit(1)
		}
		if len(cmdArgs) < len(argsConfig) {
			logError(fmt.Errorf("Not enough arguments provided."))
			printUsage(argsConfig, optionsConfig)
		}
		currentArg := argsConfig[currentArgCounter]
		name := currentArg.Name
		val, err := valFromString(arg, currentArg.Type)
		if err != nil {
			return nil, err
		}
		parsedVals[name] = val
		currentArgCounter++
		cmdArgs = slices.Delete(cmdArgs, 0, 1)
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
			logError(err)
		}
		return bool(b), nil
	case "float32":
		f, err := strconv.ParseFloat(input, 32)
		if err != nil {
			logError(err)
		}
		return float32(f), nil
	case "float64":
		f, err := strconv.ParseFloat(input, 64)
		if err != nil {
			logError(err)
		}
		return float64(f), nil
	case "int":
		i, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			logError(err)
		}
		return int(i), nil
	case "int8":
		i, err := strconv.ParseInt(input, 10, 8)
		if err != nil {
			logError(err)
		}
		return int8(i), nil
	case "int16":
		i, err := strconv.ParseInt(input, 10, 16)
		if err != nil {
			logError(err)
		}
		return int16(i), nil
	case "int32":
		i, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			logError(err)
		}
		return int32(i), nil
	case "int64":
		i, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			logError(err)
		}
		return int64(i), nil
	case "uint":
		u, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			logError(err)
		}
		return uint(u), nil
	case "uint8":
		u, err := strconv.ParseUint(input, 10, 8)
		if err != nil {
			logError(err)
		}
		return uint8(u), nil
	case "uint16":
		u, err := strconv.ParseUint(input, 10, 16)
		if err != nil {
			logError(err)
		}
		return uint16(u), nil
	case "uint32":
		u, err := strconv.ParseUint(input, 10, 32)
		if err != nil {
			logError(err)
		}
		return uint32(u), nil
	case "uint64":
		u, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			logError(err)
		}
		return uint64(u), nil
	case "complex64":
		c, err := strconv.ParseComplex(input, 64)
		if err != nil {
			logError(err)
		}
		return complex64(c), nil
	case "complex128":
		c, err := strconv.ParseComplex(input, 128)
		if err != nil {
			logError(err)
		}
		return complex128(c), nil
	case "string":
		return input, nil
	}

	return nil, fmt.Errorf("Type `%s` is unsupported, please use a supported type.", returnType)
}
