// Package postgresql содержит реализацию интерфейса MetricsStorage на базе PostgreSQL.
// Методы GetAllGauges возвращают все имеющиеся в базе gauge-метрики
// соответственно, в виде карт [имя_метрики]значение.
package postgresql

// GetAllGauges выполняет запрос queryGetAllGauges к базе данных,
// считывая все метрики типа gauge и возвращая их в виде карты map[string]float64.
// В случае ошибки чтения или парсинга строк пишет соответствующую информацию в лог
// и возвращает nil.
func (d *DataBase) GetAllGauges() map[string]float64 {
	rows, err := d.db.Query(queryGetAllGauges)
	if err != nil {
		d.log.Printf("ошибка получения gauges: %v", err)
		return nil
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			d.log.Printf("ошибка чтения строки gauge: %v", err)
			continue
		}
		result[name] = value
	}

	// Проверяем rows.Err() после итерации
	if err = rows.Err(); err != nil {
		d.log.Printf("Ошибка итерации gauges: %v", err)
		return nil
	}

	return result
}
