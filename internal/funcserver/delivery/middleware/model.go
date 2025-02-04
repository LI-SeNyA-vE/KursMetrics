/*
Package middleware предоставляет функциональность промежуточных обработчиков (middleware)
для HTTP-сервера KursMetrics. Включает в себя логику логирования запросов,
проверки целостности (HMAC SHA256), gzip-сжатия и т.д.
*/
package middleware

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/sirupsen/logrus"
)

// Middleware содержит ссылку на общий логгер (logrus.Entry) и структуру конфигурации сервера.
// Используется при инициализации набора промежуточных обработчиков (LoggingMiddleware,
// HashSHA256, GzipMiddleware и др.).
type Middleware struct {
	log *logrus.Entry
	servercfg.Server
}

// NewMiddleware создаёт новый объект Middleware с заданным логгером и конфигурацией сервера.
// Полученный объект затем может быть использован для инициализации различных middleware-функций,
// которые внедряются в цепочку обработки HTTP-запросов.
func NewMiddleware(log *logrus.Entry, cfg servercfg.Server) *Middleware {
	return &Middleware{
		log:    log,
		Server: cfg,
	}
}
