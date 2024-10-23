package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	funcAgent "github.com/LI-SeNyA-vE/KursMetrics/internal/funcAgent"
)

func main() {
	//Инициализация конфига для Агента
	config.InitializeConfig()

	//Вытаскиваем/обновляем метрики
	gaugeMetrics, counterMetrics := funcAgent.UpdateMetric()

	//Запускает горутину на отправку файлов кажные N секунд
	go func() {
		//funcAgent.SendingMetric(gaugeMetrics, counterMetrics, config.ConfigFlags.FlagPollInterval, config.ConfigFlags.FlagReportInterval, config.ConfigFlags.FlagAddressAndPort)
		funcAgent.SendingBatchMetric(gaugeMetrics, counterMetrics, config.ConfigFlags.FlagPollInterval, config.ConfigFlags.FlagReportInterval, config.ConfigFlags.FlagAddressAndPort)
	}()
	select {}

}
