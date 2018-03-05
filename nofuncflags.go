package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"path/filepath"
)

type nofuncflags struct {
	dir string
	out io.Writer
	res []funcParser
}

func new(dir string, out io.Writer) *nofuncflags {
	return &nofuncflags{
		dir: dir,
		out: out,
	}
}

func (n *nofuncflags) print() {
	for _, f := range n.res {
		f.print()
	}
}

func (n *nofuncflags) parse() {
	fset = token.NewFileSet()
	f, err := parser.ParseDir(fset, n.dir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, pkg := range f {
		for _, file := range pkg.Files {
			for _, d := range file.Decls {
				if funcDecl, ok := d.(*ast.FuncDecl); ok {
					fParser := funcParser{
						decl: funcDecl,
					}
					if hasFlags(fParser) {
						n.res = append(n.res, fParser)
					}
				}
			}
		}
	}
}

type funcParser struct {
	tokenPos token.Position
	decl     *ast.FuncDecl
}

func (f funcParser) print() {
	fmt.Printf("%s:%d:%d:%s should not take a boolean as parameter\n", f.path(), f.pos().Line, f.pos().Column, f.getName())
}

func (f funcParser) pos() *token.Position {
	if !f.tokenPos.IsValid() {
		f.tokenPos = fset.Position(f.decl.Pos())
	}
	return &f.tokenPos
}
func (f funcParser) path() string {
	path := filepath.Join(f.pos().Filename)
	return path
}
func (f funcParser) getName() string {
	return f.decl.Name.Name
}

func hasFlags(f funcParser) bool {
	fList := f.decl.Type.Params
	expr := typeFlatten(fList.List)
	for _, e := range expr {
		if ident, ok := e.(*ast.Ident); ok {
			if ident.Name == "bool" {
				return true
			}
		}
	}
	return false
}

// Turn parameter list into slice of types
// (in the ast, types are Exprs).
// Have to handle f(int, bool) and f(x, y, z int)
// so not a simple 1-to-1 conversion.
func typeFlatten(l []*ast.Field) []ast.Expr {
	var t []ast.Expr
	for _, f := range l {
		if len(f.Names) == 0 {
			t = append(t, f.Type)
			continue
		}
		for _ = range f.Names {
			t = append(t, f.Type)
		}
	}
	return t
}
