// Package goparse contains logic for parsing Go files. Specifically it parses
// source and test files into domain models for generating tests.
package cmd

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os/user"
	"strings"
	"time"
)

// ErrEmptyFile represents an empty file error.
var ErrEmptyFile = errors.New("file is empty")

// Result representats a parsed Go file.
type Result struct {
	// The package name and imports of a Go file.
	Header *Header
	// All the functions and methods in a Go file.
	Funcs []*Function
}

// Parser can parse Go files.
type Parser struct {
	// The importer to resolve packages from import paths.
	Importer types.Importer
}

// Parse parses a given Go file at srcPath, along any files that share the same
// package, into a domain model for generating tests.
func (p *Parser) Parse(srcPath string, files []Path) (*Result, error) {
	b, err := p.readFile(srcPath)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	f, err := p.parseFile(fset, srcPath)
	if err != nil {
		return nil, err
	}
	fs, err := p.parseFiles(fset, f, files)
	if err != nil {
		return nil, err
	}
	fmt.Println("srcPath: ", srcPath)
	fmt.Println("f.Package: ", f.Package)
	return &Result{
		Header: &Header{
			Comments: parsePkgComment(f, f.Package),
			Package:  f.Name.String(),
			Imports:  parseImports(f.Imports),
			Code:     goCode(b, f),
		},
		Funcs: p.parseFunctions(fset, f, fs),
	}, nil
}

func (p *Parser) readFile(srcPath string) ([]byte, error) {
	b, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	if len(b) == 0 {
		return nil, ErrEmptyFile
	}
	return b, nil
}

func (p *Parser) parseFile(fset *token.FileSet, srcPath string) (*ast.File, error) {
	f, err := parser.ParseFile(fset, srcPath, nil, parser.ParseComments)
	fmt.Println("test:", len(f.Comments))
	for _, item := range f.Comments {
		fmt.Printf("1: %v\n", item.Text())

		for _, te := range item.List {
			fmt.Printf("2: %v\n", te.Text)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("target parser.ParseFile(): %v", err)
	}

	return f, nil
}

func (p *Parser) parseFiles(fset *token.FileSet, f *ast.File, files []Path) ([]*ast.File, error) {
	pkg := f.Name.String()
	var fs []*ast.File
	for _, file := range files {
		ff, err := parser.ParseFile(fset, string(file), nil, 0)
		if err != nil {
			return nil, fmt.Errorf("other file parser.ParseFile: %v", err)
		}
		if name := ff.Name.String(); name != pkg {
			continue
		}
		fs = append(fs, ff)
	}
	return fs, nil
}

func (p *Parser) parseFunctions(fset *token.FileSet, f *ast.File, fs []*ast.File) []*Function {
	ul, el := p.parseTypes(fset, fs)
	// fmt.Printf("ul: %v\n", ul)
	var funcs []*Function
	for _, d := range f.Decls {
		fDecl, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcs = append(funcs, parseFunc(fDecl, ul, el))
	}

	return funcs
}

func (p *Parser) parseTypes(fset *token.FileSet, fs []*ast.File) (map[string]types.Type, map[*types.Struct]ast.Expr) {
	conf := &types.Config{
		Importer: p.Importer,
		// Adding a NO-OP error function ignores errors and performs best-effort
		// type checking. https://godoc.org/golang.org/x/tools/go/types#Config
		Error: func(error) {},
	}
	ti := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	// Note: conf.Check can fail, but since Info is not required data, it's ok.
	conf.Check("", fset, fs, ti)
	ul := make(map[string]types.Type)
	el := make(map[*types.Struct]ast.Expr)
	for e, t := range ti.Types {
		// Collect the underlying types.
		ul[t.Type.String()] = t.Type.Underlying()
		// Collect structs to determine the fields of a receiver.
		if v, ok := t.Type.(*types.Struct); ok {
			el[v] = e
		}
	}
	return ul, el
}

func parsePkgComment(f *ast.File, pkgPos token.Pos) []string {
	var comments []string
	var count int

	// fmt.Println("f.Comments: ", len(f.Comments))
	// for _, comment := range f.Comments {
	// 	fmt.Println("f.cmt: ", comment.Text())
	// }
	for _, comment := range f.Comments {
		// fmt.Println("f.cmt: ", f.Name.End())
		// var pkg token.Pos(len(f.Name.String())
		if comment.End() <= f.Name.End() {
			// fmt.Println("f.end: ", comment.End())
			// fmt.Println("f.pkgPos: ", pkgPos)
			// fmt.Println("skip")
			break
		}
		for _, c := range comment.List {
			// fmt.Printf("f.c: %s - %v\n", c.Text, c.End())
			count += len(c.Text) + 1 // +1 for '\n'
			if count < int(c.End()) {
				// n := int(c.End()) - count - 1
				comments = append(comments, strings.Repeat("\n", 1))
				count++ // for last of '\n'
			}
			comments = append(comments, c.Text)
		}
	}

	if int(pkgPos)-count > 1 {
		comments = append(comments, strings.Repeat("\n", int(pkgPos)-count-2))
	}
	// fmt.Printf("comments 11: %v\n", comments)
	return comments
}

// Returns the Go code below the imports block.
func goCode(b []byte, f *ast.File) []byte {
	furthestPos := f.Name.End()
	for _, node := range f.Imports {
		if pos := node.End(); pos > furthestPos {
			furthestPos = pos
		}
	}
	if furthestPos < token.Pos(len(b)) {
		furthestPos++

		// Avoid wrong output on windows-encoded files
		if b[furthestPos-2] == '\r' && b[furthestPos-1] == '\n' && furthestPos < token.Pos(len(b)) {
			furthestPos++
		}
	}
	return b[furthestPos:]
}

func parseFunc(fDecl *ast.FuncDecl, ul map[string]types.Type, el map[*types.Struct]ast.Expr) *Function {
	f := &Function{
		Name:       fDecl.Name.String(),
		IsExported: fDecl.Name.IsExported(),
		Receiver:   parseReceiver(fDecl.Recv, ul, el),
		Parameters: parseFieldList(fDecl.Type.Params, ul),
	}
	cmts := make([]string, 0)
	if fDecl.Doc != nil {
		for _, cmt := range fDecl.Doc.List {
			cmts = append(cmts, cmt.Text)
		}
	}
	user, err := user.Current()
	if err == nil {
		cmts = append(cmts, fmt.Sprintf("// Created at %s by %s", time.Now().Format("02-01-2006"), user.Name))
	} else {
		cmts = append(cmts, fmt.Sprintf("// Created at %s", time.Now().Format("02-01-2006")))
	}

	for _, p := range f.Parameters {
		if p.Type != nil && p.Type.String() == "echo.Context" {
			f.IsEcho = true
			break
		}
	}
	f.Comments = cmts

	fs := parseFieldList(fDecl.Type.Results, ul)
	i := 0
	for _, fi := range fs {
		if fi.Type.String() == "error" {
			f.ReturnsError = true
			continue
		}
		fi.Index = i
		f.Results = append(f.Results, fi)
		i++
	}
	return f
}

func parseImports(imps []*ast.ImportSpec) []*Import {
	var is []*Import
	for _, imp := range imps {
		var n string
		if imp.Name != nil {
			n = imp.Name.String()
		}
		is = append(is, &Import{
			Name: n,
			Path: imp.Path.Value,
		})
	}
	return is
}

func parseReceiver(fl *ast.FieldList, ul map[string]types.Type, el map[*types.Struct]ast.Expr) *Receiver {
	if fl == nil {
		return nil
	}
	r := &Receiver{
		Field: parseFieldList(fl, ul)[0],
	}
	t, ok := ul[r.Type.Value]
	if !ok {
		return r
	}
	s, ok := t.(*types.Struct)
	if !ok {
		return r
	}
	st, found := el[s]
	if !found {
		return r
	}
	r.Fields = append(r.Fields, parseFieldList(st.(*ast.StructType).Fields, ul)...)
	for i, f := range r.Fields {
		// https://github.com/cweill/gotests/issues/69
		if i >= s.NumFields() {
			break
		}
		f.Name = s.Field(i).Name()
	}
	return r

}

func parseFieldList(fl *ast.FieldList, ul map[string]types.Type) []*Field {
	if fl == nil {
		return nil
	}
	i := 0
	var fs []*Field
	for _, f := range fl.List {
		for _, pf := range parseFields(f, ul) {
			pf.Index = i
			fs = append(fs, pf)
			i++
		}
	}
	return fs
}

func parseFields(f *ast.Field, ul map[string]types.Type) []*Field {
	t := parseExpr(f.Type, ul)
	if len(f.Names) == 0 {
		return []*Field{{
			Type: t,
		}}
	}
	var fs []*Field
	for _, n := range f.Names {
		fs = append(fs, &Field{
			Name: n.Name,
			Type: t,
		})
	}
	return fs
}

func parseExpr(e ast.Expr, ul map[string]types.Type) *Expression {
	switch v := e.(type) {
	case *ast.StarExpr:
		val := types.ExprString(v.X)
		return &Expression{
			Value:      val,
			IsStar:     true,
			Underlying: underlying(val, ul),
		}
	case *ast.Ellipsis:
		exp := parseExpr(v.Elt, ul)
		return &Expression{
			Value:      exp.Value,
			IsStar:     exp.IsStar,
			IsVariadic: true,
			Underlying: underlying(exp.Value, ul),
		}
	default:
		val := types.ExprString(e)
		return &Expression{
			Value:      val,
			Underlying: underlying(val, ul),
			IsWriter:   val == "io.Writer",
		}
	}
}

func underlying(val string, ul map[string]types.Type) string {
	if ul[val] != nil {
		return ul[val].String()
	}
	return ""
}
