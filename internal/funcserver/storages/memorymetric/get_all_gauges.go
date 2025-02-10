// Package memorymetric предоставляет реализацию хранилища метрик в памяти.
// Методы GetAllGauges возвращают все сохранённые в памяти gauge-метрики
// соответственно, в виде карт [имя_метрики]значение.
package memorymetric

func (m *MetricStorage) GetAllGauges() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.gauge
}
