package config

import (
	"flag"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	"reflect"
)

// VarServerFlag содержит все флаги как обычные поля
type VarServerFlag struct {
	FlagAddressAndPort  string
	FlagLogLevel        string
	FlagStoreInterval   int64
	FlagFileStoragePath string
	FlagRestore         bool
	FlagDatabaseDsn     string
}

type VarServerEnv struct {
	EnvAddress         string `env:"ADDRESS"`
	EnvLogLevel        string `env:"LOG_LEVEL"`
	EnvStoreInterval   int64  `env:"STORE_INTERVAL"`
	EnvFileStoragePath string `env:"FILE_STORAGE_PATH"`
	EnvRestore         bool   `env:"RESTORE"`
	EnvDatabaseDsn     string `env:"DATABASE_DSN"`
}

type VarAgentFlag struct {
	FlagAddressAndPort string
	FlagReportInterval int64
	FlagPollInterval   int64
	FlagLogLevel       string
}

type VarAgentEnv struct {
	EnvAddress        string `env:"ADDRESS"`
	EnvReportInterval int64  `env:"REPORT_INTERVAL"`
	EnvPollInterval   int64  `env:"POLL_INTERVAL"`
	EnvLogLevel       string `env:"LOG_LEVEL"`
}

// ConfigServerFlags глобальная переменная, содержащая все флаги
var ConfigServerFlags VarServerFlag
var ConfigAgentFlags VarAgentFlag

//func InitializeConfig() {
//	//Запускает улучшенный логер
//	err := logger.Initialize("info")
//	if err != nil {
//		panic("Не удалось инициализировать логгер")
//	}
//
//	//Парсит флаги
//	ConfigServerFlags = NewVarServerFlag()
//	ConfigAgentFlags = NewVarAgentFlag()
//
//	//Парсит переменные окружения для сервера
//	var cfgServerEnv VarServerEnv
//	err = env.Parse(&cfgServerEnv)
//	if err != nil {
//		logger.Log.Info("Ошибка на этапе парсинга переменных окружения", err)
//	}
//	//Проверяет если переменные окружения не пустые, то берёт их за основные (в флаг присваивает значение перем. окруж.)
//	parseServerEnv(cfgServerEnv, &ConfigServerFlags)
//
//	//Парсит переменные окружения для агента
//	var cfgAgentEnv VarAgentEnv
//	err = env.Parse(&cfgAgentEnv)
//	if err != nil {
//		logger.Log.Info("Ошибка на этапе парсинга переменных окружения", err)
//	}
//	//Проверяет если переменные окружения не пустые, то берёт их за основные (в флаг присваивает значение перем. окруж.)
//	parseAgentEnv(cfgAgentEnv, &ConfigAgentFlags)
//}

func InitializeServerConfig() {
	//Запускает улучшенный логер
	err := logger.Initialize("info")
	if err != nil {
		panic("Не удалось инициализировать логгер")
	}

	//Парсит флаги
	ConfigServerFlags = NewVarServerFlag()

	//Парсит переменные окружения для сервера
	var cfgServerEnv VarServerEnv
	err = env.Parse(&cfgServerEnv)
	if err != nil {
		logger.Log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}
	//Проверяет если переменные окружения не пустые, то берёт их за основные (в флаг присваивает значение перем. окруж.)
	parseServerEnv(cfgServerEnv, &ConfigServerFlags)
}

func InitializeAgentConfig() {
	//Запускает улучшенный логер
	err := logger.Initialize("info")
	if err != nil {
		panic("Не удалось инициализировать логгер")
	}

	//Парсит флаги
	ConfigAgentFlags = NewVarAgentFlag()

	//Парсит переменные окружения для агента
	var cfgAgentEnv VarAgentEnv
	err = env.Parse(&cfgAgentEnv)
	if err != nil {
		logger.Log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}
	//Проверяет если переменные окружения не пустые, то берёт их за основные (в флаг присваивает значение перем. окруж.)
	parseAgentEnv(cfgAgentEnv, &ConfigAgentFlags)
}

// NewVarServerFlag инициализирует структуру VarServerFlag и парсит флаги командной строки
func NewVarServerFlag() VarServerFlag {
	cfgFlags := &VarServerFlag{}

	// Определение флагов
	flag.StringVar(&cfgFlags.FlagAddressAndPort, "a", "localhost:8080", "Указываем адрес и порт по которому будем подключаться")
	flag.StringVar(&cfgFlags.FlagLogLevel, "l", "info", "Уровень логирования")
	flag.Int64Var(&cfgFlags.FlagStoreInterval, "i", 30, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&cfgFlags.FlagFileStoragePath, "f", "C:\\Users\\Сеня\\Desktop\\KursMetrics\\cmd\\server\\metrics-db.json", "Полное имя файла, куда сохраняются текущие значения")
	flag.BoolVar(&cfgFlags.FlagRestore, "r", true, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&cfgFlags.FlagDatabaseDsn, "d", "host=localhost dbname=postgres user=Senya password=1q2w3e4r5t sslmode=disable", "Строка подключения к базе данных")
	// Парсинг флагов
	flag.Parse()

	return *cfgFlags
}

// parseServerEnv ты была рождена, что бы уменишить другую функцию
// parseServerEnv обновляет флаги на основе переменных окружения
func parseServerEnv(cfgServerEnv VarServerEnv, cfgServerFlags *VarServerFlag) {
	checkForNil(cfgServerEnv.EnvAddress, &cfgServerFlags.FlagAddressAndPort)
	checkForNil(cfgServerEnv.EnvLogLevel, &cfgServerFlags.FlagLogLevel)
	checkForNil(cfgServerEnv.EnvStoreInterval, &cfgServerFlags.FlagStoreInterval)
	checkForNil(cfgServerEnv.EnvFileStoragePath, &cfgServerFlags.FlagFileStoragePath)
	checkForNil(cfgServerEnv.EnvRestore, &cfgServerFlags.FlagRestore)
	checkForNil(cfgServerEnv.EnvDatabaseDsn, &cfgServerFlags.FlagDatabaseDsn)
}

func NewVarAgentFlag() VarAgentFlag {
	cfgAgentFlags := &VarAgentFlag{}

	// Определение флагов
	flag.StringVar(&cfgAgentFlags.FlagAddressAndPort, "a", "localhost:8080", "Указываем адрес и порт по которому будем подключаться")
	flag.Int64Var(&cfgAgentFlags.FlagReportInterval, "r", 10, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	flag.Int64Var(&cfgAgentFlags.FlagPollInterval, "p", 2, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	flag.StringVar(&cfgAgentFlags.FlagLogLevel, "l", "info", "Уровень логирования")
	// Парсинг флагов
	flag.Parse()

	return *cfgAgentFlags
}
func parseAgentEnv(cfgAgentEnv VarAgentEnv, cfgAgentFlags *VarAgentFlag) {
	checkForNil(cfgAgentEnv.EnvAddress, &cfgAgentFlags.FlagAddressAndPort)
	checkForNil(cfgAgentEnv.EnvReportInterval, &cfgAgentFlags.FlagReportInterval)
	checkForNil(cfgAgentEnv.EnvPollInterval, &cfgAgentFlags.FlagPollInterval)
	checkForNil(cfgAgentEnv.EnvLogLevel, &cfgAgentFlags.FlagLogLevel)

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
