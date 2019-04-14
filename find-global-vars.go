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

	node, err := parser.ParseFile(fset, "var-scope.go", nil, 0)

	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.GenDecl:
			if n.Tok == token.VAR {
				ast.Inspect(node, func(x ast.Node) bool {
					switch x := x.(type) {
					case *ast.FuncDecl:
						if n.TokPos < x.Body.Lbrace || n.TokPos > x.Body.Rbrace {
							fmt.Println("GLOBAL DECLARATION DETECTED")
						}
					}
					return true
				})

			}
		}
		return true
	})

}
