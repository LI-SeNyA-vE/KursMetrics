package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func (m *Middleware) HashSHA256(h http.Handler) http.Handler {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		receivedHash := r.Header.Get("HashSHA256")
		if receivedHash != "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
				return
			}

			// Вычисляем хеш от тела запроса с использованием ключа
			h := hmac.New(sha256.New, []byte(m.FlagKey)) //Подменяю ключ для теста
			h.Write(body)

			calculatedHash := hex.EncodeToString(h.Sum(nil))

			if r.Header.Get("HashSHA256") != calculatedHash {
				http.Error(w, fmt.Sprint("неверный хеш | Флаг на сервере ", m.FlagKey), http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))

		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hashFn)
}
