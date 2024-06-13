package parser

import (
	"bgo/types"
	"errors"
	"go/ast"
	"go/token"
)

// func init(){
// 	parser.ParseFile()
// }

func Parse(atree *ast.File, dst types.Packages) error {
	pkgname := atree.Name.Name
	if _, ok := dst[pkgname]; !ok {
		dst[pkgname] = types.NewPackage()
	}
	pobj := dst[pkgname]
	var err error
	for _, e := range atree.Decls {
		if ee, ok := e.(*ast.GenDecl); ok {
			//this is like type xxx xxx
			switch ee.Tok {
			case token.TYPE:
				err = parseobj(ee, dst, pobj)
			case token.CONST:
				panic("not implement")
			case token.VAR:
				panic("not implement")
			default:
				return errors.New("unknown keyword " + ee.Tok.String())
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