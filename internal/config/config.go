package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	AddressAndPort  = flag.String("a", "localhost:8080", "Указываем адресс и порт по которому будем потключаться")
	RreportInterval = flag.Int64("r", 10, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	PollInterval    = flag.Int64("p", 2, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	LogLevel        = flag.String("l", "info", "Уровень логирования")
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
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

}
