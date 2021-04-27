package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
)

const hello = `package main

import "fmt"

func main() {
        fmt.Println("Hello, world")
}`

func main() {
	fset := token.NewFileSet()

	f, _ := parser.ParseFile(fset, "hello.go", hello, 0)

	conf := types.Config{Importer: importer.Default()}
	pkg, _ := conf.Check("cmd/hello", fset, []*ast.File{f}, nil)

	fmt.Printf("Package  %q\n", pkg.Path())
	fmt.Printf("Name:    %s\n", pkg.Name())
	fmt.Printf("Imports: %s\n", pkg.Imports())
	fmt.Printf("Scope:   %s\n", pkg.Scope())
}

// END OMIT
