// Package filemetric предоставляет реализацию хранилища метрик в локальном хранилище.
// Выгружает данные из файла для начало работы
package filemetric

import (
	"encoding/json"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"os"
)

// NewFileStorage — конструктор FileStorage
func NewFileStorage(cfg servercfg.Server) (*FileStorage, error) {
	storage := &FileStorage{
		cfg: cfg,
		data: struct {
			Gauges   map[string]float64 `json:"gauges"`
			Counters map[string]int64   `json:"counters"`
		}{
			Gauges:   make(map[string]float64),
			Counters: make(map[string]int64),
		},
	}

	// Загружаем данные из файла, если он существует
	if cfg.FlagRestore {
		if _, err := os.Stat(cfg.FlagFileStoragePath); err == nil {
			file, err := os.Open(cfg.FlagFileStoragePath)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			if err := json.NewDecoder(file).Decode(&storage.data); err != nil {
				return nil, err
			}
		}
	}

	return storage, nil
}
