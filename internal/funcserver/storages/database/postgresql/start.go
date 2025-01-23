package postgresql

import (
	"database/sql"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/sirupsen/logrus"
)

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

	if !exists {
		// Создаём базу metrics
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

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Printf("не удалось установить соединение с базой данных metrics: %v", err)
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
