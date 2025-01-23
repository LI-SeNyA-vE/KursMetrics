package postgresql

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

	// Проверяем rows.Err после итерации
	if err = rows.Err(); err != nil {
		d.log.Printf("Ошибка итерации gauges: %v", err)
		return nil
	}

	return result
}
