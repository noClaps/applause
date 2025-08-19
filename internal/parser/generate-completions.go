package parser

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"slices"
	"strings"
)

func (p *Parser) GenerateCompletions(shell string) (string, error) {
	if shell == "" {
		shell = path.Base(os.Getenv("SHELL"))
	}
	switch shell {
	case "zsh":
		return fmt.Sprintf("#compdef %s\n%s", p.Name, p.generateZshCompletions(0)), nil
	}
	return "", fmt.Errorf("Shell not supported: %s", shell)
}

func (p *Parser) generateZshCompletions(indent int) string {
	indentSmall := strings.Repeat(" ", indent)
	indentMedium := strings.Repeat(" ", indent+2)
	indentLarge := strings.Repeat(" ", indent+4)
	indentXL := strings.Repeat(" ", indent+6)

	options := make([]string, len(p.Options))
	for i, opt := range p.Options {
		if opt.Type.Kind() == reflect.Bool {
			if opt.Short != "" {
				options[i] = fmt.Sprintf(
					"'(-%[1]s --%[2]s)'{-%[1]s,--%[2]s}'[%s]'",
					opt.Short, opt.Name, opt.Help)
				continue
			}
			options[i] = fmt.Sprintf("'--%s[%s]'", opt.Name, opt.Help)
			continue
		}

		completion := fmt.Sprintf("--%s[%s]:%[1]s:", opt.Name, opt.Help)
		if opt.Short != "" {
			completion = fmt.Sprintf(
				"(-%[1]s --%[2]s)'{-%[1]s,--%[2]s}'[%s]:%[2]s:",
				opt.Short, opt.Name, opt.Help)
		}
		if opt.Completion == "" {
			options[i] = fmt.Sprintf(`'%s'`, completion)
			continue
		}
		if opt.Completion == "files" || strings.HasPrefix(opt.Completion, "files[") {
			completion += "_files"
			if strings.HasPrefix(opt.Completion, "files[") {
				guard := opt.Completion[6 : len(opt.Completion)-1]
				completion += fmt.Sprintf(` -g "%s"`, guard)
			}
			options[i] = fmt.Sprintf(`'%s'`, completion)
			continue
		}
		if strings.HasPrefix(opt.Completion, "$(") && strings.HasSuffix(opt.Completion, ")") {
			cleanComp := strings.ReplaceAll(opt.Completion, "\n", `\n`)
			completion += fmt.Sprintf(`(%s)`, cleanComp)
			options[i] = fmt.Sprintf(`'%s'`, completion)
			continue
		}

		values := strings.Split(opt.Completion, " ")
		values = slices.DeleteFunc(values, func(v string) bool {
			return v == ""
		})
		completion += fmt.Sprintf(`_values "%s"`, opt.Name)
		for _, v := range values {
			completion += fmt.Sprintf(` "%s"`, v)
		}
		options[i] = fmt.Sprintf(`'%s'`, completion)
	}
	if len(p.Commands) == 0 {
		posAndOpts := fmt.Sprintf(
			"%[1]s_arguments '(-h --help)'{-h,--help}'[Display this help and exit.]' %s",
			indentSmall, strings.Join(options, " "))
		posCompletions := make([]string, 0, len(p.Positionals))
		for i, pos := range p.Positionals {
			if pos.Completion == "" {
				continue
			}
			completion := fmt.Sprintf("%d:%s:", i+1, pos.Name)
			if pos.Type.Kind() == reflect.Slice {
				completion = fmt.Sprintf("*:%s:", pos.Name)
			}
			if pos.Completion == "files" || strings.HasPrefix(pos.Completion, "files[") {
				completion += "_files"
				if strings.HasPrefix(pos.Completion, "files[") {
					guard := pos.Completion[6 : len(pos.Completion)-1]
					completion += fmt.Sprintf(` -g "%s"`, guard)
				}
				posCompletions = append(posCompletions, fmt.Sprintf(`'%s'`, completion))
				continue
			}
			if strings.HasPrefix(pos.Completion, "$(") && strings.HasSuffix(pos.Completion, ")") {
				cleanComp := strings.ReplaceAll(pos.Completion, "\n", `\n`)
				completion += fmt.Sprintf(`(%s)`, cleanComp)
				posCompletions = append(posCompletions, fmt.Sprintf(`"%s"`, completion))
				continue
			}

			values := strings.Split(pos.Completion, " ")
			values = slices.DeleteFunc(values, func(v string) bool {
				return v == ""
			})

			completion += fmt.Sprintf(`_values "%s"`, pos.Name)
			for _, v := range values {
				completion += fmt.Sprintf(` "%s"`, v)
			}
			posCompletions = append(posCompletions, fmt.Sprintf(`'%s'`, completion))
		}
		posAndOpts += strings.Join(posCompletions, " ")
		return strings.TrimSpace(posAndOpts)
	}

	commands := make([]string, len(p.Commands))
	commandCompletions := make([]string, len(p.Commands))
	for i, cmd := range p.Commands {
		commands[i] = fmt.Sprintf("'%s[%s]'", cmd.Name, cmd.Help)
		if cmd.Value.Elem().Kind() == reflect.Bool {
			commandCompletions[i] = fmt.Sprintf(
				"%s%s) _arguments '(-h --help)'{-h,--help}'[Display this help and exit.]' ;;",
				indentLarge, cmd.Name)
			continue
		}
		cmdParser := NewParser(p.Name+" "+cmd.Name, []string{}, cmd.Value)
		completions := cmdParser.generateZshCompletions(indent + 6)
		commandCompletions[i] = fmt.Sprintf("%[1]s%[2]s) %[4]s ;;", indentLarge, cmd.Name, indentXL, completions)
	}

	args := ";;"
	if len(commandCompletions) > 0 {
		args = fmt.Sprintf(
			"case $words[1] in\n%[2]s\n%[1]sesac ;;",
			indentMedium, strings.Join(commandCompletions, "\n"))
	}

	return fmt.Sprintf(`
%[1]s_arguments -C '1: :->first' '*:: :->args'
%[1]scase $state in
%[2]sfirst) _values 'global options and subcommands' '(-h --help)'{-h,--help}'[Display this help and exit.]' %[3]s %[4]s ;;
%[2]sargs) %[5]s
%[1]sesac`, indentSmall, indentMedium, strings.Join(options, " "), strings.Join(commands, " "), args)
}
