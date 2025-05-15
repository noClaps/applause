package main

import (
	"fmt"
	"log"
	"os"

	"github.com/noclaps/applause"
)

type Args struct {
	Name         string `help:"The name of the package"`
	SkipPeer     bool   `type:"option" name:"skip-peer" help:"Skip counting peer dependencies" short:"p"`
	SkipOptional bool   `type:"option" name:"skip-optional" help:"Skip counting optional dependencies" short:"o"`
	Version      string `type:"option" help:"The version of the package being fetched" value:"version"`
}

func main() {
	args := Args{}
	if err := applause.Parse(&args); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
	log.Printf("%#v", args)
}
