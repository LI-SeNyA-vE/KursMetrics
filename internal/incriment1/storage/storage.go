package storage

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

var Metric = NewMetricStorage()

// Обновление значения gauge метрики (Замена значения)
func (m *MetricStorage) UpdateGauge(name string, value float64) {
	m.Gauge[name] = value
}

// Обновление значения counter метрики (суммирование значений)
func (m *MetricStorage) UpdateCounter(name string, value int64) {
	m.Counter[name] += value
}
