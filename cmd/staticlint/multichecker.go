package main

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	mychecks := []*analysis.Analyzer{
		OsExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		inspect.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
	}
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	multichecker.Main(
		mychecks...,
	)
}

func getPackageName(file *ast.File, pos token.Pos) string {
	if file.Pos() <= pos && pos <= file.End() {
		return file.Name.Name
	}
	return "unknown"
}
