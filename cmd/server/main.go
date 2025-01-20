package main

import (
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/fileMetric"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	database_v2 "github.com/LI-SeNyA-vE/KursMetrics/internal/storeage2/database"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
	"time"
)

func main() {
	var err error
	var storage metricStorage.MetricsStorage

	//Инициализация логера
	log := logger.NewLogger()

	//Иницаилизирует все конфиги и всё в этом духе
	cfgServer := config.NewConfigServer(log)
	cfgServer.InitializeServerConfig()

	//Подключение к БД

	for i := 0; i < 3; i++ {
		storage, err = database_v2.NewConnectDB(log, cfgServer.Server)
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
	}

	//Создаёт горутину, для сохранения данных в файл
	//go func() {
	//	saveMetric.SaveMetric(cfgServer.FlagFileStoragePath, cfgServer.FlagStoreInterval)
	//}()

	//Создаёт роутер
	r := handlers.NewRouter(log, cfgServer.Server, storage)
	r.SetupRouter()

	//Старт сервера
	log.Info("Открыт сервер ", cfgServer.FlagAddressAndPort)
	err = http.ListenAndServe(cfgServer.FlagAddressAndPort, r.Mux /*TODO добавь передачу интерфейса для работы с сохранением*/)
	if err != nil {
		panic(err)
	}
}
