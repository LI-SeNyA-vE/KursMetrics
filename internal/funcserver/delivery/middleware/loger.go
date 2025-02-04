/*
Package middleware предоставляет набор промежуточных обработчиков (middleware),
используемых в цепочке обработки HTTP-запросов. LoggingMiddleware отвечает
за ведение логов: оно замеряет время выполнения запроса и, после его
обработки, выводит информацию о статусе ответа, URI и длительности обработки.
*/
package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseData хранит код статуса (HTTP) и строку запроса (URI),
// которые понадобятся для логирования.
type responseData struct {
	status int
	uri    string
}

// loggingResponseWriter оборачивает оригинальный http.ResponseWriter,
// перехватывая вызов WriteHeader для сохранения кода статуса ответа.
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// WriteHeader перехватывает вызов и записывает код статуса в responseData,
// помимо вызова метода WriteHeader оригинального ResponseWriter.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LoggingMiddleware оборачивает переданный http.Handler, замеряя длительность
// выполнения запроса. После обработки выводит в лог информацию о
// URI запроса, коде статуса, времени обработки и длине URI.
//
// Пример использования:
//
//	r := chi.NewRouter()
//	r.Use(middleware.NewMiddleware(log, cfg).LoggingMiddleware)
//	r.Get("/", handler)
func (m *Middleware) LoggingMiddleware(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // фиксируем текущее время для вычисления длительности
		responseData := &responseData{
			status: 0,
			uri:    r.RequestURI,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)           // передаём управление оригинальному хендлеру
		duration := time.Since(start) // вычисляем время выполнения

		log.Print(
			"uri:", r.RequestURI,
			" status:", responseData.status,
			" duration:", duration,
			" size:", len(responseData.uri),
		)
	}
	return http.HandlerFunc(logFn)
}
