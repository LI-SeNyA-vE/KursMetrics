package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/aeskey"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/rsakey"
	"io"
	"net/http"
)

func (m *Middleware) RsaDecoder(h http.Handler) http.Handler {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "rsa-encrypted" || r.Header.Get("Content-Encoding") == "gzip, rsa-encrypted" {
			if r.Header.Get("X-Encrypted-Hash") == "true" {
				var jsonBody struct {
					AesKey           string `json:"AES-KEY_Encode-RSA"`
					AesNode          string `json:"Rand_valu_AES-GCM"`
					EncryptedMessage string `json:"Encrypted_Message"`
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
					return
				}

				err = json.Unmarshal(body, &jsonBody)
				if err != nil {
					http.Error(w, "ошибка анмаршла тела запроса в JSON", http.StatusInternalServerError)
					return
				}

				// Декодируем AES-ключ из Base64
				encryptedAESKey, err := base64.StdEncoding.DecodeString(jsonBody.AesKey)
				if err != nil {
					http.Error(w, "Ошибка декодирования AES-ключа", http.StatusInternalServerError)
					return
				}

				// Расшифровываем AES-ключ через RSA
				keyAES, err := rsakey.DecryptMessage(m.FlagCryptoKey, encryptedAESKey)
				if err != nil {
					http.Error(w, "Ошибка расшифровки AES-ключа", http.StatusInternalServerError)
					return
				}

				// Декодируем nonce из Base64
				nonce, err := base64.StdEncoding.DecodeString(jsonBody.AesNode)
				if err != nil {
					http.Error(w, "Ошибка декодирования nonce", http.StatusInternalServerError)
					return
				}

				// Декодируем зашифрованное сообщение из Base64
				encryptedMessage, err := base64.StdEncoding.DecodeString(jsonBody.EncryptedMessage)
				if err != nil {
					http.Error(w, "Ошибка декодирования зашифрованного сообщения", http.StatusInternalServerError)
					return
				}

				// Декодируем сообщение Encrypted_Message через AES
				message, err := aeskey.DecryptMessage(encryptedMessage, keyAES, nonce)
				if err != nil {
					http.Error(w, "ошибка расшифровки переданного сообщения зашифрованного AES", http.StatusInternalServerError)
					return
				}

				// Восстанавливаем r.Body после расшифровки
				r.Body = io.NopCloser(bytes.NewBuffer(message))
			}
		}
		// Передаём управление следующему хендлеру
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hashFn)
}
