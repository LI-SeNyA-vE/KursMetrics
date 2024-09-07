package storage

import (
	"encoding/json"
	"log"
	"os"
)

func LoadMetricFromFile(fstg string) {

	res, err := os.ReadFile(fstg)
	if err != nil {
		log.Printf("Ошибка чтения файла: %s", err)
	}

	if err := json.Unmarshal(res, &StorageMetric); err != nil {
		log.Printf("Ошибка Unmarshal: %s", err)
	}
}
