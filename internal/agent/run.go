// Package agent реализует логику работы Агента (Agent) в проекте KursMetrics.
// В этом пакете происходит сбор метрик (через runtime и gopsutil) и регулярная
// отправка их на сервер по заданному адресу. Работа организована посредством
// нескольких горутин и ограниченного пула воркеров (Worker Pool).
package agent

import (
	"context"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/rsakey"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/LI-SeNyA-vE/KursMetrics/internal/agent/metrics/send"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/agent/metrics/update"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/agentcfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	_ "net/http/pprof" // Импортируем для возможности запуска pprof
	"sync"
	"time"
)

// MetricData хранит набор метрик двух типов:
//   - gauge (группа метрик float64),
//   - counter (группа метрик int64).
type MetricData struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

// Run инициализирует работу всего агента, включая:
//
//   - Чтение конфигурации (флаги и переменные окружения),
//   - Логирование,
//   - Горутины по сбору метрик (runtime, gopsutil),
//   - Работу пула воркеров (Worker Pool) для отправки собранных метрик,
//
// и блокируется до завершения, чтобы приложение не выходило из main.
// Для остановки может быть использован контекст (ctx) с cancel().
func Run() {
	var metricData MetricData
	var mu sync.Mutex
	var wg sync.WaitGroup
	metricData.gaugeMetrics = make(map[string]float64)
	metricData.counterMetrics = make(map[string]int64)

	// Инициализация логгера.
	log := logger.NewLogger()

	// Инициализация конфига для Агента.
	cfgAgent := agentcfg.NewConfigAgent(log)
	cfgAgent.InitializeAgentConfig()

	err := rsakey.CheckKey(cfgAgent.FlagCryptoKey)
	if err != nil {
		//TODO сделать горутинку, которая будет проверять правильность открытого ключа и если он не правильный, то кидать запросы на сервере на отправку открытого ключа и не выполнять никаких других действий пока не получит ключ
		log.Errorf("не найден ключ публичный ключ: %v", err)
	}

	// Создаём общий контекст, который будет использован
	// для управляемого завершения горутин (через cancel()).
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// jobs — канал, куда будут поступать "задания" (снимки метрик) для отправки.
	jobs := make(chan MetricData, 1)

	// Запускаем пул воркеров, который будет отправлять метрики (SendBatchJSONMetricsHTTP).
	startWorkerPool(ctx, &wg, jobs, *cfgAgent, log)

	// Горутина опроса runtime-метрик каждые FlagPollInterval секунд.
	go collectRuntimeMetrics(ctx, &metricData, &mu, cfgAgent, log)

	// Горутина опроса системных метрик (через gopsutil) каждые FlagPollInterval секунд.
	go collectSystemMetrics(ctx, &metricData, &mu, cfgAgent, log)

	// Горутина, формирующая каждые FlagReportInterval секунд "снимок" метрик
	// и отправляющая его в канал jobs для асинхронной отправки.
	go reportMetrics(ctx, &metricData, &mu, jobs, cfgAgent)

	// Горутина проверяющая на вызов сигнала
	go handleSignals(cancel, &wg, jobs, &metricData, &mu)

	wg.Wait()
}

// startWorkerPool запускает несколько горутин-воркеров, количество которых
// определяется параметром rateLimit (если < 1, то ставится 1).
// Каждый воркер берёт MetricData из канала jobs и вызывает SendBatchJSONMetricsHTTP.
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
		hashKey    = cfg.FlagKey
		keyRsa     = cfg.FlagCryptoKey
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
					// Для отладки выводим размер очереди.
					fmt.Printf("Запросов в очереди: %d", len(jobs))
					//send.SendBatchJSONMetricsHTTP() или SendBatchJSONMetricsRPC()
					send.SendBatchJSONMetricsRPC(mData.gaugeMetrics, mData.counterMetrics, serverAddr, hashKey, keyRsa)
				}
			}
		}(i)
	}
}

func collectRuntimeMetrics(ctx context.Context, metricData *MetricData, mu *sync.Mutex, cfgAgent *agentcfg.ConfigAgent, log *logrus.Entry) {
	ticker := time.NewTicker(time.Duration(cfgAgent.FlagPollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runtimeGauge, runtimeCounter := update.UpdateMetric()
			mu.Lock()
			for name, value := range runtimeGauge {
				metricData.gaugeMetrics[name] = value
			}
			for name, value := range runtimeCounter {
				metricData.counterMetrics[name] = value
			}
			mu.Unlock()
			log.Info("Собрали метрики runtime")
		case <-ctx.Done():
			return
		}
	}
}

func collectSystemMetrics(ctx context.Context, metricData *MetricData, mu *sync.Mutex, cfgAgent *agentcfg.ConfigAgent, log *logrus.Entry) {
	ticker := time.NewTicker(time.Duration(cfgAgent.FlagPollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			gaugeMetrics := make(map[string]float64)
			if vmStat, err := mem.VirtualMemory(); err == nil {
				gaugeMetrics["TotalMemory"] = float64(vmStat.Total)
				gaugeMetrics["FreeMemory"] = float64(vmStat.Free)
			}
			if cpuPercent, err := cpu.Percent(0, false); err == nil {
				gaugeMetrics["CPUutilization1"] = cpuPercent[0]
			}
			mu.Lock()
			for name, value := range gaugeMetrics {
				metricData.gaugeMetrics[name] = value
			}
			mu.Unlock()
			log.Info("Собрали метрики gopsutil")
		case <-ctx.Done():
			return
		}
	}
}

func reportMetrics(ctx context.Context, metricData *MetricData, mu *sync.Mutex, jobs chan<- MetricData, cfgAgent *agentcfg.ConfigAgent) {
	ticker := time.NewTicker(time.Duration(cfgAgent.FlagReportInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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

			// Обнуление локального накопления, если нужно отслеживать "прирост".
			metricData.gaugeMetrics = make(map[string]float64)
			metricData.counterMetrics = make(map[string]int64)
		case <-ctx.Done():
			return
		}
	}
}

func handleSignals(cancel context.CancelFunc, wg *sync.WaitGroup, jobs chan MetricData, metricData *MetricData, mu *sync.Mutex) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-signalChan
	fmt.Println("Агент: получен сигнал завершения, передаем оставшиеся данные...")
	flushPendingMetrics(jobs, metricData, mu)
	cancel()
	close(jobs)
	wg.Wait()
	fmt.Println("Агент: завершение работы успешно.")
}

func flushPendingMetrics(jobs chan MetricData, metricData *MetricData, mu *sync.Mutex) {
	mu.Lock()
	deferredMetrics := MetricData{
		gaugeMetrics:   copyGauge(metricData.gaugeMetrics),
		counterMetrics: copyCounter(metricData.counterMetrics),
	}
	mu.Unlock()
	jobs <- deferredMetrics
	fmt.Println("Агент: успешно передал несохранённые данные перед завершением.")
}

// copyGauge создаёт копию переданной карты gauge-метрик, чтобы
// избежать гонок при конкурентном доступе.
func copyGauge(src map[string]float64) map[string]float64 {
	dst := make(map[string]float64, len(src))
	for k, v := range src {
		dst[k] = v
	}

	return dst
}

// copyCounter создаёт копию карты counter-метрик, чтобы
// сделать "снимок" данных перед отправкой.
func copyCounter(src map[string]int64) map[string]int64 {
	dst := make(map[string]int64, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
