// Package ipandcidr предназначен для работы с IP адресами.
// Узнать локальны и глобальный IP.
// Сравнить IP с доверенной зоной
package ipandcidr

import (
	"fmt"
	"net"
)

// IsIPInCIDR Функция для проверки, входит ли IP в диапазон CIDR
func IsIPInCIDR(ipStr, cidr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("невалидный IP: %s", ipStr)
	}

	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, fmt.Errorf("невалидный CIDR: %s", cidr)
	}

	return subnet.Contains(ip), nil
}
