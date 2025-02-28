// Package filemetric предоставляет реализацию хранилища метрик в локальном хранилище.
// Методы GetAllGauges возвращают все имеющиеся в файле gauge-метрики
// соответственно, в виде карт [имя_метрики]значение.
package filemetric

func (s *FileStorage) GetAllGauges() map[string]float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make(map[string]float64, len(s.data.Gauges))
	for k, v := range s.data.Gauges {
		result[k] = v
	}
	return result
}
