package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// PostAddValue Обрабатывает полный url запрос. Если всё правильно сохраняет метрику в память
func (h *Handler) PostAddValue(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")
	countMetric := chi.URLParam(r, "countMetric")
	switch typeMetric { //Свитч для проверки что это запрос или gauge или counter
	case "gauge": //Если передано значение 'gauge'
		count, err := strconv.ParseFloat(countMetric, 64) //Проверка что переданно число и его можно перевети в float64
		if err != nil {                                   //Если не удалось перевести
			h.log.Infof("передано не число или его нельзя перевести в float64 | url: %v | оригинальная ощибка: %s", r.URL.Path, err)
			http.Error(w, "это не Float", http.StatusBadRequest) //Вывод error-ки
			return
		}
		h.storage.UpdateGauge(nameMetric, count)
	case "counter": //Если передано значение 'counter'
		count, err := strconv.ParseInt(countMetric, 10, 64) //Проверка что переданно число и его можно перевети в int64
		if err != nil {                                     //Если не удалось перевести
			//Если не удалось перевести
			h.log.Infof("передано не число или его нельзя перевести в int64 | url: %v | оригинальная ощибка: %s", r.URL.Path, err)
			http.Error(w, "Это не Float", http.StatusBadRequest) //Вывод error-ки
			return
		}
		h.storage.UpdateCounter(nameMetric, count)
	default: //Если передано другое значение значение
		h.log.Infof("передан не 'gauge' и не 'counter' | url: %s", r.URL.Path)
		http.Error(w, "это не 'gauge' и не 'counter' запросы", http.StatusBadRequest) //Вывод error-ки
		return                                                                        //Выход из ServeHTTP
	}
	h.log.Debugf("обновили метрику %s, типа %s, с значением %s", typeMetric, nameMetric, countMetric)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK) //Отправляет ответ что всё ОК
}
