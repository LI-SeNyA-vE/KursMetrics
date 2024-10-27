package storage

import (
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/errorRetriable"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"os"
)

func LoadMetricFromFile(fstg string) {
	var res []byte

	results, err := errorRetriable.ErrorRetriable(os.ReadFile, fstg)
	if err != nil {
		logger.Log.Infof("Ошибка вызова функции для повторного вызова функции: %s", err)
	}
	for _, result := range results {
		switch v := result.(type) {
		case []byte:
			res = v
		case error:
			err = v
		}
	}

	if err != nil {
		logger.Log.Infof("Ошибка чтения файла %s: %s", fstg, err)
	}

	if err := json.Unmarshal(res, &StorageMetric); err != nil {
		logger.Log.Infof("Ошибка Unmarshal: %s", err)
	}
}
