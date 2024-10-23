package loadMetric

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/dataBase"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
)

func InitializeStorage() {
	var err error

	if config.ConfigServerFlags.FlagDatabaseDsn != "" {
		err = dataBase.LoadMetricFromDB()
	}

	if err != nil && config.ConfigServerFlags.FlagRestore {
		metricStorage.LoadMetricFromFile(config.ConfigServerFlags.FlagFileStoragePath)
	}
	return
}
