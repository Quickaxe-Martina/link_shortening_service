/*
Package main
Build: go build -o staticlint ./cmd/staticlint
Usage: ./staticlint ./...
*/
package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// noOsExitAnalyzer запрещает прямой вызов os.Exit в функции main пакета main.
var noOsExitAnalyzer = &analysis.Analyzer{
	Name: "noOsExit",
	Doc:  "Запрещает использование os.Exit в функции main пакета main.",
	Run:  runNoOsExit,
}

func runNoOsExit(pass *analysis.Pass) (any, error) {
	if strings.HasSuffix(pass.Pkg.Path(), ".test") {
        return nil, nil
    }

    if pass.Pkg.Name() != "main" {
        return nil, nil
    }
	for _, file := range pass.Files {

		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" {
				return true
			}

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				pkg, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if pkg.Name == "os" && sel.Sel.Name == "Exit" {
					pass.Reportf(call.Pos(),
						"запрещён прямой вызов os.Exit в функции main")
				}

				return true
			})

			return false
		})
	}

	return nil, nil
}

func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers,
		assign.Analyzer,
		atomic.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	)

	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[:2] == "SA" {
			analyzers = append(analyzers, a.Analyzer)
		}
	}

	for _, a := range simple.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}
	for _, a := range stylecheck.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}

	analyzers = append(analyzers, noOsExitAnalyzer)

	multichecker.Main(analyzers...)
}
