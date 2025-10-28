package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer - search direct usage of os.Exit
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osExitUsageFinder",
	Doc:  "search direct usage of os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if file.Name.Name != "main" {
				return true
			}
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" && getPackageName(file, x.Pos()) == "main" {
					ast.Inspect(x, func(node ast.Node) bool {
						switch xx := node.(type) {
						case *ast.CallExpr:
							selExpr, ok := xx.Fun.(*ast.SelectorExpr)
							if !ok {
								return true
							}
							ident, ok := selExpr.X.(*ast.Ident)
							if !ok {
								return true
							}
							if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
								pass.Reportf(x.Pos(), "usage of os.Exit found on main.main")
								return false
							}
						}
						return true
					})
				}
			}
			return true
		})
	}
	return nil, nil
}
