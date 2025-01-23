package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"net/http"
)

// JSONUpdate Обновляет метрику через JSON запрос
func (h *Handler) JSONUpdate(w http.ResponseWriter, r *http.Request) {
	var metrics storages.Metrics
	var buf bytes.Buffer
	var err error

	_, err = buf.ReadFrom(r.Body) //Читает данные из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &metrics) // Разбирает данные из массива byte в структуру "metrics"
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Проверка на тип с последующим вызовом нужной функции
	switch metrics.MType {
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
		metric := h.storage.UpdateGauge(metrics.ID, *metrics.Value) //Обновляет метрику
		metrics.Value = &metric
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
		metric := h.storage.UpdateCounter(metrics.ID, *metrics.Delta) //Обновляет метрику
		metrics.Delta = &metric
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
	default:
		h.log.Debug("JSON ответа:",
			"\n{",
			"\n  ID: ", metrics.ID,
			"\n  MType: ", metrics.MType,
			"\n}\n")
		http.Error(w, "нет такого типа", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(metrics) // Запаковывает/собирает данные в массив byte
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
