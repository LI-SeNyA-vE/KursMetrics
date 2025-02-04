/*
Package postgresql предоставляет реализацию интерфейса MetricsStorage
(см. internal/funcserver/storages/metric.go) с использованием PostgreSQL
в качестве основного хранилища метрик.
*/
package postgresql

// UpdateGauge выполняет запрос queryUpdateGauge, вставляя или обновляя
// значение gauge (по сути, перезапись предыдущего значения).
// Возвращает новое значение метрики.
func (d *DataBase) UpdateGauge(name string, value float64) float64 {
	_, err := d.db.Exec(queryUpdateGauge, name, value)
	if err != nil {
		d.log.Printf("ошибка обновления gauge: %v", err)
	}
	return value
}
