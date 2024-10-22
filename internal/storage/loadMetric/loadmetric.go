package loadMetric

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/dataBase"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
)

func InitializeStorage() {
	var err error

	if config.ConfigFlags.FlagDatabaseDsn != "" {
		err = dataBase.LoadMetricFromDB()
	}

	if err != nil && config.ConfigFlags.FlagRestore {
		metricStorage.LoadMetricFromFile(config.ConfigFlags.FlagFileStoragePath)
	}
	return
}
