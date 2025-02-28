//Package send содержит функции для отправки метрик с Агента на Сервер,
//включая сериализацию в JSON, сжатие (gzip) и валидацию через HMAC.
//Функция gzipCompress выполняет сжатие переданных данных для уменьшения размера
//перед отправкой по сети.

package send

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

// gzipCompress архивирует (gzip-сжимает) переданные данные data,
// возвращая сжатый байтовый срез. Если во время записи или закрытия
// gzip.Writer возникает ошибка, она будет обёрнута и возвращена.
func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("ошибка записи данных в gzip writer: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("ошибка закрытия gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}
