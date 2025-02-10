// Package postgresql предоставляет реализацию хранилища метрик в базе данных PostgreSQL.
// Функция ensureTablesExist выполняет проверку наличия (и при необходимости
// создание) таблиц gauges и counters, необходимых для хранения метрик.
package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

// ensureTablesExist проверяет, существуют ли в базе таблицы gauges и counters.
// Если их нет — создаёт, используя запросы queryExistGaugesTable и queryExistCounterTable.
// В случае успеха логгирует информацию о том, что таблицы проверены/созданы.
// При возникновении ошибок — возвращает обёрнутую ошибку.
func ensureTablesExist(db *sql.DB, log *logrus.Entry) error {
	// Проверка существования (или создание) таблицы gauges
	gaugesQuery := queryExistGaugesTable
	_, err := db.Exec(gaugesQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы gauges: %w", err)
	}

	// Проверка существования (или создание) таблицы counters
	countersQuery := queryExistCounterTable
	_, err = db.Exec(countersQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы counters: %w", err)
	}

	log.Println("Таблицы gauges и counters проверены/созданы успешно.")
	return nil
}
