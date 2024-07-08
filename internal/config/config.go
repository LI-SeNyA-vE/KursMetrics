package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	AddressAndPort  = flag.String("a", "localhost:8080", "Указываем адресс и порт по которому будем потключаться")
	RreportInterval = flag.Int64("r", 3, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	PollInterval    = flag.Int64("p", 2, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	LogLevel        = flag.String("l", "info", "Уровень логирования")
	StoreInterval   = flag.Int64("i", 0, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	FileStoragePath = flag.String("f", "/tmp/metrics-db.json", "Полное имя файла, куда сохраняются текущие значения")
	Restore         = flag.Bool("b", true, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
)

type Config struct {
	Address         string `env:"ADDRESS"`
	ReportInterval  int64  `env:"REPORT_INTERVAL"`
	PollInterval    int64  `env:"POLL_INTERVAL"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int64  `env:STORE_INTERVAL`
	FileStoragePath string `env:FILE_STORAGE_PATH`
	Restore         bool   `env:RESTORE`
}

func GetConfig() Config {
	flag.Parse()
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func InitializeGlobals(cfg Config) {
	if cfg.Address != "" {
		*AddressAndPort = cfg.Address
	}
	if cfg.ReportInterval != 0 {
		*RreportInterval = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
		*PollInterval = cfg.PollInterval
	}
	if cfg.LogLevel != "" {
		*LogLevel = cfg.LogLevel
	}
	if cfg.StoreInterval != 0 {
		*LogLevel = cfg.LogLevel
	}
	if cfg.FileStoragePath != "" {
		*LogLevel = cfg.LogLevel
	}
	if cfg.Restore {
		*LogLevel = cfg.LogLevel
	}
}
