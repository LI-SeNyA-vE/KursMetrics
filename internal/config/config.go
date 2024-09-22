package config

import (
	"database/sql"
	"flag"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/handlers/middleware/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	FlagAddressAndPort  = flag.String("a", "localhost:8080", "Указываем адресс и порт по которому будем потключаться")
	FlagRreportInterval = flag.Int64("r", 10, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	FlagPollInterval    = flag.Int64("p", 2, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	FlagLogLevel        = flag.String("l", "info", "Уровень логирования")
	FlagStoreInterval   = flag.Int64("i", 30, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	FlagFileStoragePath = flag.String("f", "/tmp/metrics-db.json", "Полное имя файла, куда сохраняются текущие значения")
	FlagRestore         = flag.Bool("b", true, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	FlagDatabaseDsn     = flag.String("d", "host=localhost dbname=postgres user=Senya password=1q2w3e4r5t sslmode=disable", "Определяет загружать ранее сохранённые значения из базы при старте сервера")
)

// C:\GO\KursMetrics\cmd\server\metrics-db.json
type VarEnv struct {
	EnvAddress         string `env:"ADDRESS"`
	EnvReportInterval  int64  `env:"REPORT_INTERVAL"`
	EnvPollInterval    int64  `env:"POLL_INTERVAL"`
	EnvLogLevel        string `env:"LOG_LEVEL"`
	EnvStoreInterval   int64  `env:"STORE_INTERVAL"`
	EnvFileStoragePath string `env:"FILE_STORAGE_PATH"`
	EnvRestore         bool   `env:"RESTORE"`
	EnvDatabaseDsn     string `env:"DATABASE_DSN"`
}

type ConnectSQL struct {
	user      string
	password  string
	dbname    string
	sslmode   string
	tableName string
}

func ConfigSQL() string {
	var createTableSQL = `
  CREATE TABLE IF NOT EXISTS metric (
      "id" TEXT NOT NULL,
      "type" TEXT NOT NULL,
      "value" DOUBLE PRECISION NULL,
      PRIMARY KEY ("id", "type")
  );`
	return createTableSQL
}

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", *FlagDatabaseDsn)
	logger.Log.Infoln("Ссылка на подключение: %s", *FlagDatabaseDsn)
	if err != nil {
		logger.Log.Infoln("Ошибка подключения к базе данных: %v", err)
		return db, err
	}

	err = db.Ping()
	if err != nil {
		logger.Log.Infoln("Не удалось установить соединение с базой данных: %v", err)
		return db, err
	}

	return db, nil
}

// InitializeGlobals инициализирует флаги на основе значений из конфигурации
func InitializeGlobals() {
	flag.Parse()
	var cfg VarEnv
	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	checkForNil(cfg.EnvAddress, FlagAddressAndPort)
	checkForNil(cfg.EnvReportInterval, FlagRreportInterval)
	checkForNil(cfg.EnvPollInterval, FlagPollInterval)
	checkForNil(cfg.EnvLogLevel, FlagLogLevel)
	checkForNil(cfg.EnvStoreInterval, FlagStoreInterval)
	checkForNil(cfg.EnvFileStoragePath, FlagFileStoragePath)
	checkForNil(cfg.EnvRestore, FlagRestore)
	checkForNil(cfg.EnvDatabaseDsn, FlagDatabaseDsn)
}

// checkForNil проверяет значение и устанавливает его, если оно не нулевое
func checkForNil(enc interface{}, flag interface{}) {
	switch enc := enc.(type) {
	case string:
		if enc != "" {
			*flag.(*string) = enc
		}
	case int64:
		if enc != 0 {
			*flag.(*int64) = enc
		}
	case bool:
		if enc {
			*flag.(*bool) = enc
		}
	}
}
