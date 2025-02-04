/*
GetCounter и GetGauge возвращают значения отдельных метрик (counter или gauge)
из PostgreSQL. Если метрика не найдена — возвращают ошибку вида "<type> <name> not found".
При любых иных ошибках запросов также возвращают ошибку в обёрнутом формате.
*/

package postgresql

import (
	"database/sql"
	"fmt"
)

// GetGauge возвращает значение метрики типа gauge по её имени (name).
// Если метрика не найдена — возвращает ошибку "gauge <name> not found".
// Если возникла иная ошибка во время чтения из БД, она будет обёрнута и возвращена.
func (d *DataBase) GetGauge(name string) (*float64, error) {
	var value float64
	query := queryGetGauge

	err := d.db.QueryRow(query, name).Scan(&value)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("gauge %q not found", name)
	} else if err != nil {
		return nil, fmt.Errorf("failed to query gauge %q: %w", name, err)
	}

	return &value, nil
}
