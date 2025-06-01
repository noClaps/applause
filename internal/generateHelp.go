package internal

import (
	"fmt"
	"os"
	"strings"
)

func generateHelp(argsConfig []arg, optionsConfig []option) string {
	maxLen := 10 // length of `-h, --help`
	for _, arg := range argsConfig {
		maxLen = max(len(arg.Name)+2, maxLen) // add 2 for `<>`
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
		arguments += fmt.Sprintf(
			"  <%s>%s        %s\n",
			arg.Name, strings.Repeat(" ", maxLen-len(arg.Name)-2), help,
		)
	}
	arguments = strings.TrimSpace(arguments)

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

	help := fmt.Sprintf("%s\n\n%s\n\n%s", generateUsage(argsConfig, optionsConfig), arguments, options)
	return strings.TrimSpace(help)
}

func generateUsage(argsConfig []arg, optionsConfig []option) string {
	cmdName := os.Args[0]
	arguments := ""
	for _, arg := range argsConfig {
		arguments += fmt.Sprintf("<%s> ", arg.Name)
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
	return fmt.Sprintf("USAGE: %s %s %s", cmdName, strings.TrimSpace(arguments), strings.TrimSpace(options))
}
