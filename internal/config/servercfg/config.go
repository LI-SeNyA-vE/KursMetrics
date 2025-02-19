// Package servercfg реализует логику и запуск конфига для сервера.
// В нём разбираются флаги и переменные окружения
package servercfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
	"os"
)

type ConfigServer struct {
	log *logrus.Entry
	Server
}

type Server struct {
	FlagAddressAndPort  string `env:"ADDRESS" json:"address"`
	FlagLogLevel        string `env:"LOG_LEVEL"`
	FlagStoreInterval   int64  `env:"STORE_INTERVAL" json:"store_interval"`
	FlagFileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	FlagRestore         bool   `env:"RESTORE" json:"restore"`
	FlagDatabaseDsn     string `env:"DATABASE_DSN" json:"database_dsn"`
	FlagKey             string `env:"KEY"`
	FlagCryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	FlagConfig          string `env:"CONFIG"`
}

func NewConfigServer(log *logrus.Entry) *ConfigServer {
	return &ConfigServer{
		log:    log,
		Server: Server{},
	}
}
func (c *ConfigServer) InitializeServerConfig() {
	//Парсит флаги
	c.newVarServerFlag()
}

func (c *ConfigServer) newVarServerFlag() {

	// Дефолтные флаги, если ничего не указанно
	flagDefault := Server{
		//FlagAddressAndPort:  "localhost:8080",
		//FlagLogLevel:        "debug",
		//FlagStoreInterval:   30,
		//FlagFileStoragePath: "/Users/senya/GolandProjects/KursMetrics/cmd/config/metrics-db.json",
		//FlagRestore:         false,
		//FlagDatabaseDsn:     "host=localhost dbname=postgres user=Senya password=1q2w3e4r5t sslmode=disable",
		//FlagKey:             "",
		//FlagCryptoKey:       "", //rsaKeys/privateKey.pem
		//FlagConfig:          "",
	}

	// Парсинг переменных окружения
	envParse := Server{}
	err := env.Parse(&envParse)
	if err != nil {
		c.log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}

	// Парсинг флагов
	flagParse := Server{}
	flag.StringVar(&flagParse.FlagAddressAndPort, "a", c.Server.FlagAddressAndPort, "Указываем адрес и порт по которому будем подключаться")
	flag.StringVar(&flagParse.FlagLogLevel, "l", c.Server.FlagLogLevel, "Уровень логирования")
	flag.Int64Var(&flagParse.FlagStoreInterval, "i", c.Server.FlagStoreInterval, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&flagParse.FlagFileStoragePath, "f", c.Server.FlagFileStoragePath, "Полное имя файла, куда сохраняются текущие значения")
	flag.BoolVar(&flagParse.FlagRestore, "r", c.Server.FlagRestore, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&flagParse.FlagDatabaseDsn, "d", c.Server.FlagDatabaseDsn, "Строка подключения к базе данных")
	flag.StringVar(&flagParse.FlagKey, "k", c.Server.FlagKey, "Строка подключения к базе данных")
	flag.StringVar(&flagParse.FlagCryptoKey, "crypto-rsaKeys", c.Server.FlagCryptoKey, "Путь до файла с приватным ключом")
	flag.StringVar(&flagParse.FlagConfig, "c", flagParse.FlagConfig, "Путь к конфигурационному файлу")
	flag.StringVar(&flagParse.FlagConfig, "config", flagParse.FlagConfig, "Путь к конфигурационному файлу")
	flag.Parse()

	configFlag := Server{}
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
			c.log.Errorf("ошибка чтения конфигурационного файла: %w", err)
		} else {
			err = json.Unmarshal(configFile, &configFlag)
			if err != nil {
				c.log.Errorf("ошибка парсинга JSON конфигурационного файла: %w", err)
			}
		}
	}

	// Устанавливаем значение флагов в итоговый конфиг
	c.setConfigServerValue(flagParse, envParse, configFlag, flagDefault)

	fmt.Printf("%v", c.Server)
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
func (c *ConfigServer) setConfigServerValue(flagParse, envParse, configFlag, flagDefault Server) {
	var flagName string
	var flagValue interface{}

	flagValue, flagName = setConfigValue(flagParse.FlagAddressAndPort, envParse.FlagAddressAndPort, configFlag.FlagAddressAndPort, flagDefault.FlagAddressAndPort)
	c.Server.FlagAddressAndPort = flagValue.(string)
	c.log.Infof("Значение для FlagAddressAndPort было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagLogLevel, envParse.FlagLogLevel, configFlag.FlagLogLevel, flagDefault.FlagLogLevel)
	c.Server.FlagLogLevel = flagValue.(string)
	c.log.Infof("Значение для FlagLogLevel было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagStoreInterval, envParse.FlagStoreInterval, configFlag.FlagStoreInterval, flagDefault.FlagStoreInterval)
	c.Server.FlagStoreInterval = flagValue.(int64)
	c.log.Infof("Значение для FlagStoreInterval было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagFileStoragePath, envParse.FlagFileStoragePath, configFlag.FlagFileStoragePath, flagDefault.FlagFileStoragePath)
	c.Server.FlagFileStoragePath = flagValue.(string)
	c.log.Infof("Значение для FlagFileStoragePath было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagRestore, envParse.FlagRestore, configFlag.FlagRestore, flagDefault.FlagRestore)
	c.Server.FlagRestore = flagValue.(bool)
	c.log.Infof("Значение для FlagRestore было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagDatabaseDsn, envParse.FlagDatabaseDsn, configFlag.FlagDatabaseDsn, flagDefault.FlagDatabaseDsn)
	c.Server.FlagDatabaseDsn = flagValue.(string)
	c.log.Infof("Значение для FlagDatabaseDsn было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagKey, envParse.FlagKey, configFlag.FlagKey, flagDefault.FlagKey)
	c.Server.FlagKey = flagValue.(string)
	c.log.Infof("Значение для FlagKey было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagCryptoKey, envParse.FlagCryptoKey, configFlag.FlagCryptoKey, flagDefault.FlagCryptoKey)
	c.Server.FlagCryptoKey = flagValue.(string)
	c.log.Infof("Значение для FlagCryptoKey было взято из %s", flagName)

	flagValue, flagName = setConfigValue(flagParse.FlagConfig, envParse.FlagConfig, configFlag.FlagConfig, flagDefault.FlagConfig)
	c.Server.FlagConfig = flagValue.(string)
	c.log.Infof("Значение для FlagConfig было взято из %s", flagName)
}
