package loadMetric

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/dataBase"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
)

func InitializeStorage() {
	logger.Log.Infof("Флаг БД %s Флаг файла %t", config.ConfigFlags.FlagDatabaseDsn, config.ConfigFlags.FlagRestore)
	var err error

	if config.ConfigFlags.FlagDatabaseDsn != "" {
		err = dataBase.LoadMetricFromDB()
	}

	if err != nil && config.ConfigFlags.FlagRestore {
		metricStorage.LoadMetricFromFile(config.ConfigFlags.FlagFileStoragePath)
	}
	return
}
