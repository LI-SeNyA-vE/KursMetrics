package incriment1

import "net/http"

// Структура для хранения метрик в памяти
type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

// Конструктор для создания нового экземпляра MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

// Обновление значения gauge метрики (Замена значения)
func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.gauge[name] = value
}

// Обновление значения counter метрики (суммирование значений)
func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.counter[name] += value
}

// Обработка HTTP запросов
func (m *MemStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
