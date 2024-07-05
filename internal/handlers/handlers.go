package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	storageMetric "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/go-chi/chi/v5"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func PostAddValue(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL: ", r.URL.Path)
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
		storageMetric.Metric.UpdateGauge(nameMetric, count)
	case "counter": //Если передано значение 'counter'
		{
			count, err := strconv.ParseInt(countMetric, 10, 64) //Проверка что переданно число и его можно перевети в int64
			if err != nil {                                     //Если не удалось перевести
				http.Error(w, "Это не Float", http.StatusBadRequest) //Вывод error-ки
				return                                               //
			}
			storageMetric.Metric.UpdateCounter(nameMetric, count)
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

func GetReceivingMetric(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")
	typeMetric := chi.URLParam(r, "typeMetric")
	value, err := storageMetric.Metric.GetValue(typeMetric, nameMetric)
	if err != nil {
		io.WriteString(w, fmt.Sprint(value))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
}

func GetReceivingAllMetric(w http.ResponseWriter, r *http.Request) {
	gauges := storageMetric.Metric.GetAllGauges()
	counters := storageMetric.Metric.GetAllCounters()

	data := storageMetric.MetricStorage{
		Gauge:   gauges,
		Counter: counters,
	}

	tmplPath := "../../internal/templates/index.html"

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func JsonValue(w http.ResponseWriter, r *http.Request) {
	var metrics Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := storageMetric.Metric.GetValue(metrics.MType, metrics.ID)
	if err != nil {
		http.Error(w, "не найдено", http.StatusNotFound)
		return
	}

	switch metrics.MType {
	case "counter":
		v := metric.(int64)
		metrics.Delta = &v
	case "gauge":
		v := metric.(float64)
		metrics.Value = &v
	}

	resp, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func JsonUpdate(w http.ResponseWriter, r *http.Request) {
	var metrics Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metrics.MType {
	case "counter":
		storageMetric.Metric.UpdateCounter(metrics.ID, *metrics.Delta)
	case "gauge":
		storageMetric.Metric.UpdateGauge(metrics.ID, *metrics.Value)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
