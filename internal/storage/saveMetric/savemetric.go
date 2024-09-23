package saveMetric

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"time"
)

func SaveMetric(cdFile string, storeInterval int64) {
	if storeInterval == 0 {
		return
	}
	ticker1 := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker1.Stop()

	switch *config.FlagDatabaseDsn {
	case "":
	default:
		for range ticker1.C {
			metricStorage.SaveInDatabase()
		}
		return
	}

	switch *config.FlagFileStoragePath {
	case "":

	default:
		for range ticker1.C {
			metricStorage.SaveMetricToFile(cdFile)
		}
		return
	}
}
