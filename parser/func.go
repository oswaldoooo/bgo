package parser

import (
	"go/ast"

	"github.com/oswaldoooo/bgo/internal/utils"
	"github.com/oswaldoooo/bgo/types"
)

func parseFunc(src *ast.FuncDecl, packages *types.Packages, currpkg *types.Package) error {
	var ft types.Func
	ft.Kind = types.FuncType
	ft.Name = src.Name.Name
	if src.Doc != nil && len(src.Doc.List) > 0 {
		ft.Comment = make(types.Comment, len(src.Doc.List))
		utils.SliceConvert(src.Doc.List, ft.Comment, commentparse)
	}
	if src.Recv != nil && len(src.Recv.List) > 0 {
		//parse self
		expr := src.Recv.List[0]
		if len(expr.Names) > 0 {
			ft.Self = expr.Names[0].Name + " "
		}
		if exp, ok := expr.Type.(*ast.Ident); ok {
			ft.Self += exp.Name
		} else if exp, ok := expr.Type.(*ast.StarExpr); ok {
			ft.Self += "*" + exp.X.(*ast.Ident).Name
		} else {
			println("warning parse func type failed,unknown type")
		}
	}
	if src.Type != nil {
		if src.Type.Params != nil && len(src.Type.Params.List) > 0 {
			//parse params
			ft.Params = make([]string, len(src.Type.Params.List))
			utils.SliceConvert(src.Type.Params.List, ft.Params, funcparamsParse)
		}
		if src.Type.Results != nil && len(src.Type.Results.List) > 0 {
			//parse results
			ft.Resp = make([]string, len(src.Type.Results.List))
			utils.SliceConvert(src.Type.Results.List, ft.Resp, funcparamsParse)
		}
	}
	if src.Doc != nil && len(src.Doc.List) > 0 {
		//parse comment
		ft.Comment = make(types.Comment, len(src.Doc.List))
		utils.SliceConvert(src.Doc.List, ft.Comment, commentparse)
	}
	currpkg.Func[ft.Self+":"+ft.Name] = ft
	return nil
}
func funcparamsParse(src *ast.Field, dst *string) {
	if len(src.Names) > 0 {
		*dst = src.Names[0].Name + ":"
	}
	//make it easy the type transfer
	*dst += getExprStr(src.Type)
	// if ident, ok := src.Type.(*ast.Ident); ok {
	// 	*dst += ident.Name
	// } else if ident, ok := src.Type.(*ast.StarExpr); ok {
	// 	//todo: not only has ident and has index expr like xxx.xxxxx
	// 	// *dst += "*" + ident.X.(*ast.Ident).Name
	// 	if idd, ok := ident.X.(*ast.Ident); ok {
	// 		*dst = "*" + idd.Name
	// 	} else if idd, ok := ident.X.(*ast.IndexExpr); ok {
	// 		iddd := idd.Index.(*ast.Ident)
	// 		if sex, ok := idd.X.(*ast.SelectorExpr); ok {
	// 			*dst = "*" + sex.X.(*ast.Ident).Name + "." + sex.Sel.Name + "[" + iddd.Name + "]"
	// 		}

	// 	}
	// } else {
	// 	println("warning parse func param type failed")
	// }
}
func commentparse(src *ast.Comment, dst *string) {
	*dst = src.Text
}
