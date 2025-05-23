# Applause

A Go command line argument parsing library. Named after [clap-rs/clap](https://github.com/clap-rs/clap), and inspired by [clap-rs/clap](https://github.com/clap-rs/clap) and [apple/swift-argument-parser](https://github.com/apple/swift-argument-parser).

## Usage

Add the library as a dependency:

```sh
go get github.com/noclaps/applause
```

Create a struct to define your arguments, and pass in a pointer to the struct as the argument to the `applause.Parse()` function:

```go
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
	args := Args{}
	if err := applause.Parse(&args); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("MyArg: %s, MyArg2: %s, Opt1: %d, Opt2: %v\n", args.MyArg, args.MyArg2, args.Opt1, args.Opt2)
}
```

This will generate the CLI tool for you. If you run `go build -o program` and then `./program --help`, you'll see:

```
USAGE: ./program <my-arg> <my-arg-2> [--opt-1 <option>] [--opt-2]

ARGUMENTS:
  <my-arg>                    This is the help text for my-arg
  <my-arg-2>                  This is the help text for my-arg-2

OPTIONS:
  -o, --opt-1 <option>        This is the help text for opt-1
  -p, --opt-2                 This is the help text for opt-2
  -h, --help                  Display this help and exit.
```

This is automatically generated for you, and if you were to run `./program hello hi -o 5 -p`, the program would print:

```
MyArg: hello, MyArg2: hi, Opt1: 5, Opt2: true
```

The values parsed from the command line arguments are put back into the `args` struct, and can be used as needed from there.

## Configuration

The configuration struct should have fields with types and some struct tags.

The supported types for fields are:

- `bool`
- `float32`, `float64`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `complex64`, `complex128`
- `string`

Each field should have some struct tags:

- `type`: The type can be `"arg"` or `"option"`. If omitted, the default is `"arg"`. If any other type is provided, the field is ignored. Example:

  ```go
  type Args struct {
    Arg string `type:"arg"`
    Arg2 string // implicit `type:"arg"`
    Option bool `type:"option"`
  }
  ```

- `name`: The name of the argument or option. If omitted, the default is the field name converted to kebab case. If you'd like to have an option have a different name, you can write it as `name:"option-name"` in the tags. Example:

  ```go
  type Args struct {
    Arg string `name:"my-arg"` // will be displayed as <my-arg> in the help text
    Option bool `type:"option" name:"my-option"` // will be called with --my-option
    Flag bool `type:"option"` // will be called with --flag
  }
  ```

- `help`: The help text for the argument or option, will be displayed in the command help when the command is called with `--help` or `-h`. Example:

  ```go
  type Args struct {
    Arg string `help:"This is my argument"` // <arg>        This is my argument
    Option bool `type:"option" help:"This is my option"` // --option        This is my option
    Option2 bool `type:"option" name:"opt" help:"This is another option"` // --opt        This is another option
  }
  ```

- `short`: Only applicable when `type` is "option". The short form of the option. For instance, if you have a field with the tag `name:"option" short:"o"`, you can call the command with `--option` or `-o`. Example:

  ```go
  type Args struct {
    Flag  bool `type:"option" short:"f"` // can be called with --flag or -f
    Opt   int  `type:"option" name:"option" short:"o"` // can be called with --option <opt> or -o <opt>
  }
  ```

- `value`: Only applicable when `type` is `"option"` and the field type is not `"bool"`. The name of the option value to be displayed in the help text. For instance, `name:"option" value:"val"` will be displayed as `--option <val>` in the help text. Example:

  ```go
  type Args struct {
    Option int `type:"option" value:"my-opt"` // --option <my-opt>
    Flag bool `type:"option" value:"my-flag"` // --flag, value is not applicable when the field type is "bool"
  }
  ```
