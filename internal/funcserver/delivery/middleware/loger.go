package middleware

import (
	"log"
	"net/http"
	"time"
)

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

func (m *Middleware) LoggingMiddleware(h http.Handler) http.Handler {
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

		log.Print(
			"uri:", r.RequestURI,
			" status:", responseData.status, // получаем перехваченный код статуса ответа
			" duration:", duration,
			" size:", len(responseData.uri), // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn) // возвращаем функционально расширенный хендлер
}
