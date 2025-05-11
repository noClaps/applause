package internal

import (
	"fmt"
	"strings"
)

func GenerateHelp(argname string, conf map[string]Config) string {
	argsList := ""
	optionsList := ""
	maxLen := 0
	for name, c := range conf {
		if _, ok := c.(Arg); ok {
			argsList += fmt.Sprintf("<%s> ", name)
			if len(name) > maxLen {
				maxLen = len(name)
			}
			continue
		}
		if o, ok := c.(Option); ok {
			optionsList += fmt.Sprintf("[--%s] ", name)
			if len(name)+2 > maxLen {
				maxLen = len(name) + 2
			}
			if o.Short != "" {
				maxLen += 4
			}
			continue
		}
	}
	argsList = strings.TrimSpace(argsList)
	optionsList = strings.TrimSpace(optionsList)

	argsHelp := ""
	optionsHelp := ""
	for name, c := range conf {
		if a, ok := c.(Arg); ok {
			if argsHelp == "" {
				argsHelp = "\nARGUMENTS:\n"
			}
			argsHelp += fmt.Sprintf("  %s        %s\n", name+strings.Repeat(" ", maxLen-len(name)), a.Help)
			continue
		}
		if o, ok := c.(Option); ok {
			if optionsHelp == "" {
				optionsHelp = "\nOPTIONS:\n"
			}
			blankSpace := strings.Repeat(" ", maxLen-len(name)-2) // account for `--`
			shortString := ""
			if o.Short != "" {
				shortString = fmt.Sprintf("-%s, ", o.Short)
				blankSpace = strings.Repeat(" ", maxLen-len(name)-6) // account for `--` and `-c, `
			}
			optionsHelp += fmt.Sprintf("  %s--%s%s        %s\n", shortString, name, blankSpace, o.Help)
			continue
		}
	}

	return strings.TrimSpace(fmt.Sprintf(`
USAGE:
  %s %s %s%s%s
`, argname, argsList, optionsList, argsHelp, optionsHelp))
}
