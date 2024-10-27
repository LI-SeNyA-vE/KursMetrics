package loadMetric

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/dataBase"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
)

func InitializeStorage() {
	var err error

	if config.ConfigServerFlags.FlagDatabaseDsn != "" {
		err = dataBase.LoadMetricFromDB()
	}
	if err != nil {
		logger.Log.Infof("Ошибка получена из функции LoadMetricFromDB: %s | флаг FlagRestore: %v", err, config.ConfigServerFlags.FlagRestore)
	}

	if err != nil && config.ConfigServerFlags.FlagRestore {
		logger.Log.Info("Зашли в функцию LoadMetricFromFile")
		metricStorage.LoadMetricFromFile(config.ConfigServerFlags.FlagFileStoragePath)
	}
	return
}
