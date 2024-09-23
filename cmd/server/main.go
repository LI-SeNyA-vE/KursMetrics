package main

import (
	config "github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers/middleware/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/saveMetric"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func main() {
	//Иницаилизирует все конфиги и всё в этом духе
	config.InitializeConfigServer()

	//Создаёт горутину, для сохранения данных в файл
	go func() { saveMetric.SaveMetric(*config.FlagFileStoragePath, *config.FlagStoreInterval) }()

	//Создаёт роутер
	r := handlers.SetapRouter()

	//Старт сервера
	startServer(r)
}

func startServer(r *chi.Mux) {
	err := http.ListenAndServe(*config.FlagAddressAndPort, r)
	if err != nil {
		panic(err)
	}
	logger.Log.Info("Открыт сервер ", *config.FlagAddressAndPort)
}
