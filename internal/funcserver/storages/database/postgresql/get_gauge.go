package postgresql

import (
	"database/sql"
	"fmt"
)

// GetGauge возвращает значение метрики типа gauge
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
