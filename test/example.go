package test

import "context"

// gen_controller
type Greeter struct {
	Ctx      context.Context
	ModuleId int32
}

// gen_enum
type Module int32

const (
	Start  Module = 1
	End    Module = 2
	Failed Module = 3
	Panic  Module = 4
)

var (
	//start
	//
	// default_gen
	DefaultModule = Start
)
