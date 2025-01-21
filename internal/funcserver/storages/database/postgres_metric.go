package database

import (
	"database/sql"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/sirupsen/logrus"
)

type DataBase struct {
	log *logrus.Entry
	cfg config.Server
	db  *sql.DB
}

func NewConnectDB(log *logrus.Entry, cfg config.Server) (*DataBase, error) {
	sysDB, err := sql.Open("pgx", cfg.FlagDatabaseDsn)
	if err != nil {
		log.Printf("ошибка подключения к системной базе данных: %v", err)
		return nil, err
	}
	defer sysDB.Close()

	// Проверяем существование базы metric
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'metric')`
	err = sysDB.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Printf("ошибка проверки существования базы данных: %v", err)
		return nil, err
	}

	if !exists {
		// Создаём базу metric
		_, err = sysDB.Exec(`CREATE DATABASE metric`)
		if err != nil {
			log.Printf("ошибка создания базы данных metric: %v", err)
			return nil, err
		}
		log.Println("база данных metric успешно создана.")
	}

	// Подключаемся к базе metric
	db, err := sql.Open("pgx", cfg.FlagDatabaseDsn)
	if err != nil {
		log.Printf("ошибка подключения к базе данных metric: %v", err)
		return nil, err
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Printf("не удалось установить соединение с базой данных metric: %v", err)
		return nil, err
	}

	// Проверяем и создаём таблицы
	err = ensureTablesExist(db, log)
	if err != nil {
		log.Printf("ошибка проверки/создания таблиц: %v", err)
		return nil, err
	}

	return &DataBase{
		log: log,
		cfg: cfg,
		db:  db,
	}, nil
}

func ensureTablesExist(db *sql.DB, log *logrus.Entry) error {
	// Проверка существования таблицы gauges
	gaugesQuery := `
		CREATE TABLE IF NOT EXISTS gauges (
			name TEXT PRIMARY KEY,
			value DOUBLE PRECISION NOT NULL
		)`
	_, err := db.Exec(gaugesQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы gauges: %w", err)
	}

	// Проверка существования таблицы counters
	countersQuery := `
		CREATE TABLE IF NOT EXISTS counters (
			name TEXT PRIMARY KEY,
			value BIGINT NOT NULL
		)`
	_, err = db.Exec(countersQuery)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы counters: %w", err)
	}

	log.Println("таблицы gauges и counters проверены/созданы успешно.")
	return nil
}

func (d *DataBase) UpdateGauge(name string, value float64) float64 {
	_, err := d.db.Exec(`
		INSERT INTO gauges (name, value) 
		VALUES ($1, $2)
		ON CONFLICT (name) 
		DO UPDATE SET value = EXCLUDED.value
	`, name, value)
	if err != nil {
		d.log.Printf("ошибка обновления gauge: %v", err)
	}
	return value
}

func (d *DataBase) UpdateCounter(name string, value int64) int64 {
	_, err := d.db.Exec(`
		INSERT INTO counters (name, value) 
		VALUES ($1, $2)
		ON CONFLICT (name) 
		DO UPDATE SET value = counters.value + EXCLUDED.value
	`, name, value)
	if err != nil {
		d.log.Printf("ошибка обновления counter: %v", err)
	}
	return value
}

func (d *DataBase) GetAllGauges() map[string]float64 {
	rows, err := d.db.Query(`SELECT name, value FROM gauges`)

	if err != nil {
		d.log.Printf("ошибка получения gauges: %v", err)
		return nil
	}

	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			d.log.Printf("ошибка чтения строки gauge: %v", err)
			continue
		}
		result[name] = value
	}

	// Проверяем rows.Err после итерации
	if err = rows.Err(); err != nil {
		d.log.Printf("Ошибка итерации gauges: %v", err)
		return nil
	}

	return result
}

func (d *DataBase) GetAllCounters() map[string]int64 {
	rows, err := d.db.Query(`SELECT name, value FROM counters`)
	if err != nil {
		d.log.Printf("ошибка получения counters: %v", err)
		return nil
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var name string
		var value int64
		if err := rows.Scan(&name, &value); err != nil {
			d.log.Printf("ошибка чтения строки counter: %v", err)
			continue
		}
		result[name] = value
	}

	// Проверяем rows.Err после итерации
	if err = rows.Err(); err != nil {
		d.log.Printf("Ошибка итерации gauges: %v", err)
		return nil
	}

	return result
}

// GetGauge возвращает значение метрики типа gauge
func (d *DataBase) GetGauge(name string) (*float64, error) {
	var value float64
	query := `SELECT value FROM gauges WHERE name = $1`

	err := d.db.QueryRow(query, name).Scan(&value)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("gauge %q not found", name)
	} else if err != nil {
		return nil, fmt.Errorf("failed to query gauge %q: %w", name, err)
	}

	return &value, nil
}

// GetCounter возвращает значение метрики типа counter
func (d *DataBase) GetCounter(name string) (*int64, error) {
	var value int64
	query := `SELECT value FROM counters WHERE name = $1`

	err := d.db.QueryRow(query, name).Scan(&value)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("counter %q not found", name)
	} else if err != nil {
		return nil, fmt.Errorf("failed to query counter %q: %w", name, err)
	}

	return &value, nil
}

func (d *DataBase) LoadMetric() (err error) {
	return err
}
