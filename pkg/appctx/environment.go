package appctx

import "context"

type Environment struct {
	ExampleVar string
	AnotherVar string
}

type EnvironmentContext interface {
	Environment() *Environment
	context.Context
}

func (self *Environment) Environment() *Environment {
	return self
}
