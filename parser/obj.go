package parser

import (
	"bgo/internal/utils"
	"bgo/types"
	"errors"
	"fmt"
	"go/ast"
)

func parseobj(src *ast.GenDecl, packages types.Packages, currpkg *types.Package) error {
	if len(src.Specs) == 0 {
		fmt.Println("warning: not found spec at ", src.TokPos)
		return nil
	}
	stype, ok := src.Specs[0].(*ast.TypeSpec)
	if !ok {
		fmt.Println("warning: spec 0 to type spec failed")
		return nil
	}
	var structtp types.Struct
	structtp.Name = stype.Name.Name
	structtp.Kind = types.StructType
	//fill
	structType, ok := stype.Type.(*ast.StructType)
	if ok {
		parsestruct(structType, &structtp)
	} else {
		ident, ok := stype.Type.(*ast.Ident)
		if !ok {
			return errors.New(fmt.Sprint("can't parse ", src.TokPos))
		}
		structtp.Ident = ident.Name
	}
	//
	currpkg.Struct[structtp.Name] = structtp
	return nil
}
func parsestruct(input *ast.StructType, _dst *types.Struct) {
	fieldlist := make([]types.Field, len(input.Fields.List))
	utils.SliceConvert(input.Fields.List, fieldlist, func(src *ast.Field, dst *types.Field) {
		dst.Name = src.Names[0].Name
		ident, ok := src.Type.(*ast.Ident)
		if ok {
			dst.Name += ":" + ident.Name
		} else {
			se, ok := src.Type.(*ast.SelectorExpr)
			if !ok {
				fmt.Println("warning: not found suitable express")
				return
			}
			dst.Name += ":" + se.X.(*ast.Ident).Name + "." + se.Sel.Name
		}
	})
	_dst.Fields = fieldlist
}
