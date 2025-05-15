package internal

import (
	"fmt"
	"strings"
)

func GenerateHelp(argname string, conf []Config) string {
	argsList := ""
	optionsList := ""
	maxLen := 10 // Minimum 10 because of `-h, --help`
	for _, c := range conf {
		if _, ok := c.(Arg); ok {
			c := c.(Arg)
			argsList += fmt.Sprintf("<%s> ", c.Name)
			if len(c.Name) > maxLen {
				maxLen = len(c.Name)
			}
			continue
		}
		if o, ok := c.(Option); ok {
			c := c.(Option)
			optionsList += fmt.Sprintf("[--%s] ", c.Name)
			if len(c.Name)+2 > maxLen {
				maxLen = len(c.Name) + 2
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
	for _, c := range conf {
		if a, ok := c.(Arg); ok {
			c := c.(Arg)
			if argsHelp == "" {
				argsHelp = "\nARGUMENTS:\n"
			}
			argsHelp += fmt.Sprintf("  %s        %s\n", c.Name+strings.Repeat(" ", maxLen-len(c.Name)), a.Help)
			continue
		}
		if o, ok := c.(Option); ok {
			c := c.(Option)
			if optionsHelp == "" {
				optionsHelp = "\nOPTIONS:\n"
			}
			blankSpace := strings.Repeat(" ", maxLen-len(c.Name)-2) // account for `--`
			shortString := ""
			if o.Short != "" {
				shortString = fmt.Sprintf("-%s, ", o.Short)
				blankSpace = strings.Repeat(" ", maxLen-len(c.Name)-6) // account for `--` and `-c, `
			}
			optionsHelp += fmt.Sprintf("  %s--%s%s        %s\n", shortString, c.Name, blankSpace, o.Help)
			continue
		}
	}
	optionsHelp += fmt.Sprintf("  -h, --help%s        Display this help and exit.\n", strings.Repeat(" ", maxLen-10))

	return strings.TrimSpace(fmt.Sprintf(`
USAGE:
  %s %s %s%s%s
`, argname, argsList, optionsList, argsHelp, optionsHelp))
}
