// Package main запускает сервер, который принимает, хранит и предоставляет метрики.
// Вся логика «серверной» части проекта расположена в пакете funcserver.
package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver"
	_ "github.com/jackc/pgx/v5/stdlib" // Импорт драйвера PostgreSQL через pgx
)

// main является точкой входа для запуска сервера (Server).
// Он вызывает funcserver.Run(), где реализована вся основная логика принятия и хранения метрик,
// а также инициализация и запуск HTTP-сервера.
func main() {
	funcserver.Run()
}
