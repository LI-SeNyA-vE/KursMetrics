package send

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func sendMetrics(client *resty.Client, url string, compressedData []byte, flagKey string) (interface{}, error) {
	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(compressedData)

	if flagKey != "" {
		h := hmac.New(sha256.New, []byte(flagKey)) //Подменяю ключ для теста
		h.Write(compressedData)
		hash := hex.EncodeToString(h.Sum(nil))

		request.SetHeader("HashSHA256", hash)
	}
	response, err := request.Post(url)

	// Если произошла ошибка или статус-код не 2xx, возвращаем ошибку
	if err != nil || response.StatusCode() >= 400 {
		return nil, fmt.Errorf(string(response.Body()), " ||| Флаг на агенте ", flagKey)
	}

	return response, nil
}
