package utils

import (
	"os"
	"strings"
)

func GetCmdName() string {
	cmdPath := strings.Split(os.Args[0], "/")
	cmdName := cmdPath[len(cmdPath)-1]
	return cmdName
}
