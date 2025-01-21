package main

import (
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/delivery/router"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages/database"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages/filemetric"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages/memorymetric"
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
		storage, err = database.NewConnectDB(log, cfgServer.Server)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		storage, err = filemetric.NewFileStorage(cfgServer.Server)
		if err != nil {
			log.Info(fmt.Errorf("ошибка при объявление хранения в файле err: %s", err))
			storage = memorymetric.NewMetricStorage()
		}
	}

	//Если нет ошибки подключения выгружаем метрики
	err = storage.LoadMetric()
	if err != nil {
		log.Info(err)
		storage = memorymetric.NewMetricStorage()
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
