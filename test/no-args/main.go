package main

import (
	"fmt"

	"github.com/noclaps/applause"
)

type Args struct {
	Opt1 int  `type:"option" short:"o" value:"option" help:"This is the help text for opt-1"`
	Opt2 bool `type:"option" short:"p" help:"This is the help text for opt-2"`
}

func main() {
	args := Args{Opt1: 5}
	if err := applause.Parse(&args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", args)
}
