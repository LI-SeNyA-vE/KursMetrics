package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

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
	client := &http.Client{}
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
		url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%f", nameMetric, value)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Printf("Ошибка при создании запроса для метрики %s: %s", nameMetric, err)
			return
		}
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Ошибка при отправке метрики %s: %s", nameMetric, err)
			return
		}
		defer resp.Body.Close()
	}

	for nameMetric, value := range counterMetrics {
		url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%v", nameMetric, value)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Printf("Ошибка при создании запроса для метрики %s: %s", nameMetric, err)
			return
		}
		req.Header.Set("Content-Type", "text/plain")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Ошибка при отправке метрики %s: %s", nameMetric, err)
			return
		}
		defer resp.Body.Close()
	}
}

func main() {
	gaugeMetrics := &GaugeMetrics{}
	counterMetrics := &CounterMetrics{}
	for {
		collectGaugeMetrics(gaugeMetrics)
		collectCounterMetrics(counterMetrics)
		sendMetrics(gaugeMetrics, counterMetrics)

		fmt.Println("Пауза 5 сек")
		time.Sleep(5 * time.Second)
	}
}
