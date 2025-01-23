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
