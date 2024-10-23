package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/loadMetric"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/saveMetric"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func main() {
	//Иницаилизирует все конфиги и всё в этом духе
	config.InitializeServerConfig()

	//Запускается функция, которая определит откуда выгружать данные
	loadMetric.InitializeStorage()

	//Создаёт горутину, для сохранения данных в файл
	go func() {
		saveMetric.SaveMetric(config.ConfigServerFlags.FlagFileStoragePath, config.ConfigServerFlags.FlagStoreInterval, config.ConfigServerFlags.FlagDatabaseDsn)
	}()

	//Создаёт роутер
	r := handlers.SetapRouter()

	//Старт сервера
	startServer(r)
}

func startServer(r *chi.Mux) {
	logger.Log.Info("Открыт сервер ", config.ConfigServerFlags.FlagAddressAndPort)
	err := http.ListenAndServe(config.ConfigServerFlags.FlagAddressAndPort, r)
	if err != nil {
		panic(err)
	}

}
