package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"net/http"
)

// PostAddArrayMetrics Хендлер, который позволяет принимать массив метрик и сохранять его
func (h *Handler) PostAddArrayMetrics(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var arrayMetrics []storages.Metrics

	_, err := buf.ReadFrom(r.Body) //Читает данные из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &arrayMetrics) // Разбирает данные из массива byte в массив структур "metrics"
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
			res := h.storage.UpdateCounter(metrics.ID, *metrics.Delta) //Обновляет метрику
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
			res := h.storage.UpdateGauge(metrics.ID, *metrics.Value) //Обновляет метрику
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
