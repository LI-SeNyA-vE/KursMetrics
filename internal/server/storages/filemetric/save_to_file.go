// Package filemetric предоставляет реализацию хранилища метрик в локальном хранилище.
// saveToFile для сохранения в файл изменения метрик.
package filemetric

import (
	"encoding/json"
	"os"
)

func (s *FileStorage) saveToFile() {
	file, err := os.Create(s.cfg.FlagFileStoragePath)
	if err != nil {
		// Обработка ошибки сохранения
		return
	}
	defer file.Close()
	_ = json.NewEncoder(file).Encode(s.data)
}
