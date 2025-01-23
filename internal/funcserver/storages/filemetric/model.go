package filemetric

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"sync"
)

type FileStorage struct {
	cfg  servercfg.Server
	mu   sync.Mutex
	data struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]int64   `json:"counters"`
	}
}
