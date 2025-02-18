// Package rsakey содержит функции для шифровки/расшифровки запроса.
// Шифрует/расшифровывает только необходимые данные для AES шифрования
package rsakey

import (
	"os"
)

// Функция загрузки ключа
func loadKey(path string) ([]byte, error) {
	return os.ReadFile(path)
}
