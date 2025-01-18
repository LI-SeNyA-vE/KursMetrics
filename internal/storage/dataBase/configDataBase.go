package dataBase

import (
	"database/sql"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	metricStorage "github.com/LI-SeNyA-vE/KursMetrics/internal/storage/metricStorage"
	"github.com/sirupsen/logrus"
	"log"
)

type DataBase struct {
	log *logrus.Entry
	config.Server
	*sql.DB
}

func NewConnectDB(log *logrus.Entry, cfg config.Server) *DataBase {
	return &DataBase{
		log:    log,
		Server: cfg,
		DB:     nil,
	}
}

func (d *DataBase) CrereateDB(createTableSQL string) {
	_, err := d.DB.Exec(createTableSQL)
	if err != nil {
		d.log.Infof("Ошибка при создании таблицы: %v", err)

		return
	}
	d.log.Info("Таблица 'Metric' успешно создана или уже была.")
}

// ConnectDB функция для подключения к БД
func (d *DataBase) ConnectDB() (err error) {
	d.DB, err = sql.Open("pgx", d.FlagDatabaseDsn)
	if err != nil {
		d.log.Infof("Ошибка подключения к базе данных: %v", err)
		return err
	}

	err = d.DB.Ping()
	if err != nil {
		d.log.Info("Не удалось установить соединение с базой данных: %v", err)
		return err
	}

	configCreateSQL := CreateConfigSQL()
	d.CrereateDB(configCreateSQL)

	return nil
}

func (d *DataBase) LoadMetricFromDB() (err error) {
	rows, err := d.DB.Query("SELECT Id, Type, Name, Value FROM metric")
	if err != nil {
		d.log.Info("Ошибка получения данных из базы данных: %v", err)
		return err
	} else {
		for rows.Next() {
			metric := &metricStorage.MetricStorage{}
			var idMetric string
			var typeMetric string
			var nameMetric string
			var valueMetric float64
			//err := rows.Scan(&metric.Id, &metric.Type, &metric.Name, &metric.Value)
			err = rows.Scan(idMetric, typeMetric, nameMetric, valueMetric)
			if err != nil {
				d.log.Info("Ошибка сканирования строки: %v", err)
			}
			defer rows.Close()

			switch typeMetric { //Свитч для проверки, что это запрос или gauge, или counter
			case "gauge": //Если передано значение 'gauge'
				metric.UpdateGauge(nameMetric, valueMetric)
			case "counter": //Если передано значение 'counter'
				metric.UpdateCounter(nameMetric, int64(valueMetric))
			default: //Если передано другое значение
				log.Println("При вытягивание данных из БД оказалось что тип не gauge и не counter")
			}
		}
		return err
	}
}

func CreateConfigSQL() string {
	var createTableSQL = `
  CREATE TABLE IF NOT EXISTS metric (
      "id" TEXT NOT NULL,
      "type" TEXT NOT NULL,
      "name" TEXT NOT NULL,
      "value" DOUBLE PRECISION NULL,
      PRIMARY KEY ("id", "type")
  );`
	return createTableSQL
}
