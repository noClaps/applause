package internal

type Config any

type Arg struct {
	Help string
}

type Option struct {
	Value any
	Help  string
	Short string
}
