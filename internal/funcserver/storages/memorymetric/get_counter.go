// Package memorymetric предоставляет реализацию хранилища метрик в памяти.
// GetCounter возвращают значения отдельных метрик counter из памяти.
// Если метрика не найдена — возвращают ошибку вида "<type> <name> not found".
// При любых иных ошибках запросов также возвращают ошибку в обёрнутом формате.
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
