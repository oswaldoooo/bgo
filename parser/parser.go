package parser

import (
	"go/ast"
	"go/token"

	"github.com/oswaldoooo/bgo/types"
)

// func init(){
// 	parser.ParseFile()
// }

func Parse(atree *ast.File, dst *types.Packages) error {
	pkgname := atree.Name.Name
	if _, ok := dst.Pkgs[pkgname]; !ok {
		dst.Pkgs[pkgname] = types.NewPackage()
	}
	pobj := dst.Pkgs[pkgname]
	var err error
	for _, e := range atree.Decls {
		if ee, ok := e.(*ast.GenDecl); ok {
			//this is like type xxx xxx
			switch ee.Tok {
			case token.TYPE:
				err = parseobj(ee, dst, pobj)
			case token.CONST:
				err = parseconst(ee, dst, pobj)
			case token.VAR:
				err = parsevar(ee, dst, pobj)
			default:
				println("unknown keyword " + ee.Tok.String())
			}
			if err != nil {
				return err
			}
			continue
		}
		if ee, ok := e.(*ast.FuncDecl); ok {
			//this is function
			err = parseFunc(ee, dst, pobj)
			if err != nil {
				return err
			}
			continue
		}

	}
	return nil
}
