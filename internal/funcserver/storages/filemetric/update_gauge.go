package filemetric

func (s *FileStorage) UpdateGauge(name string, value float64) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Gauges[name] = value
	s.saveToFile()
	return s.data.Gauges[name]
}
