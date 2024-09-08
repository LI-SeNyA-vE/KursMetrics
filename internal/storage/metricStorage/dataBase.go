package storage

import (
	"database/sql"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers/middleware/logger"
)

func CrereateDB(db *sql.DB, createTableSQL string) {
	_, err := db.Exec(createTableSQL)
	if err != nil {
		logger.Log.Infoln("Ошибка при создании таблицы: %v", err)

		return
	}
	logger.Log.Info("Таблица 'Metric' успешно создана или уже была.")
}
