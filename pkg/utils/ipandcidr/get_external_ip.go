// Package ipandcidr предназначен для работы с IP адресами.
// Узнать локальны и глобальный IP.
// Сравнить IP с доверенной зоной
package ipandcidr

import (
	"io"
	"net/http"
)

// GetExternalIP Функция для получения внешнего IP-адреса
func GetExternalIP() string {
	resp, err := http.Get("https://api64.ipify.org?format=text")
	if err != nil {
		return "Ошибка: " + err.Error()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Ошибка при чтении ответа"
	}

	return string(body)
}
