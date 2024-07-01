package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"go.uber.org/zap"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
	"github.com/go-chi/chi/v5"
)

var sugar zap.SugaredLogger

func main() {
	log, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}
	defer log.Sync()
	sugar = *log.Sugar()

	if err := logger.Initialize("debug"); err != nil {
		panic(err)
	}

	defer logger.Log.Sync()

	cfg := config.GetConfig()
	config.InitializeGlobals(cfg)
	r := chi.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return LoggingMiddleware(h)
	})

	// Разобрать что я натворил в коде до и в логере

	r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", handlers.PostAddValue)

	r.Get("/value/{typeMetric}/{nameMetric}", handlers.GetReceivingMetric)
	r.Get("/", handlers.GetReceivingAllMetric)

	fmt.Println("Открыт сервер ", *config.AddressAndPort)
	err = http.ListenAndServe(*config.AddressAndPort, r)
	if err != nil {
		panic(err)
	}
}

// /
// /
// /
// /
// /
// /
// /
// /
// /
// /
// /
// /

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		uri    string
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode
}

func LoggingMiddleware(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // функция Now() возвращает текущее время
		responseData := &responseData{
			status: 0,
			uri:    r.RequestURI, // Захват URI из запроса
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)           // точка, где выполняется хендлер pingHandler // обслуживание оригинального запроса
		duration := time.Since(start) // Since возвращает разницу во времени между start

		sugar.Infoln(
			"uri:", r.RequestURI,
			"method:", r.Method,
			"status:", responseData.status, // получаем перехваченный код статуса ответа
			"duration:", duration,
			"size:", responseData.uri, // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn) // возвращаем функционально расширенный хендлер
}
