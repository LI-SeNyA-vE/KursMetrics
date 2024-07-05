package storage

import (
	"fmt"
)

var Metric = NewMetricStorage()

// Структура для хранения метрик в памяти
type MetricStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// Конструктор для создания нового экземпляра Metrictorage
func NewMetricStorage() *MetricStorage {
	return &MetricStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

// Обновление значения gauge метрики (Замена значения)
func (m *MetricStorage) UpdateGauge(name string, value float64) {
	m.Gauge[name] = value
}

// Обновление значения counter метрики (суммирование значений)
func (m *MetricStorage) UpdateCounter(name string, value int64) {
	m.Counter[name] += value
}

func (m *MetricStorage) GetAllGauges() map[string]float64 {
	return m.Gauge
}

func (m *MetricStorage) GetAllCounters() map[string]int64 {
	return m.Counter
}

func (m *MetricStorage) GetValue(typeMetric string, nameMetric string) (interface{}, error) {
	if typeMetric == "gauge" {
		if v, ok := m.Gauge[nameMetric]; ok {
			return v, nil
		}
	}
	if typeMetric == "counter" {
		if v, ok := m.Counter[nameMetric]; ok {
			return v, nil
		}
	}
	return nil, fmt.Errorf("нет метрики:%s, типа:%s", nameMetric, typeMetric)

	/* if v, ok := m.Gauge[nameMetric]; ok {
		return v, nil
	}
	if v, ok := m.Counter[nameMetric]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("нет метрики:%s, типа:%s", nameMetric, typeMetric) */
}
