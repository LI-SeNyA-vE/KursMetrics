// Package postgresql содержит реализацию интерфейса MetricsStorage на базе PostgreSQL.
// Модель структуры необходимая для правильной работы приложения
package postgresql

import (
	"database/sql"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/sirupsen/logrus"
)

/*
DataBase представляет собой структуру, хранящую:

  - log: логгер на базе logrus.Entry,
  - cfg: конфигурацию сервера (servercfg.Server),
  - db: объект *sql.DB (активное соединение с базой данных PostgreSQL).

Она реализует интерфейс MetricsStorage, обеспечивая методы
для чтения, записи и обновления метрик (counter, gauge).
*/
type DataBase struct {
	log *logrus.Entry
	cfg servercfg.Server
	db  *sql.DB
}
