package main

import (
	"log"

	"github.com/noclaps/applause"
)

type Args struct {
	Add struct {
		Packages []string `help:"Packages to install." completion:" files"`
		File     string   `type:"option" short:"f" help:"Install from a file" completion:"files[*.{json,jsonc}]"`
	} `help:"Install packages."`
	Update *struct {
		Packages []string `help:"Packages to update." completion:"bun go lazygit pkg"`
	} `help:"Update packages."`
	Remove struct {
		Packages []string `help:"Packages to remove." completion:"$(jq -r 'keys[]' $PKG_HOME/pkg.lock | tr '\n' ' ')"`
	} `help:"Remove packages."`
	Info struct {
		Package string `help:"The package to get the info for" completion:"$(jq -r 'keys[]' $PKG_HOME/pkg.lock | tr '\n' ' ')"`
	} `help:"Get the info for a package."`
	List string `type:"option" help:"List installed packages" completion:"installed remote"`
	Init bool   `type:"option" help:"Initialise pkg"`
}

func main() {
	args := Args{}
	if err := applause.Parse(&args); err != nil {
		log.Fatalln(err)
	}
}
