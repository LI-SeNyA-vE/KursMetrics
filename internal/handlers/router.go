package handlers

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SetapRouter() *chi.Mux {
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

	r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", PostAddValue) //Обновление по URL

	r.Post("/value/", JSONValue)             //Получение через JSON
	r.Post("/update/", JSONUpdate)           //Обновление через JSON
	r.Post("/updates/", PostAddArrayMetrics) //Обновление "батчем" через JSON

	r.Get("/value/{typeMetric}/{nameMetric}", GetReceivingMetric) //Получение по URL
	r.Get("/ping", GetReceivingAllMetric)                         //Проверка подключения к БД
	r.Get("/", GetReceivingAllMetric)                             //Получение по JSON
	r.Get("/ping", Ping)
	return r
}
