package bgo

import (
	stdparser "go/parser"
	"go/token"

	"github.com/oswaldoooo/bgo/parser"
	"github.com/oswaldoooo/bgo/types"
)

func Parse(rpath string) (types.Raw, types.Packages, error) {
	fs := token.NewFileSet()
	asf, err := stdparser.ParseFile(fs, rpath, nil, stdparser.ParseComments|stdparser.AllErrors)
	if err != nil {
		return types.Raw{}, nil, err
	}
	ans := make(types.Packages)
	err = parser.Parse(asf, ans)
	return types.Raw{
		Fs:     fs,
		Astree: asf,
	}, ans, err
}
