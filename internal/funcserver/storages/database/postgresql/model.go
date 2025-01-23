package postgresql

import (
	"database/sql"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/sirupsen/logrus"
)

type DataBase struct {
	log *logrus.Entry
	cfg servercfg.Server
	db  *sql.DB
}
