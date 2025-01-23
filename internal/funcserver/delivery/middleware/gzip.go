package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
)

type (
	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
)

func (m *Middleware) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Ошибка при создании gzip.Reader", http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			// Замена r.Body на распакованный stream
			r.Body = io.NopCloser(gz)
		}
		next.ServeHTTP(w, r)
	})
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}
