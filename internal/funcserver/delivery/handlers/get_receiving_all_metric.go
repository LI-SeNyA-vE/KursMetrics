package handlers

import (
	"fmt"
	"net/http"
)

// GetReceivingAllMetric Возвращает страничку со всеми метриками
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
	listC := h.storage.GetAllCounters()
	for k, v := range listC {
		body = body + fmt.Sprintf("<tr>\n<td>%s</td>\n", k)
		body = body + fmt.Sprintf("<td>%v</td>\n</tr>\n", v)
	}

	listG := h.storage.GetAllGauges()
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
