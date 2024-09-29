package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"io"
	"net/http"
	"strconv"

	storageMetric "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/go-chi/chi/v5"
)

var cfgFlags = config.VarFlag{}

// PostAddValue Обрабатывает полный url запрос. Если всё правильно сохраняет метрику в память
func PostAddValue(w http.ResponseWriter, r *http.Request) {
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
		storageMetric.StorageMetric.UpdateGauge(nameMetric, count)
	case "counter": //Если передано значение 'counter'
		{
			count, err := strconv.ParseInt(countMetric, 10, 64) //Проверка что переданно число и его можно перевети в int64
			if err != nil {                                     //Если не удалось перевести
				http.Error(w, "Это не Float", http.StatusBadRequest) //Вывод error-ки
				return                                               //
			}
			storageMetric.StorageMetric.UpdateCounter(nameMetric, count)
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
func GetReceivingMetric(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")
	typeMetric := chi.URLParam(r, "typeMetric")
	logger.Log.Info("Запрос с " + nameMetric + " " + typeMetric)
	value, err := storageMetric.StorageMetric.GetValue(typeMetric, nameMetric)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Log.Info("ошибка в GetReceivingMetric")
		return
	}
	io.WriteString(w, fmt.Sprint(value))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

/* func GetReceivingAllMetric(w http.ResponseWriter, r *http.Request) {
	gauges := storageMetric.Metric.GetAllGauges()
	counters := storageMetric.Metric.GetAllCounters()

	data := storageMetric.MetricStorage{
		Gauge:   gauges,
		Counter: counters,
	}

	tmplPath := filepath.Join("..", "..", "internal", "templates", "index.html")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		logger.Log.Info("Ошибка 1 " + err.Error())
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		logger.Log.Info("Ошибка 2 " + err.Error())
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
} */

// GetReceivingAllMetric Возвращает страничку со всеми метриками
func GetReceivingAllMetric(w http.ResponseWriter, r *http.Request) {
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
	listC := storageMetric.StorageMetric.GetAllCounters()
	for k, v := range listC {
		body = body + fmt.Sprintf("<tr>\n<td>%s</td>\n", k)
		body = body + fmt.Sprintf("<td>%v</td>\n</tr>\n", v)
	}

	listG := storageMetric.StorageMetric.GetAllGauges()
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
func JSONValue(w http.ResponseWriter, r *http.Request) {
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

	metric, err := storageMetric.StorageMetric.GetValue(metrics.MType, metrics.ID) // Запрашивает метрику, по данным из JSON
	if err != nil {
		logger.Log.Info(err)
		http.Error(w, "не найдено", http.StatusNotFound)
		return
	}

	//Проверка на тип
	switch metrics.MType {
	case "counter":
		v := metric.(int64)
		metrics.Delta = &v
	case "gauge":
		v := metric.(float64)
		metrics.Value = &v
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
func JSONUpdate(w http.ResponseWriter, r *http.Request) {
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
	case "counter":
		storageMetric.StorageMetric.UpdateCounter(metrics.ID, *metrics.Delta) //Обновляет метрику
	case "gauge":
		storageMetric.StorageMetric.UpdateGauge(metrics.ID, *metrics.Value) //Обновляет метрику
	}

	metric, err := storageMetric.StorageMetric.GetValue(metrics.MType, metrics.ID) // Запрашивает метрику, по данным из JSON
	if err != nil {
		http.Error(w, "не найдено", http.StatusNotFound)
		return
	}

	//Проверка на тип
	switch metrics.MType {
	case "counter":
		v := metric.(int64)
		metrics.Delta = &v
	case "gauge":
		v := metric.(float64)
		metrics.Value = &v
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
func Ping(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("pgx", config.ConfigFlags.FlagDatabaseDsn)
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
