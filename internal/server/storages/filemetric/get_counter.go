// Package filemetric предоставляет реализацию хранилища метрик в локальном хранилище.
// GetCounter возвращают значения отдельных метрик counter
// из файла. Если метрика не найдена — возвращают ошибку вида "<type> <name> not found".
// При любых иных ошибках запросов также возвращают ошибку в обёрнутом формате.
package filemetric

import (
	"fmt"
)

func (s *FileStorage) GetCounter(name string) (*int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result, ok := s.data.Counters[name]
	if !ok {
		return nil, fmt.Errorf("counter not found")
	}
	return &result, nil
}
