package handlers

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Router struct {
	log *logrus.Entry
	config.Server
	*chi.Mux
}

func NewRouter(log *logrus.Entry, cfg config.Server) *Router {
	return &Router{
		log:    log,
		Server: cfg,
		Mux:    nil,
	}
}

func (rout *Router) SetupRouter() {
	rout.Mux = chi.NewRouter()
	mw := middleware.NewMiddleware(rout.log, rout.Server)
	hl := NewHandler(rout.log, rout.Server)

	rout.Mux.Use(mw.HashSHA256)
	rout.Mux.Use(mw.LoggingMiddleware)
	rout.Mux.Use(mw.GzipMiddleware)
	rout.Mux.Use(mw.UnGzipMiddleware)

	rout.Mux.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", hl.PostAddValue) //Обновление по URL

	rout.Mux.Post("/value/", hl.JSONValue)             //Получение через JSON
	rout.Mux.Post("/update/", hl.JSONUpdate)           //Обновление через JSON
	rout.Mux.Post("/updates/", hl.PostAddArrayMetrics) //Обновление "батчем" через JSON

	rout.Mux.Get("/value/{typeMetric}/{nameMetric}", hl.GetReceivingMetric) //Получение по URL
	rout.Mux.Get("/ping", hl.GetReceivingAllMetric)                         //Проверка подключения к БД
	rout.Mux.Get("/", hl.GetReceivingAllMetric)                             //Получение по JSON
	rout.Mux.Get("/ping", hl.Ping)
}
