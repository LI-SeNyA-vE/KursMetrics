/*
Package handlers содержит функцию PostAddArrayMetrics, позволяющую обработать
JSON-массив метрик (MType gauge или counter) и сохранить каждую из них
в соответствующем хранилище.
*/
package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/storages"
	"net/http"
)

// PostAddArrayMetrics обрабатывает JSON-массив метрик, переданных в теле запроса.
// Для каждой метрики (gauge или counter) вызывается соответствующий метод обновления
// в хранилище. После успешной обработки возвращается статус 200 (OK). Если при чтении
// тела запроса или его парсинге в массив метрик возникнут ошибки, клиенту будет
// отдан код 400 (Bad Request).
func (h *Handler) PostAddArrayMetrics(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var arrayMetrics []storages.Metrics

	_, err := buf.ReadFrom(r.Body) // Читаем данные из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &arrayMetrics) // Парсим массив метрик из JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Errorf("ошибка анмаршела переданного тела: %v", err.Error())
		return
	}

	// Обходим полученный срез метрик и обновляем каждую по типу (counter/gauge)
	for _, metrics := range arrayMetrics {
		switch metrics.MType {
		case "counter":
			if metrics.Delta != nil {
				h.log.Debug("JSON запроса:",
					"\n{",
					"\n  ID: ", metrics.ID,
					"\n  MType: ", metrics.MType,
					"\n  Delta: ", *metrics.Delta,
					"\n}\n")
			} else {
				h.log.Debug("JSON запроса:",
					"\n{",
					"\n  ID: ", metrics.ID,
					"\n  MType: ", metrics.MType,
					"\n  Delta: nil",
					"\n}\n")
			}
			res := h.storage.UpdateCounter(metrics.ID, *metrics.Delta)
			metrics.Delta = &res
			if metrics.Delta != nil {
				h.log.Debug("JSON ответа:",
					"\n{",
					"\n  ID: ", metrics.ID,
					"\n  MType: ", metrics.MType,
					"\n  Delta: ", *metrics.Delta,
					"\n}\n")
			} else {
				h.log.Debug("JSON ответа:",
					"\n{",
					"\n  ID: ", metrics.ID,
					"\n  MType: ", metrics.MType,
					"\n  Delta: nil",
					"\n}\n")
			}

		case "gauge":
			if metrics.Value != nil {
				h.log.Debug("JSON запроса:",
					"\n{",
					"\n  ID: ", metrics.ID,
					"\n  MType: ", metrics.MType,
					"\n  Value: ", *metrics.Value,
					"\n}\n")
			} else {
				h.log.Debug("JSON запроса:",
					"\n{",
					"\n  ID: ", metrics.ID,
					"\n  MType: ", metrics.MType,
					"\n  Value: nil",
					"\n}\n")
			}
			res := h.storage.UpdateGauge(metrics.ID, *metrics.Value)
			metrics.Value = &res
			if metrics.Value != nil {
				h.log.Debug("JSON ответа:",
					"\n{",
					"\n  ID: ", metrics.ID,
					"\n  MType: ", metrics.MType,
					"\n  Value: ", *metrics.Value,
					"\n}\n")
			} else {
				h.log.Debug("JSON ответа:",
					"\n{",
					"\n  ID: ", metrics.ID,
					"\n  MType: ", metrics.MType,
					"\n  Value: nil",
					"\n}\n")
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}
