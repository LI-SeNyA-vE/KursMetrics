package memorymetric

// Конструктор для создания нового экземпляра Metrictorage
func NewMetricStorage() *MetricStorage {
	return &MetricStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}
