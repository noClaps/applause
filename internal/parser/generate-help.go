package parser

import (
	"fmt"
	"reflect"
	"strings"
)

func (p *Parser) generateUsage() {
	positionalUsage := ""
	for _, positional := range p.Positionals {
		if positional.Type.Kind() == reflect.Slice {
			positionalUsage += fmt.Sprintf("[%s...] ", positional.Name)
			continue
		}
		positionalUsage += fmt.Sprintf("<%s> ", positional.Name)
	}
	positionalUsage = strings.TrimSpace(positionalUsage)
	if positionalUsage != "" {
		positionalUsage += " "
	}

	optionUsage := ""
	for _, option := range p.Options {
		optionUsagePart := "["
		if option.Name != "" {
			optionUsagePart += fmt.Sprintf("--%s", option.Name)
		} else {
			if option.Short == "" {
				continue
			}
			optionUsagePart += fmt.Sprintf("-%s", option.Short)
		}
		if option.Value != "" {
			optionUsagePart += fmt.Sprintf(" <%s>", option.Value)
		}
		optionUsagePart += "] "
		optionUsage += optionUsagePart
	}
	optionUsage = strings.TrimSpace(optionUsage)

	p.Usage = fmt.Sprintf("USAGE: %s %s%s", p.Name, positionalUsage, optionUsage)
}

func (p *Parser) generateHelp() {
	maxLen := 10 // length of `-h, --help`
	for _, positional := range p.Positionals {
		maxLen = max(len(positional.Name)+2, maxLen) // add 2 for `<>`
		if positional.Type.Kind() == reflect.Slice {
			maxLen += 3 // add 3 for `...`
		}
	}
	for _, option := range p.Options {
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

	positionalHelp := ""
	for _, positional := range p.Positionals {
		if positionalHelp == "" {
			positionalHelp = "\nARGUMENTS:\n"
		}
		help := positional.Help
		if len(positional.Help) > 80 {
			help = ""
			words := strings.Split(positional.Help, " ")
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
		if positional.Type.Kind() == reflect.Slice {
			positionalHelp += fmt.Sprintf(
				"  [%s...]%s        %s\n",
				positional.Name, strings.Repeat(" ", maxLen-len(positional.Name)-5), help,
			)
			continue
		}
		positionalHelp += fmt.Sprintf(
			"  <%s>%s        %s\n",
			positional.Name, strings.Repeat(" ", maxLen-len(positional.Name)-2), help,
		)
	}
	positionalHelp = strings.TrimSpace(positionalHelp)
	if positionalHelp != "" {
		positionalHelp += "\n\n"
	}

	optionHelp := "OPTIONS:\n"
	for _, option := range p.Options {
		optLen := 0
		name := ""
		if option.Name != "" {
			name = fmt.Sprintf("--%s", option.Name)
			optLen += len(option.Name) + 2 // add `--`
		}
		short := ""
		if option.Short != "" {
			if option.Name == "" {
				short = fmt.Sprintf("-%s", option.Short)
				optLen += len(option.Short) + 1 // add `-`

			} else {
				short = fmt.Sprintf("-%s, ", option.Short)
				optLen += len(option.Short) + 3 // add `-, `
			}
		}
		value := ""
		if option.Value != "" {
			value = fmt.Sprintf(" <%s>", option.Value)
			optLen += len(option.Value) + 3 // add ` <>`
		}
		defaultStr := ""
		if !option.Default.IsZero() {
			defaultStr = fmt.Sprintf(" (default: %v)", option.Default)
		}
		help := option.Help
		if len(option.Help) > 80 {
			help = ""
			words := strings.Split(option.Help, " ")
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
		optionHelp += fmt.Sprintf(
			"  %s%s%s%s        %s%s\n",
			short, name, value, strings.Repeat(" ", maxLen-optLen), help, defaultStr,
		)
	}
	optionHelp += fmt.Sprintf("  -h, --help%s        Display this help and exit.", strings.Repeat(" ", maxLen-10))
	optionHelp = strings.TrimSpace(optionHelp)

	help := fmt.Sprintf("%s\n\n%s%s", p.Usage, positionalHelp, optionHelp)
	p.Help = strings.TrimSpace(help)
}
