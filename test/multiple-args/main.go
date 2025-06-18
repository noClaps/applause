package main

import (
	"fmt"

	"github.com/noclaps/applause"
)

type Args struct {
	MyArg    string   `help:"This is the help text for my-arg"`
	ManyArgs []string `help:"This takes multiple arguments"`
	MyArg2   string   `help:"This is the help text for my-arg-2"`
}

func main() {
	args := Args{}
	if err := applause.Parse(&args); err != nil {
		fmt.Println(err)
		fmt.Println(applause.Help)
		fmt.Println(applause.Usage)
		return
	}

	fmt.Printf("%v\n", args)
}
