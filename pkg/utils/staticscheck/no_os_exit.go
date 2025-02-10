// Package staticscheck реализует кастомный чекер для проверки вызова osExit.
package staticscheck

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// noOsExitAnalyzer – наш кастомный анализатор.
// Он проверяет, что в функции main (пакет main) нет прямого вызова os.Exit.
var noOsExitAnalyzer = &analysis.Analyzer{
	Name: "noOsExitInMain",
	Doc:  "Проверяет, что нет прямых вызовов os.Exit в функции main в пакете main",
	// Для удобства используем уже готовый анализатор inspect.
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: runNoOsExit,
}

// runNoOsExit – основная логика нашего анализатора.
// Ищем все вызовы os.Exit в функции main (pkg main).
func runNoOsExit(pass *analysis.Pass) (interface{}, error) {
	// Если пакет не main – выходим сразу
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	// Получаем инспектор из результатов зависимого анализа
	ins, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Вариант обхода всего AST
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil), // интересуют объявления функций
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		// Ищем именно функцию с именем "main"
		if fn.Name != nil && fn.Name.Name == "main" {
			// Обходим тело функции, смотрим все вызовы
			ast.Inspect(fn.Body, func(x ast.Node) bool {
				callExpr, ok := x.(*ast.CallExpr)
				if !ok {
					return true
				}

				// Проверяем, что вызов сделан из пакета os к функции Exit
				if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := sel.X.(*ast.Ident); ok {
						if ident.Name == "os" && sel.Sel.Name == "Exit" {
							// Фиксируем ошибку
							pass.Reportf(x.Pos(), "запрещён прямой вызов os.Exit в функции main пакета main")
						}
					}
				}
				return true
			})
		}
	})

	return nil, nil
}
