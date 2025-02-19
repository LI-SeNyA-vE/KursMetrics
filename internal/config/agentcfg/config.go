// Package agentcfg
package agentcfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"os"
)

type ConfigAgent struct {
	log *logrus.Entry
	Agent
}

type Agent struct {
	FlagAddressAndPort string `env:"ADDRESS" json:"address"`
	FlagReportInterval int64  `env:"REPORT_INTERVAL" json:"report_interval"`
	FlagPollInterval   int64  `env:"POLL_INTERVAL" json:"poll_interval"`
	FlagLogLevel       string `env:"LOG_LEVEL"`
	FlagKey            string `env:"KEY"`
	FlagRateLimit      int64  `env:"RATE_LIMIT"`
	FlagCryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	FlagConfig         string `env:"CONFIG"`
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
	flagDefault := Agent{
		FlagAddressAndPort: "localhost:8080",
		FlagReportInterval: 10,
		FlagPollInterval:   2,
		FlagLogLevel:       "debug",
		FlagKey:            "",
		FlagRateLimit:      1,
		FlagCryptoKey:      "", //rsaKeys/publicKey.pem
		FlagConfig:         "",
	}

	// Парсинг переменных окружения
	envParse := Agent{}
	err := env.Parse(&envParse)
	if err != nil {
		c.log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}

	// Парсинг флагов
	flagParse := Agent{}
	flag.StringVar(&flagParse.FlagAddressAndPort, "a", flagParse.FlagAddressAndPort, "Указываем адрес и порт по которому будем подключаться")
	flag.Int64Var(&flagParse.FlagReportInterval, "r", flagParse.FlagReportInterval, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	flag.Int64Var(&flagParse.FlagPollInterval, "p", flagParse.FlagPollInterval, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	flag.StringVar(&flagParse.FlagLogLevel, "g", flagParse.FlagLogLevel, "Уровень логирования")
	flag.StringVar(&flagParse.FlagKey, "k", flagParse.FlagKey, "Строка подключения к базе данных")
	flag.Int64Var(&flagParse.FlagRateLimit, "l", flagParse.FlagRateLimit, "Количество одновременно исходящих запросов на сервер")
	flag.StringVar(&flagParse.FlagCryptoKey, "crypto-rsaKeys", flagParse.FlagCryptoKey, "Путь до файла с публичным ключом")
	flag.StringVar(&flagParse.FlagConfig, "c", flagParse.FlagConfig, "Путь к конфигурационному файлу")
	flag.StringVar(&flagParse.FlagConfig, "config", flagParse.FlagConfig, "Путь к конфигурационному файлу")
	flag.Parse()

	configFlag := Agent{}
	pathConfig := flagDefault.FlagConfig
	if pathConfig == "" {
		pathConfig = envParse.FlagConfig
		if pathConfig == "" {
			pathConfig = flagParse.FlagConfig
		}
	}

	if pathConfig != "" {
		configFile, err := os.ReadFile(pathConfig)
		if err != nil {
			c.log.Errorf("ошибка чтения конфигурационного файла: %v", err)
		} else {
			err = json.Unmarshal(configFile, &configFlag)
			if err != nil {
				c.log.Errorf("ошибка парсинга JSON конфигурационного файла: %v", err)
			}
		}
	}

	// Устанавливаем значение флагов в итоговый конфиг
	c.setConfigServerValue(flagParse, envParse, configFlag, flagDefault)

	fmt.Printf("%v", c.Agent)
	return c.Agent
}

// setConfigValue - Функция установки значений с приоритетом (флаг -> env -> config -> default)
func setConfigValue(flagValue, envValue, configValue, defaultValue interface{}) (interface{}, string) {
	if flagValue != nil {
		switch v := flagValue.(type) {
		case string:
			if v != "" {
				return v, "флага установленного в консоли"
			}
		case int64:
			if v != 0 {
				return v, "флага установленного в консоли"
			}
		case bool:
			if v {
				return v, "флага установленного в консоли"
			}
		}
	}

	if envValue != nil {
		switch v := envValue.(type) {
		case string:
			if v != "" {
				return v, "переменной окружения"
			}
		case int64:
			if v != 0 {
				return v, "переменной окружения"
			}
		case bool:
			if v {
				return v, "переменной окружения"
			}
		}
	}

	if configValue != nil {
		switch v := configValue.(type) {
		case string:
			if v != "" {
				return v, "конфигурационного файла"
			}
		case int64:
			if v != 0 {
				return v, "конфигурационного файла"
			}
		case bool:
			if v {
				return v, "конфигурационного файла"
			}
		}
	}

	// Если ничего не найдено, берём значение по умолчанию
	return defaultValue, "значений по умолчанию"
}

// setConfigServerValue К каждому флагу запускает функцию на установки значений с приоритетом.
// Заполняет значениями
func (c *ConfigAgent) setConfigServerValue(flagParse, envParse, configFlag, flagDefault Agent) {
	var flagName string
	var flagValue interface{}

	flagValue, flagName = setConfigValue(flagParse.FlagAddressAndPort, envParse.FlagAddressAndPort, configFlag.FlagAddressAndPort, flagDefault.FlagAddressAndPort)
	c.Agent.FlagAddressAndPort = flagValue.(string)
	c.log.Infof("Значение для FlagAddressAndPort было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagReportInterval, envParse.FlagReportInterval, configFlag.FlagReportInterval, flagDefault.FlagReportInterval)
	c.Agent.FlagReportInterval = flagValue.(int64)
	c.log.Infof("Значение для FlagStoreInterval было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagPollInterval, envParse.FlagPollInterval, configFlag.FlagPollInterval, flagDefault.FlagPollInterval)
	c.Agent.FlagPollInterval = flagValue.(int64)
	c.log.Infof("Значение для FlagFileStoragePath было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagLogLevel, envParse.FlagLogLevel, configFlag.FlagLogLevel, flagDefault.FlagLogLevel)
	c.Agent.FlagLogLevel = flagValue.(string)
	c.log.Infof("Значение для FlagLogLevel было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagKey, envParse.FlagKey, configFlag.FlagKey, flagDefault.FlagKey)
	c.Agent.FlagKey = flagValue.(string)
	c.log.Infof("Значение для FlagRestore было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagRateLimit, envParse.FlagRateLimit, configFlag.FlagRateLimit, flagDefault.FlagRateLimit)
	c.Agent.FlagRateLimit = flagValue.(int64)
	c.log.Infof("Значение для FlagDatabaseDsn было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagCryptoKey, envParse.FlagCryptoKey, configFlag.FlagCryptoKey, flagDefault.FlagCryptoKey)
	c.Agent.FlagCryptoKey = flagValue.(string)
	c.log.Infof("Значение для FlagCryptoKey было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagConfig, envParse.FlagConfig, configFlag.FlagConfig, flagDefault.FlagConfig)
	c.Agent.FlagConfig = flagValue.(string)
	c.log.Infof("Значение для FlagConfig было взято из %s", flagName)
}
