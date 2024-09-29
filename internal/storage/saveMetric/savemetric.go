package saveMetric

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/dataBase"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"time"
)

func SaveMetric(cdFile string, storeInterval int64, flagDatabaseDsn string) {
	if storeInterval == 0 {
		return
	}
	ticker1 := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker1.Stop()

	switch flagDatabaseDsn {
	case "":
	default:
		_, err := dataBase.ConnectDB()
		if err == nil {
			for range ticker1.C {
				dataBase.SaveInDatabase()
			}
			return
		}
	}

	switch cdFile {
	case "":
	default:
		for range ticker1.C {
			metricStorage.SaveMetricToFile(cdFile)
		}
		return
	}
}
