package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

var fset = token.NewFileSet()

func main() {

	node, err := parser.ParseFile(fset, "phantom-read.go", nil, 0)

	if err != nil {
		log.Fatal(err)
	}

	writeQueryOrKey := ""
	getQueryOrKey := ""

	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncDecl:
			ast.Inspect(n.Body, func(x ast.Node) bool {
				switch x := x.(type) {
				case *ast.CallExpr:
					callExpr := nodeString(x.Fun)
					if strings.Contains(callExpr, ".") {
						putState := strings.Split(callExpr, ".")
						if putState[1] == "GetHistoryForKey" || putState[1] == "GetQueryResult" {
							getQueryOrKey = nodeString(x.Args[0])
							ast.Inspect(n.Body, func(y ast.Node) bool {
								switch y := y.(type) {
								case *ast.CallExpr:
									callExpr = nodeString(y.Fun)
									if strings.Contains(callExpr, ".") {
										putState := strings.Split(callExpr, ".")
										if putState[1] == "PutState" {
											writeQueryOrKey = nodeString(y.Args[0])
											if y.Pos() > x.Pos() && writeQueryOrKey == getQueryOrKey {
												fmt.Println("PHANTOM READ DETECTED IN FUNCTION: ")
												fmt.Println(n.Name)
											}
										}
									}
								}
								return true
							})
						}
					}
				}
				return true
			})

		}
		return true
	})

}

func nodeString(n ast.Node) string {
	var buf bytes.Buffer
	format.Node(&buf, fset, n)
	return buf.String()
}
