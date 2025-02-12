// Package main запускает сервер, который принимает, хранит и предоставляет метрики.
// Вся логика «серверной» части проекта расположена в пакете funcserver.
package main

import (
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver"
	_ "github.com/jackc/pgx/v5/stdlib" // Импорт драйвера PostgreSQL через pgx
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main является точкой входа для запуска сервера (Server).
// Он вызывает funcserver.Run(), где реализована вся основная логика принятия и хранения метрик,
// а также инициализация и запуск HTTP-сервера.
// Пример вызова go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" main.go
func main() {
	version()
	funcserver.Run()
}

func version() {
	if buildVersion != "" {
		fmt.Println("Build version: " + buildVersion)
	} else {
		fmt.Println("Build version: N/A")
	}

	if buildDate != "" {
		fmt.Println("Build date: " + buildDate)
	} else {
		fmt.Println("Build version: N/A")
	}

	if buildCommit != "" {
		fmt.Println("Build commit: " + buildCommit)
	} else {
		fmt.Println("Build version: N/A")
	}
}
