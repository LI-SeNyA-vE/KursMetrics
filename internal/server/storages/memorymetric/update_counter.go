// Package memorymetric предоставляет реализацию хранилища метрик в памяти.
package memorymetric

// UpdateCounter Обновление значения counter метрики (суммирование значений)
func (m *MetricStorage) UpdateCounter(name string, value int64) int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter[name] += value
	return m.counter[name]
}
