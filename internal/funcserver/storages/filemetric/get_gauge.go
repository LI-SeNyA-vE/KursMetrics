package filemetric

import (
	"fmt"
)

func (s *FileStorage) GetGauge(name string) (*float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result, ok := s.data.Gauges[name]
	if !ok {
		return nil, fmt.Errorf("counter not found")
	}
	return &result, nil
}
