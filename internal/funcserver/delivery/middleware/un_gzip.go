package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
)

func (m *Middleware) UnGzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(r.Header.Get("Accept-Encoding") == "gzip") {
			m.log.Info("Accept-Encoding не равен gzip")
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			m.log.Info("Ошибка при gzip.NewWriterLevel(w, gzip.BestSpeed)")
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

	})
}
