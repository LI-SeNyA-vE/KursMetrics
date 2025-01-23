package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetCounter возвращает значение метрики типа counter
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
