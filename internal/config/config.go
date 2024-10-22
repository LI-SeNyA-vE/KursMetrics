package config

import (
	"flag"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	"reflect"
)

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

// ConfigFlags глобальная переменная, содержащая все флаги
var ConfigFlags VarFlag

func InitializeConfig() {
	//Запускает улучшенный логер
	err := logger.Initialize("info")
	if err != nil {
		panic("Не удалось инициализировать логгер")
	}

	//Парсит флаги
	ConfigFlags = NewVarFlag()

	//Парсит переменные окружения
	var cfgEnv VarEnv
	err = env.Parse(&cfgEnv)
	if err != nil {
		logger.Log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}

	//Проверяет если переменные окружения не пустые, то берёт их за основные (в флаг присваивает значение перем. окруж.)
	parseAllEnv(cfgEnv, &ConfigFlags)
}

// NewVarFlag инициализирует структуру VarFlag и парсит флаги командной строки
func NewVarFlag() VarFlag {
	cfgFlags := &VarFlag{}

	// Определение флагов
	flag.StringVar(&cfgFlags.FlagAddressAndPort, "a", "localhost:8080", "Указываем адрес и порт по которому будем подключаться")
	flag.Int64Var(&cfgFlags.FlagReportInterval, "r", 10, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	flag.Int64Var(&cfgFlags.FlagPollInterval, "p", 2, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	flag.StringVar(&cfgFlags.FlagLogLevel, "l", "info", "Уровень логирования")
	flag.Int64Var(&cfgFlags.FlagStoreInterval, "i", 30, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&cfgFlags.FlagFileStoragePath, "f", "C:\\Users\\Сеня\\Desktop\\KursMetrics\\cmd\\server\\metrics-db.json", "Полное имя файла, куда сохраняются текущие значения")
	flag.BoolVar(&cfgFlags.FlagRestore, "b", true, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&cfgFlags.FlagDatabaseDsn, "d", "host=localhost dbname=postgres user=Senya password=1q2w3e4r5t sslmode=disable", "Строка подключения к базе данных")

	// Парсинг флагов
	flag.Parse()

	return *cfgFlags
}

// parseAllEnv ты была рождена, что бы уменишить другую функцию
// parseAllEnv обновляет флаги на основе переменных окружения
func parseAllEnv(cfgEnv VarEnv, cfgFlags *VarFlag) {
	checkForNil(cfgEnv.EnvAddress, &cfgFlags.FlagAddressAndPort)
	checkForNil(cfgEnv.EnvReportInterval, &cfgFlags.FlagReportInterval)
	checkForNil(cfgEnv.EnvPollInterval, &cfgFlags.FlagPollInterval)
	checkForNil(cfgEnv.EnvLogLevel, &cfgFlags.FlagLogLevel)
	checkForNil(cfgEnv.EnvStoreInterval, &cfgFlags.FlagStoreInterval)
	checkForNil(cfgEnv.EnvFileStoragePath, &cfgFlags.FlagFileStoragePath)
	checkForNil(cfgEnv.EnvRestore, &cfgFlags.FlagRestore)
	checkForNil(cfgEnv.EnvDatabaseDsn, &cfgFlags.FlagDatabaseDsn)
}

// checkForNil проверяет значение и устанавливает его, если оно не нулевое
func checkForNil(env interface{}, flag interface{}) {
	envValue := reflect.ValueOf(env)
	flagValue := reflect.ValueOf(flag)

	// Проверяем, что flag является указателем
	if flagValue.Kind() != reflect.Ptr {
		logger.Log.Info("flag должен быть указателем")
	}

	// Получаем реальное значение переменной флага
	flagElem := flagValue.Elem()

	// Проверяем типы и обновляем флаг, если env не пустое
	if !envValue.IsZero() {
		if envValue.Type().AssignableTo(flagElem.Type()) {
			flagElem.Set(envValue)
		} else {
			logger.Log.Infof("Невозможно присвоить тип %s к %s\n", envValue.Type(), flagElem.Type())
		}
	}
}
