package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

// GetReceivingMetric Позваляет получить знаачение метрики по данным: Тип метрики и Название метрики
func (h *Handler) GetReceivingMetric(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")
	typeMetric := chi.URLParam(r, "typeMetric")
	h.log.Info("Запрос с " + nameMetric + " " + typeMetric)

	switch typeMetric {
	case "gauge":
		gauge, err := h.storage.GetGauge(nameMetric) // Запрашивает метрику, по данным из JSON
		if err != nil {
			h.log.Info(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		io.WriteString(w, fmt.Sprint(*gauge))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	case "counter":
		counter, err := h.storage.GetCounter(nameMetric) // Запрашивает метрику, по данным из JSON
		if err != nil {
			h.log.Info(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		io.WriteString(w, fmt.Sprint(*counter))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	default:
		h.log.Infof("передан не 'gauge' и не 'counter' | url: %s", r.URL.Path)
		http.Error(w, "это не 'gauge' и не 'counter' запросы", http.StatusBadRequest)
		return
	}
}
