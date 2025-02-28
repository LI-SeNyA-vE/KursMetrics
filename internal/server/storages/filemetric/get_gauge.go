// Package filemetric предоставляет реализацию хранилища метрик в локальном хранилище.
// GetGauge возвращают значения отдельных метрик gauge из файла.
// Если метрика не найдена — возвращают ошибку вида "<type> <name> not found".
// При любых иных ошибках запросов также возвращают ошибку в обёрнутом формате.
package filemetric

import (
	"fmt"
)

func (s *FileStorage) GetGauge(name string) (*float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result, ok := s.data.Gauges[name]
	if !ok {
		return nil, fmt.Errorf("counter %s not found", name)
	}
	return &result, nil
}
