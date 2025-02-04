/*
Package postgresql предоставляет реализацию интерфейса MetricsStorage
(см. internal/funcserver/storages/metric.go) с использованием PostgreSQL
в качестве основного хранилища метрик.
*/
package postgresql

// LoadMetric на текущий момент не реализует логику (возвращает nil).
// В будущем можно добавить загрузку метрик из таблиц при старте приложения,
// если требуется первичная инициализация из БД.
func (d *DataBase) LoadMetric() (err error) {
	return err
}
