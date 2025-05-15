package internal

import (
	"fmt"
	"os"
)

func logError(err error) {
	fmt.Fprintln(os.Stderr, "\033[31mERROR:\033[0m", err)
}

func printUsage(argsConfig []arg, optionsConfig []option) {
	fmt.Fprintln(os.Stderr, generateUsage(argsConfig, optionsConfig))
}
