package memorymetric

import "sync"

// Структура для хранения метрик в памяти
type MetricStorage struct {
	mu      sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}
