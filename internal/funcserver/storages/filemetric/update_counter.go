// Package filemetric предоставляет реализацию хранилища метрик в локальном хранилище.
package filemetric

func (s *FileStorage) UpdateCounter(name string, value int64) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Counters[name] += value
	s.saveToFile()
	return s.data.Counters[name]
}
