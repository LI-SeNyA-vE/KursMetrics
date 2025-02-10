// Package postgresql содержит реализацию интерфейса MetricsStorage на базе PostgreSQL.
//UpdateCounter выполняет запрос queryUpdateCounter, вставляя или обновляя

package postgresql

// UpdateCounter выполняет запрос queryUpdateCounter, вставляя или обновляя
// значение counter с учётом инкремента (суммирование предыдущего и нового).
// Затем повторно запрашивает текущее значение метрики из БД. Если при запросе
// возникла ошибка, возвращает 0.
func (d *DataBase) UpdateCounter(name string, value int64) int64 {
	_, err := d.db.Exec(queryUpdateCounter, name, value)
	if err != nil {
		d.log.Printf("ошибка обновления counter: %v", err)
	}

	result, err := d.GetCounter(name)
	if err != nil {
		return 0
	}

	return *result
}
