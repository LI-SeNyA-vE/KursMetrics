//Package send содержит логику отправки метрик с Агента на Сервер.
//Функция sendMetrics формирует HTTP-запрос к эндпоинту /updates/,
//устанавливает необходимые заголовки (включая HMAC-хеш для проверки целостности,
//если предоставлен ключ flagKey) и осуществляет POST-запрос.

package send

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// sendMetrics принимает на вход:
//   - REST client (*resty.Client),
//   - url для POST-запроса,
//   - сжатые данные compressedData,
//   - ключ flagKey (используется для HMAC SHA256, если не пуст).
//
// Формирует запрос, добавляет заголовки "Content-Encoding: gzip" / "Accept-Encoding: gzip"
// и при необходимости "HashSHA256" (если flagKey не пуст). Возвращает ответ сервера или ошибку,
// если произошёл сбой запроса либо статус-код >= 400.
func sendMetrics(client *resty.Client, url string, compressedData []byte, flagKey string) (interface{}, error) {
	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(compressedData)

	// Если указан ключ, рассчитываем HMAC SHA256 по сжатым данным и устанавливаем в заголовок
	if flagKey != "" {
		h := hmac.New(sha256.New, []byte(flagKey))
		h.Write(compressedData)
		hash := hex.EncodeToString(h.Sum(nil))
		request.SetHeader("HashSHA256", hash)
	}

	response, err := request.Post(url)

	// Если произошла ошибка запроса или сервер вернул код >= 400, создаём и возвращаем ошибку
	if err != nil || response.StatusCode() >= 400 {
		return nil, fmt.Errorf(string(response.Body()), " ||| Флаг на агенте ", flagKey)
	}

	return response, nil
}
