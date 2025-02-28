/*
Package handlers содержит набор HTTP-обработчиков (Handler),
отвечающих за приём, обновление и вывод метрик.
JSONUpdate обрабатывает JSON-запросы для обновления отдельных метрик.
*/
package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/storages"
	"net/http"
)

// JSONUpdate читает данные из тела запроса в формате JSON, распаковывая их
// в структуру storages.Metrics. В зависимости от поля MType (gauge/counter)
// обновляет соответствующую метрику в хранилище. Затем формирует ответ
// также в формате JSON с актуальным значением обновлённой метрики. Если тип
// нераспознан (не gauge и не counter), вернёт статус 400 (Bad Request).
func (h *Handler) JSONUpdate(w http.ResponseWriter, r *http.Request) {
	var metrics storages.Metrics
	var buf bytes.Buffer
	var err error

	_, err = buf.ReadFrom(r.Body) // Читаем данные из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразуем массив байт из буфера в структуру Metrics
	err = json.Unmarshal(buf.Bytes(), &metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, какой тип метрики нужно обновить
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
		// Обновляем gauge-метрику
		metric := h.storage.UpdateGauge(metrics.ID, *metrics.Value)
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
		// Обновляем counter-метрику
		metric := h.storage.UpdateCounter(metrics.ID, *metrics.Delta)
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
		// Если тип метрики не поддерживается
		h.log.Debug("JSON ответа:",
			"\n{",
			"\n  ID: ", metrics.ID,
			"\n  MType: ", metrics.MType,
			"\n}\n")
		http.Error(w, "нет такого типа", http.StatusBadRequest)
		return
	}

	// Формируем ответ с обновлённым значением метрики
	resp, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем результат
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
