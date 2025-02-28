// Package postgresql содержит реализацию интерфейса MetricsStorage на базе PostgreSQL.
// GetCounter возвращают значения отдельных метрик counter
// из PostgreSQL. Если метрика не найдена — возвращают ошибку вида "<type> <name> not found".
// При любых иных ошибках запросов также возвращают ошибку в обёрнутом формате.
package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetCounter возвращает значение метрики типа counter по её имени (name).
// Если метрика не найдена — возвращает ошибку "counter <name> not found".
// Если возникла иная ошибка во время чтения из БД, она будет обёрнута и возвращена.
func (d *DataBase) GetCounter(name string) (*int64, error) {
	var value int64
	query := queryGetCounter

	err := d.db.QueryRow(query, name).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("counter %q not found", name)
	} else if err != nil {
		return nil, fmt.Errorf("failed to query counter %q: %w", name, err)
	}

	return &value, nil
}
