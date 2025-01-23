package memorymetric

import "fmt"

func (m *MetricStorage) GetCounter(name string) (*int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.counter[name]
	if !ok {
		return &v, fmt.Errorf("нет метрики:%s, типа: counter", name)
	}
	return &v, nil
}
