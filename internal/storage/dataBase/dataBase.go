package dataBase

import (
	"database/sql"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
)

func CrereateDB(db *sql.DB, createTableSQL string) {
	_, err := db.Exec(createTableSQL)
	if err != nil {
		logger.Log.Infof("Ошибка при создании таблицы: %v", err)

		return
	}
	logger.Log.Info("Таблица 'Metric' успешно создана или уже была.")
}

func SaveInDatabase() {
	allMetrics := storage.MetricStorage{
		Gauge:   storage.StorageMetric.GetAllGauges(),
		Counter: storage.StorageMetric.GetAllCounters(),
	}
	db, err := ConnectDB()
	if err != nil {
		logger.Log.Infoln("Ошибка связанная с ДБ: %v", err)
	}
	defer db.Close()
	metricGauges := allMetrics.GetAllGauges()
	for nameMetric, valueMetric := range metricGauges {
		querty := `INSERT INTO metric (id, type, value)
				   VALUES ($1, $2, $3)
				   ON CONFLICT (id, type) DO UPDATE
        		   SET value = EXCLUDED.value`
		_, err = db.Exec(querty, nameMetric, "gauge", valueMetric)
		if err != nil {
			logger.Log.Infoln("Ошибка при вставке данных: %v", err)
		}
	}

	metricsCounters := allMetrics.GetAllCounters()
	for nameMetric, valueMetric := range metricsCounters {
		querty := `INSERT INTO metric (id, type, value)
				   VALUES ($1, $2, $3)
				   ON CONFLICT (id, type) DO UPDATE
        		   SET value = EXCLUDED.value`
		_, err = db.Exec(querty, nameMetric, "counter", valueMetric)
		if err != nil {
			logger.Log.Infoln("Ошибка при вставке данных: %v", err)
		}
	}
}