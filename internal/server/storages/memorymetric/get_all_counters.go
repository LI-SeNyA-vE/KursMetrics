// Package memorymetric предоставляет реализацию хранилища метрик в памяти.
// Методы GetAllCounters возвращают все сохранённые в памяти counter-метрики
// соответственно, в виде карт [имя_метрики]значение.
package memorymetric

func (m *MetricStorage) GetAllCounters() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counter
}
