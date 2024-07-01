package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCorrectPostRequest(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name        string
		typeMetric  string
		nameMetric  string
		countMetric string
		want        want
	}{
		{"Gauge метрика", "gauge", "validGaugeName", "1.23", want{http.StatusOK, "text/plain"}},
		{"Counter метрика", "counter", "validCounterName", "123", want{http.StatusOK, "text/plain"}},
		// Добавьте другие тестовые случаи по необходимости
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/update/%s/%s/%s", tt.typeMetric, tt.nameMetric, tt.countMetric), nil)
			rw := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", PostAddValue)
			r.ServeHTTP(rw, req)

			res := rw.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestCorrectGetRequest(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name        string
		typeMetric  string
		nameMetric  string
		countMetric string
		want        want
	}{
		{"Gauge метрика", "gauge", "validGaugeName", "1.23", want{http.StatusOK, "text/plain"}},
		{"Counter метрика", "counter", "validCounterName", "123", want{http.StatusOK, "text/plain"}},
		// Добавьте другие тестовые случаи по необходимости
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/update/%s/%s/%s", tt.typeMetric, tt.nameMetric, tt.countMetric), nil)
			rw := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", PostAddValue)
			r.Get("/update/{typeMetric}/{nameMetric}", GetReceivingMetric)
			r.ServeHTTP(rw, req)

			res := rw.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestAllValue(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name        string
		typeMetric  string
		nameMetric  string
		countMetric string
		want        want
	}{
		{"Gauge метрика", "gauge", "validGaugeName", "1.23", want{http.StatusOK, "text/plain"}},
		{"Counter метрика", "counter", "validCounterName", "123", want{http.StatusOK, "text/plain"}},
		// Добавьте другие тестовые случаи по необходимости
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/update/%s/%s/%s", tt.typeMetric, tt.nameMetric, tt.countMetric), nil)
			rw := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", PostAddValue)
			r.Get("/", GetReceivingAllMetric)
			r.ServeHTTP(rw, req)

			res := rw.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
