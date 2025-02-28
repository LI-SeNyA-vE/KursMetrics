/*
Package handlers содержит HTTP-обработчики (Handler),
которые выполняют операции обновления или чтения метрик.
Функция PostAddValue обрабатывает POST-запросы на обновление метрики по URL:

	/update/{typeMetric}/{nameMetric}/{countMetric}

где:
  - typeMetric  — тип метрики (gauge или counter),
  - nameMetric  — имя метрики (строка),
  - countMetric — значение метрики (число, интерпретируемое в зависимости от типа).

Если значение не удаётся преобразовать к необходимому типу
(для gauge — float64, для counter — int64),
возвращается 400 (Bad Request).
*/
package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// PostAddValue обрабатывает URL-запрос /update/{typeMetric}/{nameMetric}/{countMetric},
// определяя тип метрики (gauge или counter) и её значение (float64 или int64).
// В случае успеха вызывает UpdateGauge или UpdateCounter в storage. Если парсинг
// числа завершился ошибкой или передан неподдерживаемый тип, клиенту возвращается
// статус 400 (Bad Request).
func (h *Handler) PostAddValue(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")
	countMetric := chi.URLParam(r, "countMetric")

	switch typeMetric {
	case "gauge":
		// Проверяем, что countMetric можно интерпретировать как float64
		count, err := strconv.ParseFloat(countMetric, 64)
		if err != nil {
			h.log.Infof("передано не число или его нельзя перевести в float64 | url: %v | оригинальная ошибка: %s", r.URL.Path, err)
			http.Error(w, "это не Float", http.StatusBadRequest)
			return
		}
		h.storage.UpdateGauge(nameMetric, count)

	case "counter":
		// Проверяем, что countMetric можно интерпретировать как int64
		count, err := strconv.ParseInt(countMetric, 10, 64)
		if err != nil {
			h.log.Infof("передано не число или его нельзя перевести в int64 | url: %v | оригинальная ошибка: %s", r.URL.Path, err)
			http.Error(w, "Это не Float", http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounter(nameMetric, count)

	default:
		// Тип метрики не поддерживается
		h.log.Infof("передан не 'gauge' и не 'counter' | url: %s", r.URL.Path)
		http.Error(w, "это не 'gauge' и не 'counter' запросы", http.StatusBadRequest)
		return
	}

	// Логируем успешное обновление
	h.log.Debugf("обновили метрику %s, типа %s, со значением %s", typeMetric, nameMetric, countMetric)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
