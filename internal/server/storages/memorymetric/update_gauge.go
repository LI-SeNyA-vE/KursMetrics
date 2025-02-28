// Package memorymetric предоставляет реализацию хранилища метрик в памяти.
package memorymetric

// UpdateGauge Обновление значения gauge метрики (Замена значения)
func (m *MetricStorage) UpdateGauge(name string, value float64) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauge[name] = value
	return m.gauge[name]
}
