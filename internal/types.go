package internal

type Config any

type Arg struct {
	Name string
	Help string
}

type Option struct {
	Name  string
	Value any
	Help  string
	Short string
}
