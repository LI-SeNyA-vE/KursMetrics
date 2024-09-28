package main

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	funcAgent "github.com/LI-SeNyA-vE/KursMetrics/internal/funcAgent"
)

func main() {
	//Инициализация конфига для Агента
	config.InitializeConfigAgent()

	//Вытаскиваем/обновляем метрики
	gaugeMetrics, counterMetrics := funcAgent.UpdateMetric()

	//Запускает горутину на отправку файлов кажные N секунд
	go func() { funcAgent.SendingMetric(gaugeMetrics, counterMetrics) }()
	select {}

}
