package storage

// Структура для хранения метрик в памяти
type MetricStorage struct {
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

var Metric = NewMetricStorage()

// Обновление значения gauge метрики (Замена значения)
func (m *MetricStorage) UpdateGauge(name string, value float64) {
	m.gauge[name] = value
}

// Обновление значения counter метрики (суммирование значений)
func (m *MetricStorage) UpdateCounter(name string, value int64) {
	m.counter[name] += value
}
