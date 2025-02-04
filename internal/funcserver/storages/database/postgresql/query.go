/*
Package postgresql предоставляет реализацию интерфейса MetricsStorage
(см. internal/funcserver/storages/metric.go) с использованием PostgreSQL
в качестве основного хранилища метрик.
*/
package postgresql

// Набор констант со SQL-запросами, необходимыми для проверки наличия базы,
// создания её при необходимости, а также создания и обновления таблиц
// counters и gauges в базе.
const (
	queryExistDatname      = `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'metrics')`
	queryCreateDatMetric   = `CREATE DATABASE metrics`
	queryGetCounter        = `SELECT value FROM counters WHERE name = $1`
	queryGetGauge          = `SELECT value FROM gauges WHERE name = $1`
	queryGetAllGauges      = `SELECT name, value FROM gauges`
	queryGetAllCounters    = `SELECT name, value FROM counters`
	queryExistCounterTable = `
		CREATE TABLE IF NOT EXISTS counters (
			name TEXT PRIMARY KEY,
			value BIGINT NOT NULL
		)`
	queryExistGaugesTable = `
		CREATE TABLE IF NOT EXISTS gauges (
			name TEXT PRIMARY KEY,
			value DOUBLE PRECISION NOT NULL
		)`
	queryUpdateGauge = `
		INSERT INTO gauges (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE SET value = EXCLUDED.value
	`
	queryUpdateCounter = `
		INSERT INTO counters (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE SET value = counters.value + EXCLUDED.value
	`
)
