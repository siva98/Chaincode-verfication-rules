package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, "find-global-vars.go", nil, 0)

	if err != nil {
		log.Fatal(err)
	}

	blacklistedImports := []string{"time"}

	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ImportSpec:
			for _, blacklistedImport := range blacklistedImports {
				blacklistedImport = "\"" + blacklistedImport + "\""

				if n.Path.Value == blacklistedImport {
					fmt.Println("BLACKLISTED IMPORT FOUND: ")
					fmt.Print(blacklistedImport)
				}
			}

		}
		return true
	})

}
