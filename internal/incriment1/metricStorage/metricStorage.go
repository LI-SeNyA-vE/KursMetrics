package matricStore

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Структура для хранения метрик в памяти
type MetricStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

// Конструктор для создания нового экземпляра Metrictorage
func NewMetricStorage() *MetricStorage {
	return &MetricStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

var Metric = NewMetricStorage()

// Обновление значения gauge метрики (Замена значения)
func (m *MetricStorage) UpdateGauge(name string, value float64) {
	m.gauge[name] = value
}

// Обновление значения counter метрики (суммирование значений)
func (m *MetricStorage) UpdateCounter(name string, value int64) {
	m.counter[name] += value
}

// Проверяет количство запросов в URL
func validationLengthsURL(parts []string, w http.ResponseWriter) bool {
	if len(parts) < 5 { //Проверка на то, что приходят все элементы //ВОЗМОЖНО НУЖНО УБРАТЬ ЧТО БЫ ОБРАБАТЫВАТЬ БОЛЬШЕ ОШИБОК
		http.Error(w, "Короткий запрос", http.StatusNotFound) //Вывод error-ки
		return true
	}
	return false
}

// Проверка на первый элемент URL
func validationFirstElementURL(update string, name string, w http.ResponseWriter) bool {
	if update != name { //Проверка что идёт запрос на обновление (update)
		http.Error(w, "No update", http.StatusBadRequest) //Вывод error-ки
		return true
	}
	return false
}

// Проверяет какую метрику передали
func validationTypeMetric(typeMetric string, nameMetric string, countMetric string, w http.ResponseWriter) bool {
	switch typeMetric { //Свитч для проверки что это запрос или gauge или counter
	case "gauge": //Если передано значение 'gauge'
		{
			count, err := strconv.ParseFloat(countMetric, 64) //Проверка что переданно число и его можно перевети в float64
			if err != nil {                                   //Если не удалось перевести
				http.Error(w, "Это не Float", http.StatusBadRequest) //Вывод error-ки
				return true                                          //
			}
			Metric.UpdateGauge(nameMetric, count)
			fmt.Println("Это gauge запрос")
			fmt.Println(Metric.gauge)
			return false
		}
	case "counter": //Если передано значение 'counter'
		{
			count, err := strconv.ParseInt(countMetric, 10, 64) //Проверка что переданно число и его можно перевети в int64
			if err != nil {                                     //Если не удалось перевести
				http.Error(w, "Это не Float", http.StatusBadRequest) //Вывод error-ки
				return true                                          //
			}
			Metric.UpdateCounter(nameMetric, count)
			fmt.Println("Это counter запрос")
			fmt.Println(Metric.counter)
			return false
		}
	default: //Если передано другое значение значение
		{
			http.Error(w, "Это не 'gauge' и не 'counter' запросы", http.StatusBadRequest) //Вывод error-ки
			return true                                                                   //Выход из ServeHTTP
		}
	}
}

// Обработка HTTP запросов
func (m *MetricStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL: ", r.URL.Path)        //Можно удалить
	parts := strings.Split(r.URL.Path, "/") //Разделение приходящего URL по знаку '/'
	if validationLengthsURL(parts, w) {
		return
	}
	update := parts[1]      //Создание переменных для удобного использования в коде
	typeMetric := parts[2]  //Создание переменных для удобного использования в коде
	nameMetric := parts[3]  //Создание переменных для удобного использования в коде
	countMetric := parts[4] //Создание переменных для удобного использования в коде
	if validationFirstElementURL(update, "update", w) {
		return
	}
	if validationTypeMetric(typeMetric, nameMetric, countMetric, w) {
		return
	}
	w.WriteHeader(http.StatusOK) //Отправляет ответ что всё ОК
}
