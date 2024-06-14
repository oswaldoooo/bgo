package parser

import (
	"bgo/internal/utils"
	"bgo/types"
	"errors"
	"fmt"
	"go/ast"
	"strings"
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
		dst.Kind = types.FieldType
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
		if src.Tag != nil {
			dst.Tag = src.Tag.Value
		}
		if src.Doc != nil && len(src.Doc.List) > 0 {
			dst.Comment = make(types.Comment, len(src.Doc.List))
			utils.SliceConvert(src.Doc.List, dst.Comment, func(src *ast.Comment, dst *string) {
				*dst = strings.Trim(src.Text, " ")
			})
		}
	})

	_dst.Fields = fieldlist
}

// parse variable
func parsevar(src *ast.GenDecl, packages types.Packages, currpkg *types.Package) error {
	if len(src.Specs) == 0 {
		return nil
	}
	for _, spec := range src.Specs {
		var vtp types.Variable
		vtp.Kind = types.VariableType
		vspec := spec.(*ast.ValueSpec)
		vtp.Name = vspec.Names[0].Name
		if exp, ok := vspec.Type.(*ast.Ident); ok {
			vtp.Name += ":" + exp.Name
		} else if exp, ok := vspec.Type.(*ast.StarExpr); ok {
			vtp.Name += ":*" + exp.X.(*ast.Ident).Name
		} else {
			println("warning parse variable type failed")
		}
		//parse value
		vtp.Value = parsevalue(vspec.Values)
		//parse comment
		if vspec.Doc != nil && len(vspec.Doc.List) > 0 {
			vtp.Comment = make(types.Comment, len(vspec.Doc.List))
			utils.SliceConvert(vspec.Doc.List, vtp.Comment, commentparse)
		}
		currpkg.Variables[getRawName(vtp.Name)] = vtp
	}
	return nil
}

// parse const
func parseconst(src *ast.GenDecl, packages types.Packages, currpkg *types.Package) error {
	if len(src.Specs) == 0 {
		return nil
	}
	for _, spec := range src.Specs {
		var vtp types.Const
		vtp.Kind = types.ConstType
		vspec := spec.(*ast.ValueSpec)
		vtp.Name = vspec.Names[0].Name
		if exp, ok := vspec.Type.(*ast.Ident); ok {
			vtp.Name += ":" + exp.Name
		} else if exp, ok := vspec.Type.(*ast.StarExpr); ok {
			vtp.Name += ":*" + exp.X.(*ast.Ident).Name
		} else {
			println("warning parse const type failed")
		}
		//parse value
		vtp.Value = parsevalue(vspec.Values)
		//parse comment
		if vspec.Doc != nil && len(vspec.Doc.List) > 0 {
			vtp.Comment = make(types.Comment, len(vspec.Doc.List))
			utils.SliceConvert(vspec.Doc.List, vtp.Comment, commentparse)
		}
		currpkg.Const[getRawName(vtp.Name)] = vtp
	}
	return nil
}

func getRawName(s string) string {
	if index := strings.IndexByte(s, ':'); index > 0 {
		return s[:index]
	}
	return s
}

// parse value
func parsevalue(src []ast.Expr) (result string) {
	if len(src) == 0 {
		return
	}
	//todo: support conflict expr
	for _, v := range src {
		if p, ok := v.(*ast.BasicLit); ok {
			result = p.Value
			return
		} else if p, ok := v.(*ast.Ident); ok {
			result = p.Name
			return
		}
		continue

	}
	return
}
