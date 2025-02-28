/*
Package handlers содержит структуру Handler и связанные с ней функции-обработчики HTTP-запросов.
Handler инкапсулирует в себе логгер, конфигурацию сервера и хранилище метрик, что позволяет
организовать доступ к ним внутри хендлеров.
*/
package handlers

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/storages"
	"github.com/sirupsen/logrus"
)

// Handler хранит ссылки на логгер (logrus.Entry), конфигурацию сервера (servercfg.Server)
// и реализацию интерфейса MetricsStorage (storage). Используется в хендлерах для взаимодействия
// с базой метрик, а также для логирования запросов/ответов.
type Handler struct {
	log     *logrus.Entry
	cfg     servercfg.Server
	storage storages.MetricsStorage
}
