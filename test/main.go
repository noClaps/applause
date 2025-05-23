package main

import (
	"fmt"

	"github.com/noclaps/applause"
)

type Args struct {
	MyArg  string `help:"This is the help text for my-arg"`
	MyArg2 string `help:"This is the help text for my-arg-2"`
	Opt1   int    `type:"option" short:"o" value:"option" help:"This is the help text for opt-1"`
	Opt2   bool   `type:"option" short:"p" help:"This is the help text for opt-2"`
}

func main() {
	args := Args{MyArg: "test-val", Opt1: 5}
	if err := applause.Parse(&args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("MyArg: %s, MyArg2: %s, Opt1: %d, Opt2: %v\n", args.MyArg, args.MyArg2, args.Opt1, args.Opt2)
}
