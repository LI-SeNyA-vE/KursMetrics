package server

import (
	"fmt"
	"net/http"
	"strconv"

	storageMetric "github.com/LI-SeNyA-vE/YaGo/internal/incriment1/storage"
)

// Проверяет количство запросов в URL
func ValidationLengthsURL(parts []string, w http.ResponseWriter) bool {
	if len(parts) < 5 { //Проверка на то, что приходят все элементы //ВОЗМОЖНО НУЖНО УБРАТЬ ЧТО БЫ ОБРАБАТЫВАТЬ БОЛЬШЕ ОШИБОК
		http.Error(w, "Короткий запрос", http.StatusNotFound) //Вывод error-ки
		return true
	}
	return false
}

// Проверка на первый элемент URL
func ValidationFirstElementURL(update string, name string, w http.ResponseWriter) bool {
	if update != name { //Проверка что идёт запрос на обновление (update)
		http.Error(w, "No update", http.StatusBadRequest) //Вывод error-ки
		return true
	}
	return false
}

// Проверяет какую метрику передали
func ValidationTypeMetric(typeMetric string, nameMetric string, countMetric string, w http.ResponseWriter) bool {
	switch typeMetric { //Свитч для проверки что это запрос или gauge или counter
	case "gauge": //Если передано значение 'gauge'
		count, err := strconv.ParseFloat(countMetric, 64) //Проверка что переданно число и его можно перевети в float64
		if err != nil {                                   //Если не удалось перевести
			http.Error(w, "Это не Float", http.StatusBadRequest) //Вывод error-ки
			return true                                          //
		}
		storageMetric.Metric.UpdateGauge(nameMetric, count)
		fmt.Println("Это gauge запрос")
		fmt.Println(storageMetric.Metric.Gauge)
		return false
	case "counter": //Если передано значение 'counter'
		{
			count, err := strconv.ParseInt(countMetric, 10, 64) //Проверка что переданно число и его можно перевети в int64
			if err != nil {                                     //Если не удалось перевести
				http.Error(w, "Это не Float", http.StatusBadRequest) //Вывод error-ки
				return true                                          //
			}
			storageMetric.Metric.UpdateCounter(nameMetric, count)
			fmt.Println("Это counter запрос")
			fmt.Println(storageMetric.Metric.Counter)
			return false
		}
	default: //Если передано другое значение значение
		{
			http.Error(w, "Это не 'gauge' и не 'counter' запросы", http.StatusBadRequest) //Вывод error-ки
			return true                                                                   //Выход из ServeHTTP
		}
	}
}
