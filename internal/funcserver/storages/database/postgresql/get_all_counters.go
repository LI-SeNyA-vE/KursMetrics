/*
Package postgresql содержит реализацию интерфейса MetricsStorage на базе PostgreSQL.
Методы GetAllCounters и GetAllGauges возвращают все имеющиеся в базе counter- и gauge-метрики
соответственно, в виде карт [имя_метрики]значение.
*/
package postgresql

// GetAllCounters выполняет запрос queryGetAllCounters к базе данных,
// считывая все метрики типа counter и возвращая их в виде карты map[string]int64.
// В случае ошибки чтения или парсинга строк пишет соответствующую информацию в лог
// и возвращает nil.
func (d *DataBase) GetAllCounters() map[string]int64 {
	rows, err := d.db.Query(queryGetAllCounters)
	if err != nil {
		d.log.Printf("ошибка получения counters: %v", err)
		return nil
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var name string
		var value int64
		if err := rows.Scan(&name, &value); err != nil {
			d.log.Printf("ошибка чтения строки counter: %v", err)
			continue
		}
		result[name] = value
	}

	// Проверяем rows.Err() для выявления ошибок итерации
	if err = rows.Err(); err != nil {
		d.log.Printf("Ошибка итерации counters: %v", err)
		return nil
	}

	return result
}
