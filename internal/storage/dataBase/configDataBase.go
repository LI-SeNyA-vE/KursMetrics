package dataBase

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers/middleware/logger"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"log"
)

func CreateConfigSQL() string {
	var createTableSQL = `
  CREATE TABLE IF NOT EXISTS metric (
      "id" TEXT NOT NULL,
      "type" TEXT NOT NULL,
      "value" DOUBLE PRECISION NULL,
      PRIMARY KEY ("id", "type")
  );`
	return createTableSQL
}

func InitializeStorage(cdFile string, resMetricBool bool, loadDataBase string) {
	if loadDataBase != "" {
		db, err := config.ConnectDB()
		if err != nil {
			logger.Log.Infoln("Ошибка связанная с ДБ: %v", err)
		}
		defer db.Close()

		configCreateSQL := CreateConfigSQL()
		metricStorage.CrereateDB(db, configCreateSQL)

		rows, err := db.Query("SELECT Id, Type, Name, Value FROM metric")
		if err != nil {
			logger.Log.Infoln("Ошибка получения данных из базы данных: %v", err)
		} else {
			for rows.Next() {
				metric := &metricStorage.MetricStorage{}
				var idMetric string
				var typeMetric string
				var nameMetric string
				var valueMetric float64
				//err := rows.Scan(&metric.Id, &metric.Type, &metric.Name, &metric.Value)
				err := rows.Scan(idMetric, typeMetric, nameMetric, valueMetric)
				if err != nil {
					logger.Log.Infoln("Ошибка сканирования строки: %v", err)
				}
				defer rows.Close()

				switch typeMetric { //Свитч для проверки что это запрос или gauge или counter
				case "gauge": //Если передано значение 'gauge'
					metric.UpdateGauge(nameMetric, valueMetric)
				case "counter": //Если передано значение 'counter'
					metric.UpdateCounter(nameMetric, int64(valueMetric))
				default: //Если передано другое значение значение
					log.Println("При вытягивание данных из БД оказалось что тип не gauge и не counter")
				}
			}
			return
		}
	}
	if resMetricBool {
		metricStorage.LoadMetricFromFile(cdFile)
	}
	return
}
