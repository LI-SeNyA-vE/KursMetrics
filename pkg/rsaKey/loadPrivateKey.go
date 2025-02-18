package rsaKey

import (
	"os"
)

// Функция загрузки ключа
func loadKey(path string) ([]byte, error) {
	return os.ReadFile(path)
}
