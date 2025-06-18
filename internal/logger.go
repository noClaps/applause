package internal

import (
	"fmt"
	"os"
)

func printUsage(argsConfig []arg, optionsConfig []option) {
	fmt.Fprintln(os.Stderr, generateUsage(argsConfig, optionsConfig))
}
