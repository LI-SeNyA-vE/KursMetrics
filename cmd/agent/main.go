package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	funcAgent "github.com/LI-SeNyA-vE/KursMetrics/internal/funcAgent"
)

var ConfigAgent = &config.ConfigAgentFlags

func main() {
	//Инициализация конфига для Агента
	config.InitializeAgentConfig()

	//Вытаскиваем/обновляем метрики
	gaugeMetrics, counterMetrics := funcAgent.UpdateMetric()

	//Запускает горутину на отправку файлов каждые N секунд
	go func() {
		//funcAgent.SendingMetric(gaugeMetrics, counterMetrics, config.ConfigServerFlags.FlagPollInterval, config.ConfigServerFlags.FlagReportInterval, config.ConfigServerFlags.FlagAddressAndPort)
		funcAgent.SendingBatchMetric(gaugeMetrics, counterMetrics, ConfigAgent.FlagPollInterval, ConfigAgent.FlagReportInterval, ConfigAgent.FlagAddressAndPort, ConfigAgent.FlagKey)
	}()
	select {}

}
