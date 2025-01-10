package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	funcAgent "github.com/LI-SeNyA-vE/KursMetrics/internal/funcAgent"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
)

func main() {
	//Инициализация логера
	log := logger.NewLogger()

	//Инициализация конфига для Агента
	cfgAgent := config.NewConfigAgent(log)
	cfgAgent.InitializeAgentConfig()
	//Вытаскиваем/обновляем метрики
	gaugeMetrics, counterMetrics := funcAgent.UpdateMetric()

	//Запускает горутину на отправку файлов каждые N секунд
	go func() {
		//funcAgent.SendingMetric(gaugeMetrics, counterMetrics, config.ConfigServerFlags.FlagPollInterval, config.ConfigServerFlags.FlagReportInterval, config.ConfigServerFlags.FlagAddressAndPort)
		funcAgent.SendingBatchMetric(gaugeMetrics, counterMetrics, cfgAgent.Agent)
	}()
	select {}

}
