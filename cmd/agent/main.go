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

	log.Infof(" FlagAddressAndPort = %s | FlagReportInterval = %s | FlagPollInterval = %s| FlagLogLevel = %s | FlagKey = %s", cfgAgent.Agent.FlagAddressAndPort, cfgAgent.Agent.FlagReportInterval, cfgAgent.Agent.FlagPollInterval, cfgAgent.Agent.FlagLogLevel, cfgAgent.Agent.FlagKey)
	//Вытаскиваем/обновляем метрики
	gaugeMetrics, counterMetrics := services.UpdateMetric()

	//Запускает горутину на отправку файлов каждые N секунд
	go func() {
		services.SendingBatchMetric(gaugeMetrics, counterMetrics, cfgAgent.Agent)
	}()
	select {}

}
