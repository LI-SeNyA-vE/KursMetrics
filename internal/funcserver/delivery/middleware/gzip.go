/*
Package middleware содержит набор промежуточных обработчиков (middleware),
которые перехватывают входящие HTTP-запросы и/или исходящие ответы,
позволяя модифицировать или анализировать их до или после передачи к конечному хендлеру.
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

// GzipMiddleware проверяет, содержит ли запрос заголовок "Content-Encoding: gzip".
// Если да, то r.Body декомпрессируется "на лету" с помощью gzip.NewReader,
// и распакованные данные подставляются обратно в r.Body.
// Затем управление передаётся следующему хендлеру в цепочке.
func (m *Middleware) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, запакован ли Body gzip'ом
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Ошибка при создании gzip.Reader", http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			// Подменяем r.Body на поток распакованных данных
			r.Body = io.NopCloser(gz)
		}
		// Вызываем следующий хендлер в цепочке
		next.ServeHTTP(w, r)
	})
}

// Write переопределяет метод записи для gzipWriter, направляя байты
// в gzip-сжатый поток (w.Writer), вместо исходного ResponseWriter.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
