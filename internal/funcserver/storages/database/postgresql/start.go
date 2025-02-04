/*
Package postgresql предоставляет реализацию интерфейса MetricsStorage
(см. internal/funcserver/storages/metric.go) с использованием PostgreSQL
в качестве основного хранилища метрик.
*/
package postgresql

import (
	"database/sql"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/sirupsen/logrus"
)

// NewConnectDB осуществляет подключение к базе данных PostgreSQL.
// Если база "metrics" не существует, она создаётся. После чего происходит проверка/создание
// необходимых таблиц gauges и counters (ensureTablesExist). Возвращает объект DataBase
// или ошибку, если возникли проблемы с подключением или созданием структуры.
func NewConnectDB(log *logrus.Entry, cfg servercfg.Server) (*DataBase, error) {
	sysDB, err := sql.Open("pgx", cfg.FlagDatabaseDsn)
	if err != nil {
		log.Printf("ошибка подключения к системной базе данных: %v", err)
		return nil, err
	}
	defer sysDB.Close()

	// Проверяем существование базы metrics
	var exists bool
	query := queryExistDatname
	err = sysDB.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Printf("ошибка проверки существования базы данных: %v", err)
		return nil, err
	}

	// Если базы metrics нет — создаём
	if !exists {
		_, err = sysDB.Exec(queryCreateDatMetric)
		if err != nil {
			log.Printf("ошибка создания базы данных metrics: %v", err)
			return nil, err
		}
		log.Println("база данных metrics успешно создана.")
	}

	// Подключаемся к базе metrics
	db, err := sql.Open("pgx", cfg.FlagDatabaseDsn)
	if err != nil {
		log.Printf("ошибка подключения к базе данных metrics: %v", err)
		return nil, err
	}

	// Проверка соединения (ping)
	if err := db.Ping(); err != nil {
		log.Printf("не удалось установить соединение с базой данных metrics: %v", err)
		return nil, err
	}

	// Проверяем и создаём таблицы (gauges, counters)
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
