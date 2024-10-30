package middleware

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"io"
	"net/http"
	"time"
)

type (
	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}

	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		uri    string
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func GzipMiddleware(next http.Handler) http.Handler {
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

func UnGzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info("Провалились в функцию UnGzipMiddleware")
		if !(r.Header.Get("Accept-Encoding") == "gzip") {
			logger.Log.Info("Accept-Encoding не равен gzip")
			next.ServeHTTP(w, r)
			return
		}
		logger.Log.Info("Accept-Encoding равен gzip")
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			logger.Log.Info("Ошибка при gzip.NewWriterLevel(w, gzip.BestSpeed)")
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()
		logger.Log.Info("Нет ошибки при gzip.NewWriterLevel(w, gzip.BestSpeed)")
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

	})
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggingMiddleware(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // функция Now() возвращает текущее время
		responseData := &responseData{
			status: 0,
			uri:    r.RequestURI, // Захват URI из запроса
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)           // обслуживание оригинального запроса
		duration := time.Since(start) // Since возвращает разницу во времени между start

		logger.Log.Infoln(
			"uri:", r.RequestURI,
			"method:", r.Method,
			"status:", responseData.status, // получаем перехваченный код статуса ответа
			"duration:", duration,
			"size:", len(responseData.uri), // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn) // возвращаем функционально расширенный хендлер
}

func HashSHA256(next http.Handler, flagKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHash := r.Header.Get("HashSHA256")
		if receivedHash != "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
				return
			}

			// Вычисляем хеш от тела запроса с использованием ключа
			h := hmac.New(sha256.New, []byte(flagKey))
			h.Write(body)

			calculatedHash := hex.EncodeToString(h.Sum(nil))

			if r.Header.Get("HashSHA256") != calculatedHash {
				http.Error(w, "Неверный хеш", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))

		}
		next.ServeHTTP(w, r)
	})
}
