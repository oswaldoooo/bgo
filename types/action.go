package types

import "unsafe"

func To[T any](t *Type) *T {
	return (*T)(unsafe.Pointer(t))
}
func NewPackage() *Package {
	return &Package{
		Variables: make(map[string]Variable),
		Const:     make(map[string]Const),
		Func:      make(map[string]Func),
		Struct:    make(map[string]Struct),
		Interface: make(map[string]Interface),
	}
}
