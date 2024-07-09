package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"time"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Ошибка при создании gzip.Reader", http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			// Замена r.Body на распакованный stream
			//зачем я это ПИСАЛ??!?!?!?!?!?!
			r.Body = io.NopCloser(gz)
		}
		next.ServeHTTP(w, r)
		w.Header().Set("Content-Encoding", "gzip")
	})
}

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
	r.ResponseWriter.WriteHeader(statusCode)
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

		h.ServeHTTP(&lw, r)           // обслуживание оригинального запроса
		duration := time.Since(start) // Since возвращает разницу во времени между start

		logger.Log.Sugar().Infoln(
			"uri:", r.RequestURI,
			"method:", r.Method,
			"status:", responseData.status, // получаем перехваченный код статуса ответа
			"duration:", duration,
			"size:", len(responseData.uri), // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn) // возвращаем функционально расширенный хендлер
}
