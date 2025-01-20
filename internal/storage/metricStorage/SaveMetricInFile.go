package storage

//func SaveMetricToFile(cdFile string) {
//	allMetrics := MetricStorage{
//		gauge:   Metric.GetAllGauges(),
//		counter: Metric.GetAllCounters(),
//	}
//
//	data, err := json.Marshal(allMetrics)
//	if err != nil {
//		log.Print(err)
//	}
//
//	dir := filepath.Dir(cdFile)
//	os.MkdirAll(dir, 0755)
//	os.WriteFile(cdFile, data, 0666)
//	if err != nil {
//		log.Print(err)
//		return
//	}
//}
