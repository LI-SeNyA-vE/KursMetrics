package funcagent

import (
	"context"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/agentcfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcagent/metrics/send"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcagent/metrics/update"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	_ "net/http/pprof"
	"sync"
	"time"
)

type MetricData struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

func Run() {
	var metricData MetricData
	metricData.gaugeMetrics = make(map[string]float64)
	metricData.counterMetrics = make(map[string]int64)
	var mu sync.Mutex
	//Инициализация логера
	log := logger.NewLogger()

	//Инициализация конфига для Агента
	cfgAgent := agentcfg.NewConfigAgent(log)
	cfgAgent.InitializeAgentConfig()

	// Вывод догов на дебаге, для отслеживания
	log.Debugf("Адрес сервера: %s | Интервал отправки: %d | Интервал опроса: %d | Уровень логирования: %s | Key: %s | RateLimit: %d",
		cfgAgent.FlagAddressAndPort,
		cfgAgent.FlagReportInterval,
		cfgAgent.FlagPollInterval,
		cfgAgent.FlagLogLevel,
		cfgAgent.FlagKey,
		cfgAgent.FlagRateLimit,
	)

	// Создадим общий контекст (на случай, если захотите останавливать горутины по cancel)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для «сырых» данных (jobs), которые нужно отправлять
	jobs := make(chan MetricData) // Можно засунуть в метод startWorkerPool

	// Запускаем Worker Pool с ограничением concurrency = FlagRateLimit
	// которая будет отправлять метрики
	var wg sync.WaitGroup
	startWorkerPool(ctx, &wg, jobs, *cfgAgent, log)

	// Запускаем горутину опроса runtime
	// runtime
	go func() {
		ticker := time.NewTicker(time.Duration(cfgAgent.FlagPollInterval) * time.Second)
		defer ticker.Stop()

		// Локальные переменные под данные метрик
		var runtimeGauge map[string]float64
		var runtimeCounter map[string]int64

		for {
			select {
			case <-ticker.C:
				// Собираем runtime-метрики
				runtimeGauge, runtimeCounter = update.UpdateMetric()

				// Блокируем переменную для записи метрик
				mu.Lock()
				// Записываем в переменную значения полученные с UpdateMetric
				for name, value := range runtimeGauge {
					metricData.gaugeMetrics[name] = value
				}
				// Записываем в переменную значения полученные с UpdateMetric
				for name, value := range runtimeCounter {
					metricData.counterMetrics[name] = value
				}
				// Снимаем блокировку
				mu.Unlock()
				log.Info("Собрали метрики runtime")
			case <-ctx.Done():
				return
			}
		}
	}()

	// Запускаем горутину опроса системных метрик через gopsutil
	// system
	go func() {
		ticker := time.NewTicker(time.Duration(cfgAgent.FlagPollInterval) * time.Second)
		defer ticker.Stop()

		gaugeMetrics := make(map[string]float64)

		for {
			select {
			case <-ticker.C:
				// Пример использования gopsutil
				if vmStat, err := mem.VirtualMemory(); err == nil {
					gaugeMetrics["TotalMemory"] = float64(vmStat.Total)
					gaugeMetrics["FreeMemory"] = float64(vmStat.Free)
				}
				if cpuPercent, err := cpu.Percent(0, false); err == nil {
					gaugeMetrics["CPUutilization1"] = cpuPercent[0]
				}

				// Блокируем переменную для записи метрик
				mu.Lock()
				// Записываем в переменную значения полученные с VirtualMemory и Percent
				for name, value := range gaugeMetrics {
					metricData.gaugeMetrics[name] = value
				}
				// Снимаем блокировку
				mu.Unlock()
				log.Info("Собрали метрики gopsutil")
			case <-ctx.Done():
				return
			}
		}
	}()

	// Запускаем горутину, которая периодически будет «формировать» финальный объект MetricData
	// и складывать его в канал jobs, чтобы воркеры могли отправлять его на сервер.
	go func() {
		ticker := time.NewTicker(time.Duration(cfgAgent.FlagReportInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// Сформировали структуру и отправили «задание» на отправку
				mu.Lock()
				gauge := copyGauge(metricData.gaugeMetrics)
				counter := copyCounter(metricData.counterMetrics)
				mu.Unlock()
				jobs <- MetricData{
					gaugeMetrics:   gauge,
					counterMetrics: counter,
				}

				fmt.Printf("Послали job на отправку. Gauge=%d, Counter=%d, запросов в очереди: %d\n",
					len(gauge),
					len(counter),
					len(jobs),
				)

				// После отправки можно обнулить/очистить локальные данные,
				// чтобы «начинать с чистого листа» (если требуется).
				metricData.gaugeMetrics = make(map[string]float64)
				metricData.counterMetrics = make(map[string]int64)
			case <-ctx.Done():
				return
			}
		}
	}()

	// «Зависаем» в select{}, чтобы main не завершался
	select {}
}

// startWorkerPool запускает нужное количество воркеров (равное rateLimit).
// Каждый воркер читает из канала jobs структуру MetricData и отправляет метрики на сервер.
func startWorkerPool(
	ctx context.Context,
	wg *sync.WaitGroup,
	jobs <-chan MetricData,
	cfg agentcfg.ConfigAgent,
	log interface{ Infof(string, ...interface{}) },
) {
	var (
		rateLimit  = cfg.FlagRateLimit
		serverAddr = cfg.FlagAddressAndPort
		key        = cfg.FlagKey
	)
	if rateLimit < 1 {
		rateLimit = 1
	}
	for i := int64(0); i < rateLimit; i++ {
		wg.Add(1)
		go func(workerID int64) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case mData, ok := <-jobs:
					if !ok {
						return
					}
					fmt.Printf("Запросов в очереди: %d", len(jobs))
					send.SendBatchJSONMetrics(mData.gaugeMetrics, mData.counterMetrics, serverAddr, key)
				}
			}
		}(i)
	}
}

// sendMetrics — пример функции отправки набора метрик.
// Вы можете разделять на gauge/counter, как у вас было в send.SendBatchJSONMetricsGauge/Counter.
// Здесь для наглядности общее.

// Функции для копирования map, чтобы не было гонок,
// если мы хотим «снимок» данных перед отправкой:
func copyGauge(src map[string]float64) map[string]float64 {
	dst := make(map[string]float64, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
func copyCounter(src map[string]int64) map[string]int64 {
	dst := make(map[string]int64, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
