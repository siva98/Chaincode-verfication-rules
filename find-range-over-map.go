package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

var fset = token.NewFileSet()

func main() {

	node, err := parser.ParseFile(fset, "functions.go", nil, 0)
	if err != nil {
		log.Fatal(err) // parse error
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{Types: make(map[ast.Expr]types.TypeAndValue),
		Defs: make(map[*ast.Ident]types.Object)}
	if _, err := conf.Check("cmd/hello", fset, []*ast.File{node}, info); err != nil {
		log.Fatal(err) // type error
	}
	rangeOverMapName := ""
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.RangeStmt:
			rangeOverMapName = types.ExprString(n.X)
			ast.Inspect(node, func(x ast.Node) bool {
				if expr, ok := x.(ast.Expr); ok {
					if tv, ok := info.Types[expr]; ok {
						mapString := tv.Type.String()[0:3]
						if rangeOverMapName == nodeString(expr) && n.X.Pos() == expr.Pos() && mapString == "map" {
							fmt.Println("RANGE OVER MAP DETECTED")
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
