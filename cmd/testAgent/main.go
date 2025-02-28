package main

import "github.com/LI-SeNyA-vE/KursMetrics/internal/agent/metrics/send"

var (
	mapMetricGauge     map[string]float64 = map[string]float64{}
	mapMetricCounter   map[string]int64   = map[string]int64{}
	flagAddressAndPort string             = "localhost:8080"
	flagHashKey        string             = ""
	flagRsaKey         string             = "rsaKeys/publicKey.pem"
)

func main() {
	mapMetricGauge["test_gauge"] = 1.4
	mapMetricCounter["test_count"] = 1
	send.SendBatchJSONMetricsHTTP(mapMetricGauge, mapMetricCounter, flagAddressAndPort, flagHashKey, flagRsaKey)
}
