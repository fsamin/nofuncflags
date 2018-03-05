package main

import (
	"go/token"
	"os"
)

var (
	fset *token.FileSet
	dir  string
)

func main() {
	if len(os.Args) == 1 {
		new(".", os.Stdout).parse()
		return
	}

	for _, arg := range os.Args {
		fs, err := os.Stat(arg)
		if err != nil {
			panic(err)
		}
		if fs.IsDir() {
			linter := new(arg, os.Stdout)
			linter.parse()
			linter.print()
		}
	}
}
