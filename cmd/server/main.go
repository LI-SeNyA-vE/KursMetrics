package main

import (
	config "github.com/LI-SeNyA-vE/KursMetrics/internal/config"
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
	cfgFlags := config.InitializeConfig()

	//Запускается функция, которая определит куда сохранять данные
	dataBase.InitializeStorage()

	//Создаёт горутину, для сохранения данных в файл
	go func() { saveMetric.SaveMetric(cfgFlags.FlagFileStoragePath, cfgFlags.FlagStoreInterval) }()

	//Создаёт роутер
	r := handlers.SetapRouter()

	//Старт сервера
	startServer(r, cfgFlags)
}

func startServer(r *chi.Mux, cfgFlags config.VarFlag) {
	logger.Log.Info("Открыт сервер ", cfgFlags.FlagAddressAndPort)
	err := http.ListenAndServe(cfgFlags.FlagAddressAndPort, r)
	if err != nil {
		panic(err)
	}

}
