package middleware

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	log *logrus.Entry
	servercfg.Server
}

func NewMiddleware(log *logrus.Entry, cfg servercfg.Server) *Middleware {
	return &Middleware{
		log:    log,
		Server: cfg,
	}
}
