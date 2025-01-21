package main

import (
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcServer/delivery/router"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcServer/storages"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcServer/storages/database"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcServer/storages/fileMetric"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcServer/storages/memoryMetric"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
	"time"
)

func main() {
	var err error
	var storage storages.MetricsStorage

	//Инициализация логера
	log := logger.NewLogger()

	//Иницаилизирует все конфиги и всё в этом духе
	cfgServer := config.NewConfigServer(log)
	cfgServer.InitializeServerConfig()

	//Подключение к БД

	for i := 0; i < 3; i++ {
		storage, err = postgresMetric.NewConnectDB(log, cfgServer.Server)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		storage, err = fileMetric.NewFileStorage(cfgServer.Server)
		if err != nil {
			log.Info(fmt.Errorf("NewFileStorage err: %s", err))

		}
	}

	//Если нет ошибки подключения выгружаем метрики
	err = storage.LoadMetric()
	if err != nil {
		log.Info(err)
		storage = memoryMetric.NewMetricStorage()
	}

	//Создаёт горутину, для сохранения данных в файл
	//go func() {
	//	saveMetric.SaveMetric(cfgServer.FlagFileStoragePath, cfgServer.FlagStoreInterval)
	//}()

	//Создаёт роутер
	r := router.NewRouter(log, cfgServer.Server, storage)
	r.SetupRouter()

	//Старт сервера
	log.Info("Открыт сервер ", cfgServer.FlagAddressAndPort)
	err = http.ListenAndServe(cfgServer.FlagAddressAndPort, r.Mux /*TODO добавь передачу интерфейса для работы с сохранением*/)
	if err != nil {
		panic(err)
	}
}
