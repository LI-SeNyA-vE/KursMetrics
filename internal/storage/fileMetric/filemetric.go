package fileMetric

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type FileStorage struct {
	filePath string
	mu       sync.Mutex
	data     struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]int64   `json:"counters"`
	}
}

// NewFileStorage — конструктор FileStorage
func NewFileStorage(filePath string) (*FileStorage, error) {
	storage := &FileStorage{
		filePath: filePath,
		data: struct {
			Gauges   map[string]float64 `json:"gauges"`
			Counters map[string]int64   `json:"counters"`
		}{
			Gauges:   make(map[string]float64),
			Counters: make(map[string]int64),
		},
	}

	// Загружаем данные из файла, если он существует
	if _, err := os.Stat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&storage.data); err != nil {
			return nil, err
		}
	}
	return storage, nil
}

func (s *FileStorage) UpdateGauge(name string, value float64) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Gauges[name] = value
	s.saveToFile()
	return s.data.Gauges[name]
}

func (s *FileStorage) UpdateCounter(name string, value int64) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Counters[name] += value
	s.saveToFile()
	return s.data.Counters[name]
}

func (s *FileStorage) GetAllGauges() map[string]float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make(map[string]float64, len(s.data.Gauges))
	for k, v := range s.data.Gauges {
		result[k] = v
	}
	return result
}

func (s *FileStorage) GetAllCounters() map[string]int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make(map[string]int64, len(s.data.Counters))
	for k, v := range s.data.Counters {
		result[k] = v
	}
	return result
}

func (s *FileStorage) GetGauge(name string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	//TODO сделай меня
	return s.data.Gauges[name], nil
}

func (s *FileStorage) GetCounter(name string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	//TODO сделай меня
	return s.data.Counters[name], nil
}

func (s *FileStorage) saveToFile() {
	file, err := os.Create(s.filePath)
	if err != nil {
		// Обработка ошибки сохранения
		return
	}
	defer file.Close()
	_ = json.NewEncoder(file).Encode(s.data)
}

func (s *FileStorage) LoadMetric() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("Ошибка чтения файла: %s", err)
	}
	if err = json.Unmarshal(res, &s.data); err != nil {
		return fmt.Errorf("Ошибка Unmarshal: %s", err)
	}
	return err
}
