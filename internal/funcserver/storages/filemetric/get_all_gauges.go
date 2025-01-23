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
