# Applause

A Go command line argument parsing library. Named after and inspired by [clap-rs/clap](https://github.com/clap-rs/clap).

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
USAGE: program <my-arg> <my-arg-2> [--opt-1 <option>] [--opt-2]

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

You can also get the usage and help strings by using `applause.Usage` and `applause.Help`, after you have called `applause.Parse()`:

```go
// ...

func main() {
	args := Args{}
	_ = applause.Parse(&args)
	help := applause.Help
	usage := applause.Usage

	fmt.Println("help:\n", help)
	fmt.Println()
	fmt.Println("usage:\n", usage)
}
```

This will output:

```
help:
USAGE: program <my-arg> <my-arg-2> [--opt-1 <option>] [--opt-2]

ARGUMENTS:
  <my-arg>                    This is the help text for my-arg
  <my-arg-2>                  This is the help text for my-arg-2

OPTIONS:
  -o, --opt-1 <option>        This is the help text for opt-1 (default: 5)
  -p, --opt-2                 This is the help text for opt-2
  -h, --help                  Display this help and exit.

usage:
USAGE: program <my-arg> <my-arg-2> [--opt-1 <option>] [--opt-2]
```

## Configuration

The configuration struct should have fields with types and some struct tags. All fields you'd like to be parsed should be exported in the struct.

The supported types for fields are:

- `bool`
- `float32`, `float64`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `complex64`, `complex128`
- `string`

Each field should have some struct tags:

- `type`: The type can be `"arg"`, `"option"` or `"command"`. If omitted, the default is `"arg"`. If any other type is provided, the field is ignored. Example:

  ```go
  type Args struct {
    Arg string `type:"arg"`
    Arg2 string // implicit `type:"arg"`
    Option bool `type:"option"`
  }
  ```

- `name`: The name of the argument, option or command. If omitted, the default is the field name converted to kebab case. If you'd like to have an option have a different name, you can write it as `name:"option-name"` in the tags. Example:

  ```go
  type Args struct {
    Arg string `name:"my-arg"` // will be displayed as <my-arg> in the help text
    Option bool `type:"option" name:"my-option"` // will be called with --my-option
    Flag bool `type:"option"` // will be called with --flag
  }
  ```

  You can also set `name:""` to use only the short form. If no `short` tag is set then the option will be inaccessible and will not appear in the help menu.

- `help`: The help text for the argument, option or command, will be displayed in the command help when the command is called with `--help` or `-h`. Example:

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

### Multiple arguments

If you'd like an argument to take multiple values, you can use a slice. The supported types are:

- `[]bool`
- `[]float32`, `[]float64`
- `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`
- `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`
- `[]complex64`, `[]complex128`
- `[]string`

For example, if you have:

```go
package main

import (
	"fmt"
)

type Args struct {
	First string `help:"First argument"`
	Multiple []string `help:"Multiple arguments"`
	Last string `help:"Last argument"`
}

func main() {
	args := Args{}
	_ = applause.Parse(&args)

	fmt.Println(args.First, args.Multiple, args.Last)
}
```

and you run:

```sh
./program first second third fourth fifth
```

your output will be:

```
first [second third fourth] fifth
```

Note that the argument that takes multiple values is optional, since the slice can simply be empty:

```
USAGE: program <first> [multiple...] <last>

ARGUMENTS:
  <first>              First argument
  [multiple...]        Multiple arguments
  <last>               Last argument

OPTIONS:
  -h, --help           Display this help and exit.
```

If you run:

```sh
./program first second
```

your output will be:

```
first [] second
```

### Default options

You can also set default values for options and they will appear in the help menu. For example, if you have:

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
	args := Args{Opt1: 5}
	applause.Parse(&args)
}
```

and you run `./program --help`, you'll see:

```
USAGE: program <my-arg> <my-arg-2> [--opt-1 <option>] [--opt-2]

ARGUMENTS:
  <my-arg>                    This is the help text for my-arg
  <my-arg-2>                  This is the help text for my-arg-2

OPTIONS:
  -o, --opt-1 <option>        This is the help text for opt-1 (default: 5)
  -p, --opt-2                 This is the help text for opt-2
  -h, --help                  Display this help and exit.
```

### Commands

You can define commands by using a struct as the field type:

```go
package main

import (
	"github.com/noclaps/applause"
)

type Args struct {
	Add struct {
		Names []string `help:"Packages to install"`
		Quiet bool     `type:"option" short:"q" help:"Make the output quiet."`
	} `help:"Add a package"`
}

func main() {
	args := Args{}
	_ = applause.Parse(args)

	args.Add.Names
}
```

If you run `./program --help` on this, you'll get:

```
USAGE: program [add]

COMMANDS:
  add               Add a package

OPTIONS:
  -h, --help        Display this help and exit.
```

You can also run `./program add --help` to get the help menu for the command:

```
USAGE: program add [names...] [--quiet]

ARGUMENTS:
  [names...]           Packages to install

OPTIONS:
  -q, --quiet          Make the output quiet.
  -h, --help           Display this help and exit.
```

You can add subcommands by nesting structs:

```go
package main

import (
	"github.com/noclaps/applause"
)

type Args struct {
	Update struct {
		Upgrade struct {
			All bool `type:"option" short:"A" help:"Upgrade all packages"`
		} `help:"Upgrade packages"`
		All bool `type:"option" short:"A" help:"Update all packages"`
	} `help:"Update packages"`
}

func main() {
	args := Args{}
	_ = applause.Parse(&args)

	args
}
```

This will allow you to run subcommands like `./program update upgrade`.

You can also have commands without any arguments using a boolean and the `type:"command"` tag:

```go
package main

import (
	"github.com/noclaps/applause"
)

type Args struct {
	List bool `type:"command" help:"List files"`
}

func main() {
	args := Args{}
	_ = applause.Parse(&args)

	args
}
```

Running `./program list` will set `args.List` to `true`.
