/*
Package main
Build: go build -o staticlint ./cmd/staticlint
Usage: ./staticlint ./...
*/
package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// noOsExitAnalyzer запрещает прямой вызов os.Exit в функции main пакета main.
var noOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc: `Запрещает использование os.Exit в функции main пакета main.

Использование os.Exit в main нарушает корректное завершение программы:
- не выполняются deferred-вызовы
- усложняется тестирование
- ухудшается читаемость

Рекомендуется возвращать код ошибки из main через логирование или обработку ошибок.`,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: runNoOsExit,
}

func runNoOsExit(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		// проверяем пакет main
		if pass.Pkg.Name() != "main" {
			continue
		}

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

	// Стандартные анализаторы golang.org/x/tools
	analyzers = append(analyzers,
		assign.Analyzer,
		atomic.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	)

	// Все анализаторы класса SA (staticcheck)
	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[:2] == "SA" {
			analyzers = append(analyzers, a.Analyzer)
		}
	}

	// Анализаторы других классов staticcheck
	for _, a := range simple.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}
	for _, a := range stylecheck.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}

	// Собственный анализатор
	analyzers = append(analyzers, noOsExitAnalyzer)

	multichecker.Main(analyzers...)
}
