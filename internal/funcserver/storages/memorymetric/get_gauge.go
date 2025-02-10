// Package memorymetric предоставляет реализацию хранилища метрик в памяти.
// GetGauge возвращают значения отдельных метрик gauge из памяти.
// Если метрика не найдена — возвращают ошибку вида "<type> <name> not found".
// При любых иных ошибках запросов также возвращают ошибку в обёрнутом формате.
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
