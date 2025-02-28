//Package handlers содержит набор HTTP-обработчиков (Handler),
//которые отвечают за приём, обновление и вывод метрик.

package handlers

import (
	"fmt"
	"net/http"
)

// GetReceivingAllMetric формирует HTML-страницу со всеми имеющимися в хранилище
// метриками (как counter, так и gauge) и возвращает её клиенту.
// Для каждой метрики в таблице отображаются имя и текущее значение.
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
	// Получаем все counter-метрики
	listC := h.storage.GetAllCounters()
	for k, v := range listC {
		body += fmt.Sprintf("<tr>\n<td>%s</td>\n", k)
		body += fmt.Sprintf("<td>%v</td>\n</tr>\n", v)
	}

	// Получаем все gauge-метрики
	listG := h.storage.GetAllGauges()
	for k, v := range listG {
		body += fmt.Sprintf("<tr>\n<td>%s</td>\n", k)
		body += fmt.Sprintf("<td>%v</td>\n</tr>\n", v)
	}

	body += " </table>\n </body>\n</html>"

	// Устанавливаем заголовки и отправляем готовую HTML-страницу
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
