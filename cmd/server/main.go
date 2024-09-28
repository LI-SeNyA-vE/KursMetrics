package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/dataBase"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/saveMetric"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func main() {
	//Иницаилизирует все конфиги и всё в этом духе
	config.InitializeConfig()

	//Запускается функция, которая определит откуда выгружать данные
	dataBase.InitializeStorage()

	//Создаёт горутину, для сохранения данных в файл
	go func() {
		saveMetric.SaveMetric(config.ConfigFlags.FlagFileStoragePath, config.ConfigFlags.FlagStoreInterval, config.ConfigFlags.FlagDatabaseDsn)
	}()

	//Создаёт роутер
	r := handlers.SetapRouter()

	//Старт сервера
	startServer(r)
}

func startServer(r *chi.Mux) {
	logger.Log.Info("Открыт сервер ", config.ConfigFlags.FlagAddressAndPort)
	err := http.ListenAndServe(config.ConfigFlags.FlagAddressAndPort, r)
	if err != nil {
		panic(err)
	}

}
