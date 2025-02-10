// Package memorymetric предоставляет реализацию хранилища метрик в памяти.
// Модель структуры необходимая для правильной работы приложения
package memorymetric

import "sync"

// MetricStorage Структура для хранения метрик в памяти
type MetricStorage struct {
	mu      sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}
