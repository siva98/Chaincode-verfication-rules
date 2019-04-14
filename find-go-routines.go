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

	node, err := parser.ParseFile(fset, "goroutines.go", nil, 0)

	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.GoStmt:
			fmt.Println("GOROUTINE DETECTED AT POSITION: ")
			fmt.Print(n.Go)
		}
		return true
	})

}
