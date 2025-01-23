package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
)

func ensureTablesExist(db *sql.DB, log *logrus.Entry) error {
	// Проверка существования таблицы gauges
	gaugesQuery := queryExistGaugesTable
	_, err := db.Exec(gaugesQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы gauges: %w", err)
	}

	// Проверка существования таблицы counters
	countersQuery := queryExistCounterTable
	_, err = db.Exec(countersQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы counters: %w", err)
	}

	log.Println("таблицы gauges и counters проверены/созданы успешно.")
	return nil
}
