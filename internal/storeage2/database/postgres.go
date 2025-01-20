package database_v2

import (
	"database/sql"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config"
	"github.com/sirupsen/logrus"
	"log"
)

type DataBase struct {
	log *logrus.Entry
	cfg config.Server
	db  *sql.DB
}

func NewConnectDB(log *logrus.Entry, cfg config.Server) (*DataBase, error) {
	sysDB, err := sql.Open("pgx", cfg.FlagDatabaseDsn)
	if err != nil {
		log.Printf("Ошибка подключения к системной базе данных: %v", err)
		return nil, err
	}
	defer sysDB.Close()

	// Проверяем существование базы metric
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'metric')`
	err = sysDB.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Printf("Ошибка проверки существования базы данных: %v", err)
		return nil, err
	}

	if !exists {
		// Создаём базу metric
		_, err = sysDB.Exec(`CREATE DATABASE metric`)
		if err != nil {
			log.Printf("Ошибка создания базы данных metric: %v", err)
			return nil, err
		}
		log.Println("База данных metric успешно создана.")
	}

	// Подключаемся к базе metric
	db, err := sql.Open("pgx", "host=localhost dbname=metric user=Senya password=1q2w3e4r5t sslmode=disable")
	if err != nil {
		log.Printf("Ошибка подключения к базе данных metric: %v", err)
		return nil, err
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Printf("Не удалось установить соединение с базой данных metric: %v", err)
		return nil, err
	}

	// Проверяем и создаём таблицы
	err = ensureTablesExist(db)
	if err != nil {
		log.Printf("Ошибка проверки/создания таблиц: %v", err)
		return nil, err
	}

	return &DataBase{
		log: log,
		cfg: cfg,
		db:  db,
	}, nil
}

func ensureTablesExist(db *sql.DB) error {
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

	log.Println("Таблицы gauges и counters проверены/созданы успешно.")
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
		d.log.Printf("Ошибка обновления gauge: %v", err)
	}

	//TODO сделать вывод новой переменной
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
		d.log.Printf("Ошибка обновления counter: %v", err)
	}

	//TODO сделать вывод новой переменной
	return value
}

func (d *DataBase) GetAllGauges() map[string]float64 {
	rows, err := d.db.Query(`SELECT name, value FROM gauges`)
	if err != nil {
		d.log.Printf("Ошибка получения gauges: %v", err)
		return nil
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			d.log.Printf("Ошибка чтения строки gauge: %v", err)
			continue
		}
		result[name] = value
	}
	return result
}

func (d *DataBase) GetAllCounters() map[string]int64 {
	rows, err := d.db.Query(`SELECT name, value FROM counters`)
	if err != nil {
		d.log.Printf("Ошибка получения counters: %v", err)
		return nil
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var name string
		var value int64
		if err := rows.Scan(&name, &value); err != nil {
			d.log.Printf("Ошибка чтения строки counter: %v", err)
			continue
		}
		result[name] = value
	}
	return result
}

func (d *DataBase) GetGauge(name string) (float64, error) {
	//TODO сделай меня
	return 0, nil
}

func (d *DataBase) GetCounter(name string) (int64, error) {
	//TODO сделай меня
	return 0, nil
}

func (d *DataBase) LoadMetric() (err error) {
	//TODO сделай меня
	return err
}
