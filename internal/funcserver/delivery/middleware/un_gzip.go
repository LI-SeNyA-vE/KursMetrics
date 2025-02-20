/*
Package middleware предоставляет набор функций промежуточных обработчиков (middleware),
используемых в проекте KursMetrics для модификации и анализа входящих/исходящих данных.

GzipMiddleware, вопреки названию, проверяет заголовок Accept-Encoding
и, если он равен "gzip", сжимает (а не распаковывает) исходящий ответ,
настраивая соответствующий заголовок "Content-Encoding: gzip".
*/
package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
)

// gzipWriter обёртывает http.ResponseWriter, заменяя операцию записи
// на запись в Writer, который осуществляет gzip-сжатие.
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write переопределяет метод записи для gzipWriter, направляя байты
// в gzip-сжатый поток (w.Writer), вместо исходного ResponseWriter.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipMiddleware проверяет заголовок "Accept-Encoding". Если он равен "gzip",
// оборачивает http.ResponseWriter в gzip.Writer, таким образом сжимая весь
// последующий контент. Если заголовок не равен "gzip", просто вызывает
// следующий хендлер без изменений.
//
// Параметры:
//   - next: следующий http.Handler в цепочке.
//
// Поведение:
//   - При "Accept-Encoding: gzip": устанавливает "Content-Encoding: gzip"
//     и оборачивает ответ в gzipWriter с уровнем BestSpeed.
//   - Иначе: выполняет next без изменений.
func (m *Middleware) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept-Encoding") != "gzip" {
			m.log.Info("Accept-Encoding не равен gzip")
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			m.log.Info("Ошибка при создании gzip.Writer с уровнем BestSpeed")
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		// Устанавливаем заголовок для ответа
		w.Header().Set("Content-Encoding", "gzip")

		// Передаём управление следующему хендлеру,
		// оборачивая ResponseWriter в gzipWriter
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
