// Package ipandcidr предназначен для работы с IP адресами.
// Узнать локальны и глобальный IP.
// Сравнить IP с доверенной зоной
package ipandcidr

import (
	"net"
)

// GetLocalIP Функция для получения локального IP-адреса
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Ошибка: " + err.Error()
	}

	for _, addr := range addrs {
		// Проверяем, является ли адрес IPv4 и не является ли loopback (127.0.0.1)
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "Не удалось определить локальный IP"
}
