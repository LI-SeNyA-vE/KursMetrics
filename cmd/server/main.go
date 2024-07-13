package main

import (
	"net/http"
	"time"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers/middleware"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/go-chi/chi/v5"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func main() {
	if err := logger.Initialize("debug"); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	sugar := *logger.Log.Sugar()

	cfg := config.GetConfig()
	config.InitializeGlobals(cfg)

	initializeStorage(*config.FlagFileStoragePath, *config.FlagRestore)

	go func() { startTicker(*config.FlagFileStoragePath, *config.FlagStoreInterval) }()

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

	r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", handlers.PostAddValue)

	r.Post("/value/", handlers.JSONValue)
	r.Post("/update/", handlers.JSONUpdate)

	r.Get("/value/{typeMetric}/{nameMetric}", handlers.GetReceivingMetric)
	r.Get("/", handlers.GetReceivingAllMetric)

	sugar.Log(logger.Log.Level(), "Открыт сервер ", *config.FlagAddressAndPort)
	err := http.ListenAndServe(*config.FlagAddressAndPort, r)
	if err != nil {
		panic(err)
	}
}

func initializeStorage(cdFile string, resMetricBool bool) {
	if resMetricBool {
		metricStorage.LoadMetricFromFile(cdFile)
	}
}

func startTicker(cdFile string, storeInterval int64) {
	ticker1 := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker1.Stop()

	for range ticker1.C {
		metricStorage.SaveMetricToFile(cdFile)
	}
}
