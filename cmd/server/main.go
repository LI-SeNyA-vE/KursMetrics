package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
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
	r := chi.NewRouter()

	r.Use(func(h http.Handler) http.Handler {
		return logger.LoggingMiddleware(h)
	})

	// Разобрать что я натворил в коде до и в логере

	r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", handlers.PostAddValue)

	r.Get("/value/{typeMetric}/{nameMetric}", handlers.GetReceivingMetric)
	r.Get("/", handlers.GetReceivingAllMetric)

	sugar.Log(logger.Log.Level(), "Открыт сервер ", *config.AddressAndPort)
	err := http.ListenAndServe(*config.AddressAndPort, r)
	if err != nil {
		panic(err)
	}
}

func JSONhandler(w http.ResponseWriter, req *http.Request) {
	var id string

	if req.Method == http.MethodPost {
		var metrics Metrics
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id = strconv.Itoa(metrics.ID)
		metrics[id] = metrics
	}

	id = req.URL.Query().Get("id")
	resp, err := json.Marshal(visitors[id])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
