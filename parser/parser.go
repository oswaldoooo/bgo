package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"

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

type _selector ast.SelectorExpr

func (s _selector) String() (content string) {
	content += s.X.(*ast.Ident).Name
	content += "." + s.Sel.String()
	return
}

type _star ast.StarExpr

func (s _star) String() (content string) {
	content = "*" + getExprStr(s.X)
	// if id, ok := s.X.(*ast.Ident); ok {
	// 	content += id.Name
	// 	return
	// } else if id, ok := s.X.(*ast.SelectorExpr); ok {
	// 	content += _selector(*id).String()
	// } else if id, ok := s.X.(*ast.IndexExpr); ok {
	// 	content += _index(*id).String()
	// 	// id := s.X.(*ast.IndexExpr)
	// 	// xxx[xxx]
	// 	if idd, ok := id.X.(*ast.Ident); ok {
	// 		content += idd.Name
	// 	} else if idd, ok := id.X.(*ast.SelectorExpr); ok {
	// 		// xxx.xxxx[xxxx]
	// 		content += _selector(*idd).String()
	// 	} else {
	// 		panic("can't parse ast tree")
	// 	}
	// 	content += "["
	// 	if idd, ok := id.Index.(*ast.Ident); ok {
	// 		content += idd.Name
	// 	} else if idd, ok := id.Index.(*ast.SelectorExpr); ok {
	// 		content += _selector(*idd).String()
	// 	} else if idd, ok := id.Index.(*ast.IndexExpr); ok {
	// 		content += _index(*idd).String()
	// 	}
	// 	content += "]"

	// }
	return
}

type _index ast.IndexExpr

func (i _index) String() (content string) {
	content = getExprStr(i.X) + "[" + getExprStr(i.Index) + "]"
	return
}

func getExprStr(a ast.Expr) (content string) {
	if id, ok := a.(*ast.Ident); ok {
		content = id.Name
	} else if id, ok := a.(*ast.SelectorExpr); ok {
		content = _selector(*id).String()
	} else if id, ok := a.(*ast.IndexExpr); ok {
		content = _index(*id).String()
	} else if id, ok := a.(*ast.StarExpr); ok {
		content = _star(*id).String()
	} else {
		fmt.Fprintln(os.Stderr, "don't know expr type")
	}
	return
}
