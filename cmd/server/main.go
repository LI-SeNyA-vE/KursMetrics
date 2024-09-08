package main

import (
	config "github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers/middleware"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers/middleware/logger"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"time"
)

func main() {
	config.InitializeGlobals()
	err := logger.Initialize("info")
	if err != nil {
		panic("Не удалось инициализировать логгер")
	}
	initializeStorage(*config.FlagFileStoragePath, *config.FlagRestore, *config.FlagDatabaseDsn)
	go func() { startTicker(*config.FlagFileStoragePath, *config.FlagStoreInterval) }()

	r := setapRouter()

	logger.Log.Info(logger.Log.Level(), "Открыт сервер ", *config.FlagAddressAndPort)
	startServer(r)
}

func initializeStorage(cdFile string, resMetricBool bool, loadDataBase string) {
	if loadDataBase != "" {
		db, err := config.ConnectDB()
		if err != nil {
			logger.Log.Infoln("Ошибка связанная с ДБ: %v", err)
		}
		defer db.Close()

		configCreateSQL := config.ConfigSQL()
		metricStorage.CrereateDB(db, configCreateSQL)

		rows, err := db.Query("SELECT Id, Type, Name, Value FROM your_table_name")
		if err != nil {
			logger.Log.Infoln("Ошибка получения данных из базы данных: %v", err)
		} else {
			for rows.Next() {
				metric := &metricStorage.MetricStorage{}
				var idMetric string
				var typeMetric string
				var nameMetric string
				var valueMetric float64
				//err := rows.Scan(&metric.Id, &metric.Type, &metric.Name, &metric.Value)
				err := rows.Scan(idMetric, typeMetric, nameMetric, valueMetric)
				if err != nil {
					logger.Log.Infoln("Ошибка сканирования строки: %v", err)
				}
				defer rows.Close()

				switch typeMetric { //Свитч для проверки что это запрос или gauge или counter
				case "gauge": //Если передано значение 'gauge'
					metric.UpdateGauge(nameMetric, valueMetric)
				case "counter": //Если передано значение 'counter'
					metric.UpdateCounter(nameMetric, int64(valueMetric))
				default: //Если передано другое значение значение
					log.Println("При вытягивание данных из БД оказалось что тип не gauge и не counter")
				}
			}
			return
		}

		//// Проверка на ошибки, которые могли произойти при итерировании по строкам
		//if err = rows.Err(); err != nil {
		//	log.Fatalf("Ошибка при итерировании по строкам: %v", err)
		//}
	}
	if resMetricBool {
		metricStorage.LoadMetricFromFile(cdFile)
	}
	return
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
