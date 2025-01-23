package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	funcserver.Run()
}
