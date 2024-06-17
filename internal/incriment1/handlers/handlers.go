package handlers

import (
	"fmt"
	"io"
	"net/http"

	serverMetric "github.com/LI-SeNyA-vE/KursMetrics/internal/incriment1/server"
	storageMetric "github.com/LI-SeNyA-vE/KursMetrics/internal/incriment1/storage"
	"github.com/go-chi/chi/v5"
)

func CorrectPostRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL: ", r.URL.Path)
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")
	countMetric := chi.URLParam(r, "countMetric")
	if serverMetric.ValidationTypeMetric(typeMetric, nameMetric, countMetric, w) {
		return
	}
	w.WriteHeader(http.StatusOK) //Отправляет ответ что всё ОК
}

func CorrectGetRequest(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")
	typeMetric := chi.URLParam(r, "typeMetric")
	value, err := storageMetric.Metric.GetValue(typeMetric, nameMetric)
	if !err {
		io.WriteString(w, fmt.Sprint(value))
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func AllValue(w http.ResponseWriter, r *http.Request) {
	v := storageMetric.Metric
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprint(v))
}
