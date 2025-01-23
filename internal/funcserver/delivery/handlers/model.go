package handlers

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	log     *logrus.Entry
	cfg     servercfg.Server
	storage storages.MetricsStorage
}
