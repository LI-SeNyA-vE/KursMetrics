package memorymetric

func (m *MetricStorage) GetAllCounters() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counter
}
