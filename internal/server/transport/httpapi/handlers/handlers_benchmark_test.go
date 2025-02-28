package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/storages"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/storages/memorymetric"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/transport/httpapi/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkPostAddArrayMetrics проверяет производительность хендлера PostAddArrayMetrics.
// Он эмулирует POST-запрос с JSON-массивом метрик и многократно вызывает хендлер.
func BenchmarkPostAddArrayMetrics(b *testing.B) {
	// 1. Готовим "заглушки" (mock/фейки) для логгера, конфига и хранилища
	log := logger.NewLogger() // можно упростить, если слишком «тяжёлый» логгер
	cfg := servercfg.Server{} // пустой конфиг, если тестируем чистую логику
	storage := memorymetric.NewMetricStorage()

	// 2. Создаём сам Handler
	h := handlers.NewHandler(log, cfg, storage)

	// 3. Подготовим пример данных (JSON-массив метрик),
	//    чтобы не формировать их на каждом шаге внутри цикла b.N
	arrayMetrics := []storages.Metrics{
		{
			ID:    "CounterMetric",
			MType: "counter",
			Delta: ptrInt64(123),
		},
		{
			ID:    "GaugeMetric",
			MType: "gauge",
			Value: ptrFloat64(45.67),
		},
		// при желании можно добавить больше метрик в массив
	}

	// Сериализуем в JSON
	data, err := json.Marshal(arrayMetrics)
	if err != nil {
		b.Fatalf("Ошибка маршалинга JSON: %v", err)
	}

	// 4. Сбрасываем таймер перед самим тестом, чтобы не учитывать время на подготовку
	b.ResetTimer()

	// 5. Запускаем цикл бенчмарка
	for i := 0; i < b.N; i++ {
		// Каждый раз создаём новый Recorder и новый Request
		req := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBuffer(data))
		w := httptest.NewRecorder()

		// 6. Вызываем тестируемую функцию
		h.PostAddArrayMetrics(w, req)

		//b.StopTimer()
		//if w.Result().StatusCode != http.StatusOK {
		//	b.Errorf("unexpected status code: %d", w.Result().StatusCode)
		//}
		//b.StartTimer()
	}
}

// ptrInt64 — маленький вспомогательный хелпер
func ptrInt64(v int64) *int64 {
	return &v
}
func ptrFloat64(v float64) *float64 {
	return &v
}
