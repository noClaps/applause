package main

import (
	"log"

	"github.com/noclaps/applause"
)

type Args struct {
	Add struct {
		Packages []string `help:"Packages to install."`
	} `help:"Install packages."`
	Update *struct {
		Packages []string `help:"Packages to update."`
	} `help:"Update packages."`
	Remove struct {
		Packages []string `help:"Packages to remove."`
	} `help:"Remove packages."`
	Info struct {
		Package string `help:"The package to get the info for"`
	} `help:"Get the info for a package."`
	List bool `type:"command" help:"List installed packages"`
	Init bool `type:"option" help:"Initialise pkg"`
}

func main() {
	args := Args{}
	if err := applause.Parse(&args); err != nil {
		log.Fatalln(err)
	}

	if len(args.Add.Packages) > 0 {
		log.Println(args.Add.Packages)
		return
	}
	if args.Update != nil {
		log.Println(args.Update.Packages)
		return
	}

	log.Println(args)
}
