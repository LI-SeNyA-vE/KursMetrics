package memoryMetric

import (
	"fmt"
	"sync"
)

// Структура для хранения метрик в памяти
type MetricStorage struct {
	mu      sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

// Конструктор для создания нового экземпляра Metrictorage
func NewMetricStorage() *MetricStorage {
	return &MetricStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

// Обновление значения gauge метрики (Замена значения)
func (m *MetricStorage) UpdateGauge(name string, value float64) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauge[name] = value
	return m.gauge[name]
}

// Обновление значения counter метрики (суммирование значений)
func (m *MetricStorage) UpdateCounter(name string, value int64) int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter[name] += value
	return m.counter[name]
}

func (m *MetricStorage) GetAllGauges() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.gauge
}

func (m *MetricStorage) GetAllCounters() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counter
}

func (m *MetricStorage) GetGauge(name string) (*float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.gauge[name]
	if !ok {
		return &v, fmt.Errorf("нет метрики:%s, типа: gauge", name)
	}
	return &v, nil
}

func (m *MetricStorage) GetCounter(name string) (*int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.counter[name]
	if !ok {
		return &v, fmt.Errorf("нет метрики:%s, типа: counter", name)
	}
	return &v, nil
}

func (m *MetricStorage) LoadMetric() (err error) {
	return err
}
