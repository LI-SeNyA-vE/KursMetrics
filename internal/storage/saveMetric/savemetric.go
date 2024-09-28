package saveMetric

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/dataBase"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"time"
)

var cfgFlags = config.VarFlag{}

func SaveMetric(cdFile string, storeInterval int64) {
	if storeInterval == 0 {
		return
	}
	ticker1 := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker1.Stop()

	switch cfgFlags.FlagDatabaseDsn {
	case "":
	default:
		for range ticker1.C {
			dataBase.SaveInDatabase()
		}
		return
	}

	switch cfgFlags.FlagFileStoragePath {
	case "":

	default:
		for range ticker1.C {
			metricStorage.SaveMetricToFile(cdFile)
		}
		return
	}
}
