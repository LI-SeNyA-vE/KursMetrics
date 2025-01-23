package memorymetric

import "fmt"

func (m *MetricStorage) GetGauge(name string) (*float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.gauge[name]
	if !ok {
		return &v, fmt.Errorf("нет метрики:%s, типа: gauge", name)
	}
	return &v, nil
}
