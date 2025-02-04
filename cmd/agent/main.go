// Package main запускает агент, который собирает метрики и отправляет их на сервер.
// Вся логика сбора и отправки метрик сосредоточена в пакете funcagent.
package main

import "github.com/LI-SeNyA-vE/KursMetrics/internal/funcagent"

// main является точкой входа для бинарного файла агента.
// Она вызывает функцию Run из пакета funcagent, где реализована основная логика работы агента.
func main() {
	funcagent.Run()
}
