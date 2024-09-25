package config

import (
	"database/sql"
	"flag"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/storage/dataBase"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
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

//type VarFlag struct {
//	FlagAddressAndPort  flag.Flag
//	FlagRreportInterval flag.Flag
//	FlagPollInterval    flag.Flag
//	FlagLogLevel        flag.Flag
//	FlagStoreInterval   flag.Flag
//	FlagFileStoragePath flag.Flag
//	FlagRestore         flag.Flag
//	FlagDatabaseDsn     flag.Flag
//}

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

// ConnectDB функция для проверки подключения к БД
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

// InitializeConfigAgent Конфиг для  Агента
func InitializeConfigAgent() {
	//Запускает улучшенный логер
	err := logger.Initialize("info")
	if err != nil {
		panic("Не удалось инициализировать логгер")
	}

	//Парсит флаги
	flag.Parse()
	var cfg VarEnv

	//Парсит переменные окружения
	err = env.Parse(&cfg)
	if err != nil {
		logger.Log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}
}

// InitializeConfigServer Конфиг для Сервера
func InitializeConfigServer() {
	//Запускает улучшенный логер
	err := logger.Initialize("info")
	if err != nil {
		panic("Не удалось инициализировать логгер")
	}

	//Парсит флаги
	flag.Parse()
	var cfg VarEnv

	//Парсит переменные окружения
	err = env.Parse(&cfg)
	if err != nil {
		logger.Log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}

	//Проверяет если переменные окружения не пустые, то берёт их за основные (в флаг присваивает значение перем. окруж.)
	parseAllEnv(cfg)

	//Запускается функция, которая определит куда сохранять данные
	dataBase.InitializeStorage(*FlagFileStoragePath, *FlagRestore, *FlagDatabaseDsn)
}

// parseAllEnv ты была рождена, что бы уменишьт другую функцию
func parseAllEnv(cfg VarEnv) {
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
func checkForNil(env interface{}, flag interface{}) {
	switch enc := env.(type) {
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
