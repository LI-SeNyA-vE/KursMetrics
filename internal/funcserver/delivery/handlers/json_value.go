/*
Package handlers содержит набор HTTP-обработчиков (Handler),
отвечающих за приём, обновление и вывод метрик.
JSONValue обрабатывает запросы на получение значения конкретной метрики в формате JSON.
*/
package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"net/http"
)

// JSONValue принимает JSON-запрос, в котором указывается тип метрики (gauge/counter) и её ID.
// В ответ возвращает текущее значение этой метрики (Value или Delta) также в JSON-формате.
// Если в хранилище метрика отсутствует – возвращает статус 404 (Not Found).
func (h *Handler) JSONValue(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var metrics storages.Metrics

	_, err := buf.ReadFrom(r.Body) // Читаем данные из тела запроса
	if err != nil {
		h.log.Infof("ошибка чтения запроса %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразуем прочитанные данные в структуру Metrics
	err = json.Unmarshal(buf.Bytes(), &metrics)
	if err != nil {
		h.log.Infof("ошибка анмаршлинга запроса в JSON, %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// В зависимости от типа метрики (MType), получаем её значение из хранилища
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
		// Получаем gauge-метрику из хранилища
		metrics.Value, err = h.storage.GetGauge(metrics.ID)
		if err != nil {
			h.log.Info(err)
			http.Error(w, "не найдено", http.StatusNotFound)
			return
		}
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
				"\n}\n")
		} else {
			h.log.Debug("JSON запроса:",
				"\n{",
				"\n  ID: ", metrics.ID,
				"\n  Delta: nil",
				"\n}\n")
		}
		// Получаем counter-метрику из хранилища
		metrics.Delta, err = h.storage.GetCounter(metrics.ID)
		if err != nil {
			h.log.Info(err)
			http.Error(w, "не найдено", http.StatusNotFound)
			return
		}
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
		// Если тип не поддерживается (не gauge и не counter)
		h.log.Debug("JSON запроса и ответа:",
			"\n{",
			"\n  ID: ", metrics.ID,
			"\n  MType: ", metrics.MType,
			"\n}\n")
		http.Error(w, "не найдено", http.StatusNotFound)
		return
	}

	// Формируем ответ с актуальным значением метрики
	resp, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем результат в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
