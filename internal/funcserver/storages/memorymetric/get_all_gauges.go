package memorymetric

func (m *MetricStorage) GetAllGauges() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.gauge
}
