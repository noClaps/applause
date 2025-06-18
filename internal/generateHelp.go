package internal

import (
	"fmt"
	"os"
	"strings"
)

func GenerateHelp(argsConfig []arg, optionsConfig []option) string {
	maxLen := 10 // length of `-h, --help`
	for _, arg := range argsConfig {
		maxLen = max(len(arg.Name)+2, maxLen) // add 2 for `<>`
		if arg.Type[0:2] == "[]" {
			maxLen += 3 // add 3 for `...`
		}
	}
	for _, option := range optionsConfig {
		optLen := 0
		if option.Name != "" {
			optLen += len(option.Name) + 2 // add 2 for `--`
		}
		if option.Value != "" {
			optLen += len(option.Value) + 3 // add 3 for ` <>`
		}
		if option.Short != "" {
			optLen += len(option.Short) + 3 // add 3 for `-, `
		}
		maxLen = max(optLen, maxLen)
	}
	arguments := ""
	for _, arg := range argsConfig {
		if arguments == "" {
			arguments = "\nARGUMENTS:\n"
		}
		help := arg.Help
		if len(arg.Help) > 80 {
			help = ""
			words := strings.Split(arg.Help, " ")
			spaceLen := 10 + maxLen // 2 for starting spaces + 8 for middle spaces
			i := 0
			for i < len(words) {
				line := ""
				for i < len(words) && len(line+words[i]) < 80 {
					line += words[i] + " "
					i++
				}
				line += "\n" + strings.Repeat(" ", spaceLen)
				help += line
			}
			help = strings.TrimSpace(help)
		}
		if arg.Type[0:2] == "[]" {
			arguments += fmt.Sprintf(
				"  [%s...]%s        %s\n",
				arg.Name, strings.Repeat(" ", maxLen-len(arg.Name)-5), help,
			)
			continue
		}
		arguments += fmt.Sprintf(
			"  <%s>%s        %s\n",
			arg.Name, strings.Repeat(" ", maxLen-len(arg.Name)-2), help,
		)
	}
	arguments = strings.TrimSpace(arguments)
	if arguments != "" {
		arguments += "\n\n"
	}

	options := "OPTIONS:\n"
	for _, opt := range optionsConfig {
		optLen := 0
		name := ""
		if opt.Name != "" {
			name = fmt.Sprintf("--%s", opt.Name)
			optLen += len(opt.Name) + 2 // add `--`
		}
		short := ""
		if opt.Short != "" {
			if opt.Name == "" {
				short = fmt.Sprintf("-%s", opt.Short)
				optLen += len(opt.Short) + 1 // add `-`

			} else {
				short = fmt.Sprintf("-%s, ", opt.Short)
				optLen += len(opt.Short) + 3 // add `-, `
			}
		}
		value := ""
		if opt.Value != "" {
			value = fmt.Sprintf(" <%s>", opt.Value)
			optLen += len(opt.Value) + 3 // add ` <>`
		}
		defaultStr := ""
		if opt.Default != "" {
			defaultStr = fmt.Sprintf(" (default: %s)", opt.Default)
		}
		help := opt.Help
		if len(opt.Help) > 80 {
			help = ""
			words := strings.Split(opt.Help, " ")
			spaceLen := 10 + maxLen // 2 for starting spaces + 8 for middle spaces
			i := 0
			for i < len(words) {
				line := ""
				for i < len(words) && len(line+words[i]) < 80 {
					line += words[i] + " "
					i++
				}
				line += "\n" + strings.Repeat(" ", spaceLen)
				help += line
			}
			help = strings.TrimSpace(help)
		}
		options += fmt.Sprintf(
			"  %s%s%s%s        %s%s\n",
			short, name, value, strings.Repeat(" ", maxLen-optLen), help, defaultStr,
		)
	}
	options += fmt.Sprintf("  -h, --help%s        Display this help and exit.", strings.Repeat(" ", maxLen-10))
	options = strings.TrimSpace(options)

	help := fmt.Sprintf("%s\n\n%s%s", GenerateUsage(argsConfig, optionsConfig), arguments, options)
	return strings.TrimSpace(help)
}

func GenerateUsage(argsConfig []arg, optionsConfig []option) string {
	cmdPath := strings.Split(os.Args[0], "/")
	cmdName := cmdPath[len(cmdPath)-1]
	arguments := ""
	for _, arg := range argsConfig {
		if arg.Type[0:2] == "[]" {
			arguments += fmt.Sprintf("[%s...] ", arg.Name)
			continue
		}
		arguments += fmt.Sprintf("<%s> ", arg.Name)
	}
	arguments = strings.TrimSpace(arguments)
	if arguments != "" {
		arguments += " "
	}

	options := ""
	for _, option := range optionsConfig {
		optionHelp := "["
		if option.Name != "" {
			optionHelp += fmt.Sprintf("--%s", option.Name)
		} else {
			if option.Short == "" {
				continue
			}
			optionHelp += fmt.Sprintf("-%s", option.Short)
		}
		if option.Value != "" {
			optionHelp += fmt.Sprintf(" <%s>", option.Value)
		}
		optionHelp += "] "
		options += optionHelp
	}
	options = strings.TrimSpace(options)
	return fmt.Sprintf("USAGE: %s %s%s", cmdName, arguments, options)
}
