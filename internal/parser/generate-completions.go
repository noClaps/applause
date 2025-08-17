package parser

import (
	"fmt"
	"os"
	"path"
	"reflect"
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
		if opt.Short != "" {
			switch opt.Type.Kind() {
			case reflect.Bool:
				options[i] = fmt.Sprintf(
					"'(-%[1]s --%[2]s)'{-%[1]s,--%[2]s}'[%s]'",
					opt.Short, opt.Name, opt.Help)
			default:
				options[i] = fmt.Sprintf(
					"'(-%[1]s --%[2]s)'{-%[1]s+,--%[2]s=}'[%s]:option:'",
					opt.Short, opt.Name, opt.Help)
			}
		} else {
			switch opt.Type.Kind() {
			case reflect.Bool:
				options[i] = fmt.Sprintf("'--%s[%s]'", opt.Name, opt.Help)
			default:
				options[i] = fmt.Sprintf("'--%s=[%s]:option:'", opt.Name, opt.Help)
			}
		}
	}
	if len(p.Commands) == 0 {
		return strings.TrimSpace(
			fmt.Sprintf(
				"%[1]s_arguments '(-h --help)'{-h,--help}'[Display this help and exit.]' %s",
				indentSmall, strings.Join(options, " ")),
		)
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
