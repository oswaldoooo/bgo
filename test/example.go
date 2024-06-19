package test

import (
	"context"

	"github.com/oswaldoooo/bgo/types"
)

// gen_controller
type Greeter struct {
	Ctx      context.Context
	ModuleId int32
}
type X interface {
	X() int32
}

// // this is validator
// type Validator interface {
// 	X
// 	Validate(any) error
// }

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

type Greeter2 struct{}

func GetMyName() *types.Group[string] {
	return &types.Group[string]{}
}

//template test

type Result[T, E any] struct {
	T T
	E E
}
type PrefixKey[T any] struct {
	Key T
}

var result Result[string, PrefixKey[int]]
