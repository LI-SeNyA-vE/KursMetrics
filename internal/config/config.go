package config

import (
	"flag"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
)

/*var (
	FlagAddressAndPort  = flag.String("a", "localhost:8080", "Указываем адресс и порт по которому будем потключаться")
	FlagRreportInterval = flag.Int64("r", 10, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	FlagPollInterval    = flag.Int64("p", 2, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	FlagLogLevel        = flag.String("l", "info", "Уровень логирования")
	FlagStoreInterval   = flag.Int64("i", 30, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	FlagFileStoragePath = flag.String("f", "/tmp/metrics-db.json", "Полное имя файла, куда сохраняются текущие значения")
	FlagRestore         = flag.Bool("b", true, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	FlagDatabaseDsn     = flag.String("d", "host=localhost dbname=postgres user=Senya password=1q2w3e4r5t sslmode=disable", "Определяет загружать ранее сохранённые значения из базы при старте сервера")
)*/

// VarFlag содержит все флаги как обычные поля
type VarFlag struct {
	FlagAddressAndPort  string
	FlagReportInterval  int64
	FlagPollInterval    int64
	FlagLogLevel        string
	FlagStoreInterval   int64
	FlagFileStoragePath string
	FlagRestore         bool
	FlagDatabaseDsn     string
}

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

// InitializeConfigAgent Конфиг для  Агента
func InitializeConfigAgent() {
	//Запускает улучшенный логер
	err := logger.Initialize("info")
	if err != nil {
		panic("Не удалось инициализировать логгер")
	}

	//Парсит флаги
	flag.Parse()

	//Парсит переменные окружения
	var cfg VarEnv
	err = env.Parse(&cfg)
	if err != nil {
		logger.Log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}
}

// InitializeConfigServer Конфиг для Сервера
func InitializeConfigServer() (cfgFlags *VarFlag) {
	//Запускает улучшенный логер
	err := logger.Initialize("info")
	if err != nil {
		panic("Не удалось инициализировать логгер")
	}

	//Парсит флаги
	cfgFlags = NewVarFlag()

	//Парсит переменные окружения
	var cfgEnv VarEnv
	err = env.Parse(&cfgEnv)
	if err != nil {
		logger.Log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}

	//Проверяет если переменные окружения не пустые, то берёт их за основные (в флаг присваивает значение перем. окруж.)
	parseAllEnv(cfgEnv, *cfgFlags)
	return cfgFlags
}

// NewVarFlag инициализирует структуру VarFlag и парсит флаги командной строки
func NewVarFlag() *VarFlag {
	cfgFlags := &VarFlag{}

	// Определение флагов
	flag.StringVar(&cfgFlags.FlagAddressAndPort, "a", "localhost:8080", "Указываем адрес и порт по которому будем подключаться")
	flag.Int64Var(&cfgFlags.FlagReportInterval, "r", 10, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	flag.Int64Var(&cfgFlags.FlagPollInterval, "p", 2, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	flag.StringVar(&cfgFlags.FlagLogLevel, "l", "info", "Уровень логирования")
	flag.Int64Var(&cfgFlags.FlagStoreInterval, "i", 30, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&cfgFlags.FlagFileStoragePath, "f", "/tmp/metrics-db.json", "Полное имя файла, куда сохраняются текущие значения")
	flag.BoolVar(&cfgFlags.FlagRestore, "b", true, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&cfgFlags.FlagDatabaseDsn, "d", "host=localhost dbname=postgres user=Senya password=1q2w3e4r5t sslmode=disable", "Строка подключения к базе данных")

	// Парсинг флагов
	flag.Parse()

	return cfgFlags
}

// parseAllEnv ты была рождена, что бы уменишьт другую функцию
func parseAllEnv(cfgEnv VarEnv, cfgFlags VarFlag) {
	checkForNil(cfgEnv.EnvAddress, cfgFlags.FlagAddressAndPort)
	checkForNil(cfgEnv.EnvReportInterval, cfgFlags.FlagReportInterval)
	checkForNil(cfgEnv.EnvPollInterval, cfgFlags.FlagPollInterval)
	checkForNil(cfgEnv.EnvLogLevel, cfgFlags.FlagLogLevel)
	checkForNil(cfgEnv.EnvStoreInterval, cfgFlags.FlagStoreInterval)
	checkForNil(cfgEnv.EnvFileStoragePath, cfgFlags.FlagFileStoragePath)
	checkForNil(cfgEnv.EnvRestore, cfgFlags.FlagRestore)
	checkForNil(cfgEnv.EnvDatabaseDsn, cfgFlags.FlagDatabaseDsn)
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
