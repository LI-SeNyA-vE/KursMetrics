package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"net/http"
)

// JSONValue Запрашивает метрику через JSON формат
func (h *Handler) JSONValue(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var metrics storages.Metrics

	_, err := buf.ReadFrom(r.Body) //Читает данные из тела запроса
	if err != nil {
		h.log.Infof("ошибка чтения запроса %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &metrics) // Разбирает данные из массива byte в структуру "metrics"
	if err != nil {
		h.log.Infof("ошибка анмаршлинга запроса в JSON, %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metrics.MType {
	case "gauge":
		if metrics.Value != nil {
			h.log.Debug("JSON запроса:",
				"\n{",
				"\n  ID: ", metrics.ID,
				"\n  MType: ", metrics.MType,
				"\n}\n")
		} else {
			h.log.Debug("JSON запроса:",
				"\n{",
				"\n  ID: ", metrics.ID,
				"\n  MType: ", metrics.MType,
				"\n}\n")
		}
		metrics.Value, err = h.storage.GetGauge(metrics.ID) // Запрашивает метрику, по данным из JSON
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
		if err != nil {
			h.log.Info(err)
			http.Error(w, "не найдено", http.StatusNotFound)
			return
		}
	case "counter":
		if metrics.Delta != nil {
			h.log.Debug("JSON запроса:",
				"\n{",
				"\n  ID: ", metrics.ID,
				"\n  MType: ", metrics.MType,
				"\n}\n")
		} else {
			h.log.Debug("JSON запроса:",
				"\n{",
				"\n  ID: ", metrics.ID,
				"\n  Delta: nil",
				"\n}\n")
		}
		metrics.Delta, err = h.storage.GetCounter(metrics.ID) // Запрашивает метрику, по данным из JSON
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
		if err != nil {
			h.log.Info(err)
			http.Error(w, "не найдено", http.StatusNotFound)
			return
		}
	default:
		h.log.Debug("JSON запроса и ответа:",
			"\n{",
			"\n  ID: ", metrics.ID,
			"\n  MType: ", metrics.MType,
			"\n}\n")
		http.Error(w, "не найдено", http.StatusNotFound)
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
