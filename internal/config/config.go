package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	FlagAddressAndPort  = flag.String("a", "localhost:8080", "Указываем адресс и порт по которому будем потключаться")
	FlagRreportInterval = flag.Int64("r", 10, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	FlagPollInterval    = flag.Int64("p", 2, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	FlagLogLevel        = flag.String("l", "info", "Уровень логирования")
	FlagStoreInterval   = flag.Int64("i", 0, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	FlagFileStoragePath = flag.String("f", "/tmp/metrics-db.json", "Полное имя файла, куда сохраняются текущие значения")
	FlagRestore         = flag.Bool("b", true, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
)

type Config struct {
	EnvAddress         string `env:"ADDRESS"`
	EnvReportInterval  int64  `env:"REPORT_INTERVAL"`
	EnvPollInterval    int64  `env:"POLL_INTERVAL"`
	EnvLogLevel        string `env:"LOG_LEVEL"`
	EnvStoreInterval   int64  `env:"STORE_INTERVAL"`
	EnvFileStoragePath string `env:"FILE_STORAGE_PATH"`
	EnvRestore         bool   `env:"RESTORE"`
}

// InitializeGlobals инициализирует флаги на основе значений из конфигурации
func InitializeGlobals() {
	flag.Parse()
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	checkForNil(cfg.EnvAddress, FlagAddressAndPort)
	checkForNil(cfg.EnvReportInterval, FlagRreportInterval)
	checkForNil(cfg.EnvPollInterval, FlagPollInterval)
	checkForNil(cfg.EnvLogLevel, FlagLogLevel)
	checkForNil(cfg.EnvStoreInterval, FlagStoreInterval)
	checkForNil(cfg.EnvFileStoragePath, FlagFileStoragePath)
	checkForNil(cfg.EnvRestore, FlagRestore)
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
