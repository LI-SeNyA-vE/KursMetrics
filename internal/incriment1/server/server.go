package server

import (
	"fmt"
	"net/http"
	"strconv"

	storageMetric "github.com/LI-SeNyA-vE/KursMetrics/internal/incriment1/storage"
)

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
