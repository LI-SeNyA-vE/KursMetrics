package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcagent/services"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
)

func main() {
	//Инициализация логера
	log := logger.NewLogger()

	//Инициализация конфига для Агента
	cfgAgent := config.NewConfigAgent(log)
	cfgAgent.InitializeAgentConfig()
	//Вытаскиваем/обновляем метрики
	gaugeMetrics, counterMetrics := services.UpdateMetric()

	//Запускает горутину на отправку файлов каждые N секунд
	go func() {
		//funcagent.SendingMetric(gaugeMetrics, counterMetrics, config.ConfigServerFlags.FlagPollInterval, config.ConfigServerFlags.FlagReportInterval, config.ConfigServerFlags.FlagAddressAndPort)
		services.SendingBatchMetric(gaugeMetrics, counterMetrics, cfgAgent.Agent)
	}()
	select {}

}
