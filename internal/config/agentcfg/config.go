// Package agentcfg
package agentcfg

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type ConfigAgent struct {
	log *logrus.Entry
	Agent
}

type Agent struct {
	FlagAddressAndPort string `env:"ADDRESS"`
	FlagReportInterval int64  `env:"REPORT_INTERVAL"`
	FlagPollInterval   int64  `env:"POLL_INTERVAL"`
	FlagLogLevel       string `env:"LOG_LEVEL"`
	FlagKey            string `env:"KEY"`
	FlagRateLimit      int64  `env:"RATE_LIMIT"`
	FlagCryptoKey      string `env:"CRYPTO_KEY"`
}

func NewConfigAgent(log *logrus.Entry) *ConfigAgent {
	return &ConfigAgent{
		log:   log,
		Agent: Agent{},
	}
}

func (c *ConfigAgent) InitializeAgentConfig() Agent {
	//Парсит флаги и переменные окружения
	return c.newAgentFlag()
}

func (c *ConfigAgent) newAgentFlag() Agent {
	c.Agent = Agent{
		FlagAddressAndPort: "localhost:8080",
		FlagReportInterval: 10,
		FlagPollInterval:   2,
		FlagLogLevel:       "debug",
		FlagKey:            "",
		FlagRateLimit:      1,
		FlagCryptoKey:      "rsaKeys/publicKey.pem", //
	}

	// Определение флагов
	flag.StringVar(&c.Agent.FlagAddressAndPort, "a", c.Agent.FlagAddressAndPort, "Указываем адрес и порт по которому будем подключаться")
	flag.Int64Var(&c.Agent.FlagReportInterval, "r", c.Agent.FlagReportInterval, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	flag.Int64Var(&c.Agent.FlagPollInterval, "p", c.Agent.FlagPollInterval, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	flag.StringVar(&c.Agent.FlagLogLevel, "g", c.Agent.FlagLogLevel, "Уровень логирования")
	flag.StringVar(&c.Agent.FlagKey, "k", c.Agent.FlagKey, "Строка подключения к базе данных")
	flag.Int64Var(&c.Agent.FlagRateLimit, "l", c.Agent.FlagRateLimit, "Количество одновременно исходящих запросов на сервер")
	flag.StringVar(&c.Agent.FlagCryptoKey, "crypto-rsaKeys", c.Agent.FlagCryptoKey, "Путь до файла с публичным ключом")
	// Парсинг флагов
	flag.Parse()

	//Парсит переменные окружения для агента
	err := env.Parse(&c.Agent)
	if err != nil {
		c.log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}

	return c.Agent
}
