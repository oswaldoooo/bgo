package test

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"runtime"
	"strconv"
	"testing"

	bparser "github.com/oswaldoooo/bgo/parser"
	"github.com/oswaldoooo/bgo/types"
)

func TestParse(t *testing.T) {
	af := parse("example.go")
	pkg := types.NewPackages()
	err := bparser.Parse(af, pkg)
	throw(err)
	content, _ := json.MarshalIndent(pkg, "", "    ")
	err = os.WriteFile("test.json", content, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
func TestAst(t *testing.T) {
	fs := token.NewFileSet()
	asf, err := parser.ParseFile(fs, "example.go", nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile("example.ast", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = ast.Fprint(f, fs, asf, nil)
	if err != nil {
		t.Fatal(err)
	}
}
func TestShowOrigin(t *testing.T) {
	fs := token.NewFileSet()
	af, err := parser.ParseFile(fs, "example.go", nil, parser.AllErrors|parser.ParseComments)
	throw(err)
	ast.Print(fs, af)
}
func parse(rpath string) *ast.File {
	fs := token.NewFileSet()
	af, err := parser.ParseFile(fs, rpath, nil, parser.AllErrors|parser.ParseComments)
	throw(err)
	return af
}
func throw(err error) {
	if err != nil {
		_, f, line, _ := runtime.Caller(1)
		fmt.Fprintln(os.Stderr, f+":"+strconv.Itoa(line)+" error "+err.Error())
		os.Exit(-1)
	}
}
