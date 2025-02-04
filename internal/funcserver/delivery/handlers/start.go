/*
Package handlers содержит набор конструкторов и функций-обработчиков,
которые позволяют работать с HTTP-запросами к серверу метрик.
*/
package handlers

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"github.com/sirupsen/logrus"
)

// NewHandler создаёт новый Handler, инициализируя его логгером, конфигурацией сервера
// и интерфейсом хранилища метрик. Возвращает указатель на готовую структуру Handler,
// которую затем можно использовать в роутере (для регистрации HTTP-хендлеров).
func NewHandler(log *logrus.Entry, cfg servercfg.Server, storage storages.MetricsStorage) *Handler {
	return &Handler{
		log:     log,
		cfg:     cfg,
		storage: storage,
	}
}
