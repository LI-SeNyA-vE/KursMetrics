package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"

	storageMetric "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	log     *logrus.Entry
	cfg     config.Server
	storage storageMetric.MetricsStorage
}

func NewHandler(log *logrus.Entry, cfg config.Server, storage storageMetric.MetricsStorage) *Handler {
	return &Handler{
		log:     log,
		cfg:     cfg,
		storage: storage,
	}
}

// PostAddValue Обрабатывает полный url запрос. Если всё правильно сохраняет метрику в память
func (h *Handler) PostAddValue(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")
	countMetric := chi.URLParam(r, "countMetric")
	switch typeMetric { //Свитч для проверки что это запрос или gauge или counter
	case "gauge": //Если передано значение 'gauge'
		count, err := strconv.ParseFloat(countMetric, 64) //Проверка что переданно число и его можно перевети в float64
		if err != nil {                                   //Если не удалось перевести
			http.Error(w, "Это не Float", http.StatusBadRequest) //Вывод error-ки
			return                                               //
		}
		h.storage.UpdateGauge(nameMetric, count)
	case "counter": //Если передано значение 'counter'
		{
			count, err := strconv.ParseInt(countMetric, 10, 64) //Проверка что переданно число и его можно перевети в int64
			if err != nil {                                     //Если не удалось перевести
				http.Error(w, "Это не Float", http.StatusBadRequest) //Вывод error-ки
				return                                               //
			}
			h.storage.UpdateCounter(nameMetric, count)
		}
	default: //Если передано другое значение значение
		{
			http.Error(w, "Это не 'gauge' и не 'counter' запросы", http.StatusBadRequest) //Вывод error-ки
			return                                                                        //Выход из ServeHTTP
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK) //Отправляет ответ что всё ОК
}

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
		io.WriteString(w, fmt.Sprint(gauge))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	case "counter":
		counter, err := h.storage.GetCounter(nameMetric) // Запрашивает метрику, по данным из JSON
		if err != nil {
			h.log.Info(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		io.WriteString(w, fmt.Sprint(counter))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

// GetReceivingAllMetric Возвращает страничку со всеми метриками
func (h *Handler) GetReceivingAllMetric(w http.ResponseWriter, r *http.Request) {
	body := `
        <!DOCTYPE html>
        <html>
            <head>
                <title>All tuples</title>
            </head>
            <body>
            <table>
                <tr>
                    <td>Metric</td>
                    <td>Value</td>
                </tr>
    `
	listC := h.storage.GetAllCounters()
	for k, v := range listC {
		body = body + fmt.Sprintf("<tr>\n<td>%s</td>\n", k)
		body = body + fmt.Sprintf("<td>%v</td>\n</tr>\n", v)
	}

	listG := h.storage.GetAllGauges()
	for k, v := range listG {
		body = body + fmt.Sprintf("<tr>\n<td>%s</td>\n", k)
		body = body + fmt.Sprintf("<td>%v</td>\n</tr>\n", v)
	}

	body = body + " </table>\n </body>\n</html>"

	// respond to agent
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

// JSONValue Запрашивает метрику через JSON формат
func (h *Handler) JSONValue(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var metrics storageMetric.Metrics

	_, err := buf.ReadFrom(r.Body) //Читает данные из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &metrics) // Разбирает данные из массива byte в структуру "metrics"
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metrics.MType {
	case "gauge":
		metrics.Value, err = h.storage.GetGauge(metrics.ID) // Запрашивает метрику, по данным из JSON
		if err != nil {
			h.log.Info(err)
			http.Error(w, "не найдено", http.StatusNotFound)
			return
		}
	case "counter":
		metrics.Delta, err = h.storage.GetCounter(metrics.ID) // Запрашивает метрику, по данным из JSON
		if err != nil {
			h.log.Info(err)
			http.Error(w, "не найдено", http.StatusNotFound)
			return
		}
	default:
		http.Error(w, err.Error(), http.StatusBadRequest)
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

// JSONUpdate Обновляет метрику через JSON запрос
func (h *Handler) JSONUpdate(w http.ResponseWriter, r *http.Request) {
	var metrics storageMetric.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body) //Читает данные из тела запроса
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
		metric := h.storage.UpdateGauge(metrics.ID, metrics.Value) //Обновляет метрику
		metrics.Value = metric
	case "counter":
		metric := h.storage.UpdateCounter(metrics.ID, metrics.Delta) //Обновляет метрику
		metrics.Delta = metric
	default:
		http.Error(w, err.Error(), http.StatusBadRequest)
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

// Ping Кидает запрос в базу, для прорки её наличия
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", h.cfg.FlagDatabaseDsn)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// PostAddArrayMetrics Хендлер, который позволяет принимать массив метрик и сохранять его
func (h *Handler) PostAddArrayMetrics(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var arrayMetrics []storageMetric.Metrics

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
			h.storage.UpdateCounter(metrics.ID, metrics.Delta) //Обновляет метрику
		case "gauge":
			h.storage.UpdateGauge(metrics.ID, metrics.Value) //Обновляет метрику
		}
	}
	w.WriteHeader(http.StatusOK)
}
