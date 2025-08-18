package main

import (
	"log"

	"github.com/noclaps/applause"
)

type Args struct {
	Add struct {
		Names []string `help:"Packages to install"`
		Quiet bool     `type:"option" short:"q" help:"Make the output quiet."`
	} `help:"Add a package"`
	Remove struct {
		Names []string `help:"Packages to remove"`
		Quiet bool     `type:"option" short:"q" help:"Make the output quiet."`
	} `help:"Remove a package"`
	Update struct {
		Upgrade struct {
			All bool `type:"option" short:"A" help:"Upgrade all packages"`
		} `help:"Upgrade packages"`
		All bool `type:"option" short:"A" help:"Update all packages"`
	} `help:"Update packages"`
	List bool `type:"command" help:"List installed packages"`
}

func main() {
	args := Args{}
	err := applause.Parse(&args)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(args)
}
