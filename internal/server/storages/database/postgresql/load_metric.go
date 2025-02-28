// Package postgresql содержит реализацию интерфейса MetricsStorage на базе PostgreSQL.
package postgresql

// LoadMetric на текущий момент не реализует логику (возвращает nil).
// В будущем можно добавить загрузку метрик из таблиц при старте приложения,
// если требуется первичная инициализация из БД.
func (d *DataBase) LoadMetric() (err error) {
	return err
}
