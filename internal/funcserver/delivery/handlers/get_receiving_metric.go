/*
Package handlers содержит набор HTTP-обработчиков (Handler),
отвечающих за приём, обновление и вывод метрик.
*/
package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

// GetReceivingMetric обрабатывает GET-запрос на получение значения метрики
// по её типу (gauge или counter) и названию. Параметры извлекаются из
// URL с помощью chi.URLParam. Если метрика не найдена в хранилище –
// возвращается статус 404 (Not Found). В противном случае отдаётся
// её текущее значение в теле ответа.
func (h *Handler) GetReceivingMetric(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")
	typeMetric := chi.URLParam(r, "typeMetric")
	h.log.Info("Запрос с " + nameMetric + " " + typeMetric)

	switch typeMetric {
	case "gauge":
		// Запрос значения gauge-метрики по имени
		gauge, err := h.storage.GetGauge(nameMetric)
		if err != nil {
			h.log.Info(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		io.WriteString(w, fmt.Sprint(*gauge))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	case "counter":
		// Запрос значения counter-метрики по имени
		counter, err := h.storage.GetCounter(nameMetric)
		if err != nil {
			h.log.Info(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		io.WriteString(w, fmt.Sprint(*counter))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	default:
		// Если тип метрики не gauge или counter, возвращаем ошибку 400 (Bad Request)
		h.log.Infof("передан не 'gauge' и не 'counter' | url: %s", r.URL.Path)
		http.Error(w, "это не 'gauge' и не 'counter' запросы", http.StatusBadRequest)
		return
	}
}
