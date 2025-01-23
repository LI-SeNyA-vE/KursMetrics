package send

import (
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/agentcfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcagent/metrics/update"
	"time"
)

func SendingBatchMetric(gaugeMetrics map[string]float64, counterMetrics map[string]int64, cfg agentcfg.Agent) {
	ticker1 := time.NewTicker(time.Duration(cfg.FlagPollInterval) * time.Second)
	ticker2 := time.NewTicker(time.Duration(cfg.FlagReportInterval) * time.Second)
	defer ticker1.Stop()
	defer ticker2.Stop()

	for {
		select {
		case <-ticker1.C:
			gaugeMetrics, counterMetrics = update.UpdateMetric()
			fmt.Printf("Пауза в %d секунд между сборкой метрик\n", cfg.FlagPollInterval)
		case <-ticker2.C:
			SendBatchJSONMetricsGauge(gaugeMetrics, cfg.FlagAddressAndPort, cfg.FlagKey)
			SendBatchJSONMetricsCounter(counterMetrics, cfg.FlagAddressAndPort, cfg.FlagKey)
			fmt.Printf("Пауза в %d секунд между отправкой 'батчей' метрик на сервер\n", cfg.FlagReportInterval)
		}
	}
}
