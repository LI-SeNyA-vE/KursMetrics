// Package main запускает агент, который собирает метрики и отправляет их на сервер.
// Вся логика сбора и отправки метрик сосредоточена в пакете agent.
package main

import (
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/agent"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main является точкой входа для бинарного файла агента.
// Она вызывает функцию Run из пакета agent, где реализована основная логика работы агента.
// Пример вызова go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" main.go
func main() {
	version()
	agent.Run()
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
