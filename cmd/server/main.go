package main

import (
	config "github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers/middleware"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
	"time"
)

func main() {
	config.InitializeGlobals()
	if err := logger.Initialize("debug"); err != nil {
		panic(err)
	}

	sugar := *logger.Log.Sugar()
	//Причесать логер

	if *config.FlagDatabaseDsn != "" {
		db, err := config.ConnectDB()
		if err != nil {
			sugar.Log(logger.Log.Level(), "Ошибка связанная с ДБ ", err)
		}
		_, configCreateSQL := config.ConfigSQL()
		metricStorage.CrereateDB(db, configCreateSQL)
	}

	initializeStorage(*config.FlagFileStoragePath, *config.FlagRestore, *config.FlagDatabaseDsn)
	go func() { startTicker(*config.FlagFileStoragePath, *config.FlagStoreInterval) }()

	r := setapRouter()

	sugar.Log(logger.Log.Level(), "Открыт сервер ", *config.FlagAddressAndPort)
	startServer(r)
}

func initializeStorage(cdFile string, resMetricBool bool, loadDataBase string) {
	if loadDataBase != "" {

	}
	if resMetricBool {
		metricStorage.LoadMetricFromFile(cdFile)
	}

}

func startTicker(cdFile string, storeInterval int64) {
	if storeInterval == 0 {
		return
	}
	ticker1 := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker1.Stop()

	for range ticker1.C {
		metricStorage.SaveMetricToFile(cdFile)
	}
}

func setapRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return middleware.LoggingMiddleware(h)
	})
	r.Use(func(h http.Handler) http.Handler {
		return middleware.GzipMiddleware(h)
	})
	r.Use(func(h http.Handler) http.Handler {
		return middleware.UnGzipMiddleware(h)
	})

	r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", handlers.PostAddValue) //Обновление по URL

	r.Post("/value/", handlers.JSONValue)   //Обновлени через JSON
	r.Post("/update/", handlers.JSONUpdate) //Обновлени через JSON

	r.Get("/value/{typeMetric}/{nameMetric}", handlers.GetReceivingMetric) //Получение по URL
	r.Get("/ping", handlers.GetReceivingAllMetric)                         //Проверка подключения к БД
	r.Get("/", handlers.GetReceivingAllMetric)                             //Получение по JSON
	r.Get("/ping", handlers.Ping)
	return r
}

func startServer(r *chi.Mux) {
	err := http.ListenAndServe(*config.FlagAddressAndPort, r)
	if err != nil {
		panic(err)
	}
}
