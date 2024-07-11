package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/oswaldoooo/bgo/internal/utils"
	"github.com/oswaldoooo/bgo/types"
)

func parseobj(ctx *context, src *ast.GenDecl, packages *types.Packages, currpkg *types.Package) error {
	if len(src.Specs) == 0 {
		fmt.Println("warning: not found spec at ", src.TokPos)
		return nil
	}
	stype, ok := src.Specs[0].(*ast.TypeSpec)
	if !ok {
		fmt.Println("warning: spec 0 to type spec failed")
		return nil
	}
	ctx.set("start", int64(src.Pos())-1)
	ctx.set("end", int64(src.End())-1)
	var structtp types.Struct
	structtp.Name = stype.Name.Name
	structtp.Kind = types.StructType
	if src.Doc != nil && len(src.Doc.List) > 0 {
		structtp.Comment = make(types.Comment, len(src.Doc.List))
		utils.SliceConvert(src.Doc.List, structtp.Comment, commentparse)
	}
	//fill
	if structType, ok := stype.Type.(*ast.StructType); ok {
		parsestruct(ctx, structType, &structtp)
	} else if iftype, ok := stype.Type.(*ast.InterfaceType); ok {
		//support for interface type
		var dsttype types.Interface
		dsttype.Kind = types.InterfaceType
		dsttype.Name = stype.Name.Name
		dsttype.Comment = structtp.Comment
		parseInterface(iftype, &dsttype)
		currpkg.Interface[dsttype.Name] = dsttype
		return nil
	} else {
		ident, ok := stype.Type.(*ast.Ident)
		if !ok {
			ast.Print(nil, stype)
			return errors.New(fmt.Sprint("can't parse ", src.TokPos))
		}
		structtp.Ident = ident.Name
	}
	//
	currpkg.Struct[structtp.Name] = structtp
	return nil
}
func parsestruct(ctx *context, input *ast.StructType, _dst *types.Struct) {
	_dst.Start = ctx.get("start").(int64)
	_dst.End = ctx.get("end").(int64)
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
func parsevar(ctx *context, src *ast.GenDecl, packages *types.Packages, currpkg *types.Package) error {
	if len(src.Specs) == 0 {
		return nil
	}
	//todo: add to var group if count >1
	var (
		add2group = len(src.Specs) > 1
		group     types.Group[[]types.Variable]
	)
	if add2group {
		if src.Doc != nil && len(src.Doc.List) > 0 {
			group.Comments = make(types.Comment, len(src.Doc.List))
			utils.SliceConvert(src.Doc.List, group.Comments, commentparse)
		}
	}
	for _, spec := range src.Specs {
		var vtp types.Variable
		vtp.Kind = types.VariableType
		vspec := spec.(*ast.ValueSpec)
		vtp.Name = vspec.Names[0].Name
		fmt.Println("parse variable " + vtp.Name)
		tpstr := getExprStr(vspec.Type)
		if len(tpstr) > 0 {
			vtp.Name += ":" + tpstr
		}
		//parse value
		vtp.Value = parsevalue(vspec.Names, vspec.Values)
		//parse comment
		if vspec.Doc != nil && len(vspec.Doc.List) > 0 {
			vtp.Comment = make(types.Comment, len(vspec.Doc.List))
			utils.SliceConvert(vspec.Doc.List, vtp.Comment, commentparse)
		}
		vtp.Start = int64(spec.Pos())
		vtp.End = int64(spec.End())
		currpkg.Variables[getRawName(vtp.Name)] = vtp
		group.Members = append(group.Members, vtp)
	}
	if add2group {
		packages.VarGroups = append(packages.VarGroups, group)
	}
	return nil
}

// parse const
func parseconst(ctx *context, src *ast.GenDecl, packages *types.Packages, currpkg *types.Package) error {
	if len(src.Specs) == 0 {
		return nil
	}
	var (
		last_type string
	)
	//todo: add to const group if count >1
	var (
		add2group = len(src.Specs) > 1
		group     types.Group[[]types.Const]
	)
	if add2group {
		if src.Doc != nil && len(src.Doc.List) > 0 {
			group.Comments = make(types.Comment, len(src.Doc.List))
			utils.SliceConvert(src.Doc.List, group.Comments, commentparse)
		}
	}
	for _, spec := range src.Specs {
		var vtp types.Const
		vtp.Start = int64(spec.Pos())
		vtp.End = int64(spec.End())
		vtp.Kind = types.ConstType
		vspec := spec.(*ast.ValueSpec)
		vtp.Name = vspec.Names[0].Name
		thistype := ""
		if exp, ok := vspec.Type.(*ast.Ident); ok {
			thistype = exp.Name
		} else if exp, ok := vspec.Type.(*ast.StarExpr); ok {
			thistype = "*" + exp.X.(*ast.Ident).Name
		} else {
			println("warning parse const type failed")
		}
		if len(thistype) > 0 {
			last_type = thistype
		}
		if len(last_type) > 0 {
			vtp.Name += ":" + last_type
		}

		//parse value
		vtp.Value = parsevalue(vspec.Names, vspec.Values)
		//parse comment
		if vspec.Doc != nil && len(vspec.Doc.List) > 0 {
			vtp.Comment = make(types.Comment, len(vspec.Doc.List))
			utils.SliceConvert(vspec.Doc.List, vtp.Comment, commentparse)
		}
		currpkg.Const[getRawName(vtp.Name)] = vtp
		if add2group {
			group.Members = append(group.Members, vtp)
		}
	}
	if add2group {
		packages.ConstGroup = append(packages.ConstGroup, group)
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
func parsevalue(stp []*ast.Ident, src []ast.Expr) (result string) {
	if src == nil || len(src) == 0 {
		//try get data from obj
		if stp[0].Obj != nil {
			result = fmt.Sprint(stp[0].Obj.Data)
		}
		return
	}
	defer func() {
		if strings.Contains(result, "iota") {
			result = fmt.Sprint(stp[0].Obj.Data)
		}
	}()
	//todo: support conflict expr
	v := src[0]
	// for _, v := range src {
	if p, ok := v.(*ast.BasicLit); ok {
		result = p.Value
		return
	} else if p, ok := v.(*ast.Ident); ok {
		result = p.Name
		return
	}
	// continue

	// }
	return
}

// parse function info to interface struct
func parseInterface(input *ast.InterfaceType, _dst *types.Interface) {
	if input.Methods != nil && input.Methods.List != nil && len(input.Methods.List) > 0 {
		_dst.Funcs = make([]types.Func, len(input.Methods.List))
		utils.SliceConvert(input.Methods.List, _dst.Funcs, func(src *ast.Field, dst *types.Func) {
			if src.Doc != nil && len(src.Doc.List) > 0 {
				dst.Comment = make(types.Comment, len(src.Doc.List))
				utils.SliceConvert(src.Doc.List, dst.Comment, commentparse)
			}
			dst.Kind = types.FuncType
			dst.Name = src.Names[0].Name
			if ftp, ok := src.Type.(*ast.FuncType); ok {
				if ftp.Params != nil && len(ftp.Params.List) > 0 {
					dst.Params = make([]string, len(ftp.Params.List))
					utils.SliceConvert(ftp.Params.List, dst.Params, funcparamsParse)
				}
				if ftp.Results != nil && len(ftp.Results.List) > 0 {
					dst.Resp = make([]string, len(ftp.Results.List))
					utils.SliceConvert(ftp.Results.List, dst.Resp, funcparamsParse)
				}
			}
		})

	}
}
