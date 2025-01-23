package handlers

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"github.com/sirupsen/logrus"
)

func NewHandler(log *logrus.Entry, cfg servercfg.Server, storage storages.MetricsStorage) *Handler {
	return &Handler{
		log:     log,
		cfg:     cfg,
		storage: storage,
	}
}
