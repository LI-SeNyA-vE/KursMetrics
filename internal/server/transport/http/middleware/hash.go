/*
Package middleware содержит цепочку промежуточных обработчиков (middleware),
которые перехватывают и обрабатывают запросы и/или ответы сервера.
HashSHA256 проверяет целостность тела запроса на основе заголовка "HashSHA256".
*/
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

// HashSHA256 проверяет, содержит ли запрос заголовок HashSHA256. Если да,
// читает всё тело запроса в память и вычисляет для него HMAC SHA256 с помощью
// ключа m.FlagKey. Сравнивает результат с пришедшим в заголовке. В случае
// несоответствия выдаётся ошибка 400 (Bad Request), а при успехе –
// запрос передаётся дальше, причём r.Body восстанавливается для последующих обработчиков.
func (m *Middleware) HashSHA256(h http.Handler) http.Handler {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		receivedHash := r.Header.Get("HashSHA256")
		if receivedHash != "" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
				return
			}

			// Вычисляем хеш HMAC SHA256
			hasher := hmac.New(sha256.New, []byte(m.FlagKey))
			hasher.Write(body)
			calculatedHash := hex.EncodeToString(hasher.Sum(nil))

			// Сравниваем полученный хеш с переданным клиентом
			if receivedHash != calculatedHash {
				http.Error(w, fmt.Sprint("неверный хеш | Флаг на сервере ", m.FlagKey), http.StatusBadRequest)
				return
			}

			// Восстанавливаем r.Body после чтения
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		// Передаём управление следующему хендлеру
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hashFn)
}
