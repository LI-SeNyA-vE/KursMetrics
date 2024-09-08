package storage

import (
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers/middleware/logger"
	"os"
)

func LoadMetricFromFile(fstg string) {

	res, err := os.ReadFile(fstg)
	if err != nil {
		logger.Log.Info("Ошибка чтения файла: %s", err)
	}

	if err := json.Unmarshal(res, &StorageMetric); err != nil {
		logger.Log.Info("Ошибка Unmarshal: %s", err)
	}
}
