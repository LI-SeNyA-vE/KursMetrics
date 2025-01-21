package storages

type MetricsStorage interface {
	UpdateGauge(name string, value float64) float64
	UpdateCounter(name string, value int64) int64
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
	GetGauge(name string) (*float64, error)
	GetCounter(name string) (*int64, error)
	LoadMetric() error
}

// Структура метрики для отправки JSON
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
