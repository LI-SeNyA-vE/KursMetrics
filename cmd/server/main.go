package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/dataBase"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/saveMetric"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func main() {
	var err error
	//Инициализация логера
	log := logger.NewLogger()

	//Иницаилизирует все конфиги и всё в этом духе
	cfgServer := config.NewConfigServer(log)
	cfgServer.InitializeServerConfig()

	//Подключение к БД
	db := dataBase.NewConnectDB(log, cfgServer.Server)
	for i := 0; i < 3; i++ {
		err = db.ConnectDB()
		if err == nil {
			break
		}
	}

	//Если ошибка подключения к БД, берём значение из файла
	if err != nil {
		err = metricStorage.LoadMetricFromFile(cfgServer.FlagFileStoragePath)
		log.Info(err)
	} else {
		err = db.LoadMetricFromDB()
		if err != nil {
			return
		}
	}

	//Создаёт горутину, для сохранения данных в файл
	go func() {
		saveMetric.SaveMetric(cfgServer.FlagFileStoragePath, cfgServer.FlagStoreInterval)
	}()

	//Создаёт роутер
	r := handlers.NewRouter(log, cfgServer.Server)
	r.SetupRouter()

	//Старт сервера
	log.Info("Открыт сервер ", cfgServer.FlagAddressAndPort)
	err = http.ListenAndServe(cfgServer.FlagAddressAndPort, r.Mux)
	if err != nil {
		panic(err)
	}
}
