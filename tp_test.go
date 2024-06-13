package bgo_test

import (
	"bgo/types"
	"testing"
)

func TestTp(t *testing.T) {
	var im types.Packages = make(types.Packages)
	change(im)
	t.Log(im)
}

func change(input types.Packages) {
	input["main"] = &types.Package{}
}
