/*
Package storages определяет общий интерфейс для работы с хранилищем метрик (MetricsStorage)
и структуру Metrics, используемую при передаче/приёме данных о метриках (как gauge, так и counter)
в JSON-формате.
*/
package storages

// MetricsStorage описывает методы для создания и обновления метрик (gauge, counter),
// а также для их получения как по отдельности, так и списком. Дополнительно содержит метод
// LoadMetric, позволяющий загружать метрики из внешнего хранилища (например, из файла или БД)
// при инициализации.
type MetricsStorage interface {
	UpdateGauge(name string, value float64) float64
	UpdateCounter(name string, value int64) int64
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
	GetGauge(name string) (*float64, error)
	GetCounter(name string) (*int64, error)
	LoadMetric() error
}

// Metrics представляет структуру для передачи информации о метрике:
//   - ID: название метрики,
//   - MType: тип метрики (gauge или counter),
//   - Delta: значение для counter,
//   - Value: значение для gauge.
//
// Поля Delta и Value сделаны указателями, чтобы в JSON не отправлялось лишнее,
// а также была возможность разграничивать типы в одном объекте.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // "gauge" или "counter"
	Delta *int64   `json:"delta,omitempty"` // значение метрики, если MType == "counter"
	Value *float64 `json:"value,omitempty"` // значение метрики, если MType == "gauge"
}
