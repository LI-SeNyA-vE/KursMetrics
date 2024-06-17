package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"runtime"

	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
)

var (
	addressAndPort = flag.String("a", "localhost:8080", "Указываем адресс и порт по которому будем потключаться")
	reportInterval = flag.Int64("r", 10, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	pollInterval   = flag.Int64("p", 2, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
)

type Config struct {
	address        string `env:"ADDRESS"`
	reportInterval int64  `env:"REPORT_INTERVAL"`
	pollInterval   int64  `env:"POLL_INTERVAL"`
}

type CounterMetrics struct {
	PollCount int64
}

type GaugeMetrics struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
	RandomValue   float64
}

// Функия сборки и записи метрик
func collectGaugeMetrics(m *GaugeMetrics) {
	var memStats runtime.MemStats   //Присваиваешь к переменной структуру runtime.MemStats
	runtime.ReadMemStats(&memStats) // Функция которая заполняет структуру runtime.MemStats
	// Заполнение и перевод в нужный пит данных
	m.Alloc = float64(memStats.Alloc)
	m.BuckHashSys = float64(memStats.BuckHashSys)
	m.Frees = float64(memStats.Frees)
	m.GCCPUFraction = float64(memStats.GCCPUFraction)
	m.GCSys = float64(memStats.GCSys)
	m.HeapAlloc = float64(memStats.HeapAlloc)
	m.HeapIdle = float64(memStats.HeapIdle)
	m.HeapInuse = float64(memStats.HeapInuse)
	m.HeapObjects = float64(memStats.HeapObjects)
	m.HeapReleased = float64(memStats.HeapReleased)
	m.HeapSys = float64(memStats.HeapSys)
	m.LastGC = float64(memStats.LastGC)
	m.Lookups = float64(memStats.Lookups)
	m.MCacheInuse = float64(memStats.MCacheInuse)
	m.MCacheSys = float64(memStats.MCacheSys)
	m.MSpanInuse = float64(memStats.MSpanInuse)
	m.MSpanSys = float64(memStats.MSpanSys)
	m.Mallocs = float64(memStats.Mallocs)
	m.NextGC = float64(memStats.NextGC)
	m.NumForcedGC = float64(memStats.NumForcedGC)
	m.NumGC = float64(memStats.NumGC)
	m.OtherSys = float64(memStats.OtherSys)
	m.PauseTotalNs = float64(memStats.PauseTotalNs)
	m.StackInuse = float64(memStats.StackInuse)
	m.StackSys = float64(memStats.StackSys)
	m.Sys = float64(memStats.Sys)
	m.TotalAlloc = float64(memStats.TotalAlloc)
	m.RandomValue = rand.Float64()
}

func collectCounterMetrics(m *CounterMetrics) {
	m.PollCount += 1
}

func sendMetrics(g *GaugeMetrics, c *CounterMetrics) {
	client := resty.New()
	gaugeMetrics := map[string]float64{
		"Alloc":         g.Alloc,
		"BuckHashSys":   g.BuckHashSys,
		"Frees":         g.Frees,
		"GCCPUFraction": g.GCCPUFraction,
		"GCSys":         g.GCSys,
		"HeapAlloc":     g.HeapAlloc,
		"HeapIdle":      g.HeapIdle,
		"HeapInuse":     g.HeapInuse,
		"HeapObjects":   g.HeapObjects,
		"HeapReleased":  g.HeapReleased,
		"HeapSys":       g.HeapSys,
		"LastGC":        g.LastGC,
		"Lookups":       g.Lookups,
		"MCacheInuse":   g.MCacheInuse,
		"MCacheSys":     g.MCacheSys,
		"MSpanInuse":    g.MSpanInuse,
		"MSpanSys":      g.MSpanSys,
		"Mallocs":       g.Mallocs,
		"NextGC":        g.NextGC,
		"NumForcedGC":   g.NumForcedGC,
		"NumGC":         g.NumGC,
		"OtherSys":      g.OtherSys,
		"PauseTotalNs":  g.PauseTotalNs,
		"StackInuse":    g.StackInuse,
		"StackSys":      g.StackSys,
		"Sys":           g.Sys,
		"TotalAlloc":    g.TotalAlloc,
		"RandomValue":   g.RandomValue,
	}

	counterMetrics := map[string]int64{
		"PollCount": c.PollCount,
	}

	for nameMetric, value := range gaugeMetrics {
		url := fmt.Sprintf("http://%s/update/gauge/%s/%f", *addressAndPort, nameMetric, value)
		_, err := client.R().SetHeader("Content-Type", "text/plain").Post(url)
		if err != nil {
			fmt.Printf("Ошибка при создании запроса для метрики %s: %s", nameMetric, err)
			return
		}
	}

	for nameMetric, value := range counterMetrics {
		url := fmt.Sprintf("http://%s/update/counter/%s/%b", *addressAndPort, nameMetric, value)
		_, err := client.R().SetHeader("Content-Type", "text/plain").Post(url)
		if err != nil {
			fmt.Printf("Ошибка при создании запроса для метрики %s: %s", nameMetric, err)
			return
		}
	}
}

func main() {
	var cfg Config

	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.address != "" {
		addressAndPort = &cfg.address
	}
	if cfg.reportInterval != 0 {
		reportInterval = &cfg.reportInterval
	}
	if cfg.pollInterval != 0 {
		pollInterval = &cfg.pollInterval
	}

	gaugeMetrics := &GaugeMetrics{}
	counterMetrics := &CounterMetrics{}
	ticker1 := time.NewTicker(time.Duration(*pollInterval) * time.Second)
	ticker2 := time.NewTicker(time.Duration(*reportInterval) * time.Second)
	for {
		select {
		case <-ticker1.C:
			collectGaugeMetrics(gaugeMetrics)
			collectCounterMetrics(counterMetrics)
			fmt.Printf("Пауза в %d секунд между сборкой метрик\n", *pollInterval)
		case <-ticker2.C:
			sendMetrics(gaugeMetrics, counterMetrics)
			fmt.Printf("Пауза в %d секунд между отправкой метрик на сервер\n", *reportInterval)
		}

	}
}
