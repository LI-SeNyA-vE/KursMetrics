package filemetric

func (s *FileStorage) GetAllCounters() map[string]int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make(map[string]int64, len(s.data.Counters))
	for k, v := range s.data.Counters {
		result[k] = v
	}
	return result
}
