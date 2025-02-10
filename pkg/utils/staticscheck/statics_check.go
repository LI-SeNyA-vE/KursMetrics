// Package staticscheck реализует мультичекер.
package staticscheck

import (
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/analysis/facts/directives"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/stylecheck"

	// Модуль multichecker из стандартного пакета инструментов
	"golang.org/x/tools/go/analysis/multichecker"

	// Анализаторы от staticcheck
	//
	// Список доступных анализаторов см. в репозитории https://github.com/dominikh/go-tools
	// или https://pkg.go.dev/honnef.co/go/tools/staticcheck (для SA, ST, S и т.д.)
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"

	// Пример публичных анализаторов (можно выбрать другие):
	// 1. github.com/gostaticanalysis/nilerr
	// 2. github.com/gostaticanalysis/sqlrows
	"github.com/gostaticanalysis/nilerr"
	"github.com/gostaticanalysis/sqlrows/passes/sqlrows"
	// Импорт кастомного анализатора, который мы сами реализуем
	// "noosexit.go" располагается в том же пакете main.
)

// MyMulticheck - точка входа для нашего multichecker.
func MyMulticheck() {
	var analyzers []*analysis.Analyzer

	var myAnalyzers = []*analysis.Analyzer{
		// Стандартные анализаторы:
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		directives.Analyzer,
		//embedcfg.Analyzer,
		errorsas.Analyzer,
		//fieldalignment.Analyzer,
		findcall.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}

	// Мы собираем все нужные анализаторы
	analyzers = append(analyzers,
		myAnalyzers...)

	for _, v := range staticcheck.Analyzers {
		// добавляем в массив все проверки
		analyzers = append(analyzers,
			v.Analyzer)
	}

	for _, v := range stylecheck.Analyzers {
		// добавляем в массив все проверки
		analyzers = append(analyzers,
			v.Analyzer)
	}

	for _, v := range simple.Analyzers {
		// добавляем в массив все проверки
		analyzers = append(analyzers,
			v.Analyzer)
	}

	// Добавляем публичные анализаторы (пример с nilerr и sqlrows)
	analyzers = append(analyzers,
		nilerr.Analyzer,
		sqlrows.Analyzer,
	)

	// Добавляем наш кастомный анализатор (см. файл noosexit.go)
	analyzers = append(analyzers,
		noOsExitAnalyzer, // Важно, чтобы имя совпадало с тем, как мы экспортируем его
	)

	// Запускаем multichecker
	multichecker.Main(analyzers...)
}
