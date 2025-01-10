package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

type ConfigServer struct {
	log *logrus.Entry
	Server
}

type ConfigAgent struct {
	log *logrus.Entry
	Agent
}

type Server struct {
	FlagAddressAndPort  string `env:"ADDRESS"`
	FlagLogLevel        string `env:"LOG_LEVEL"`
	FlagStoreInterval   int64  `env:"STORE_INTERVAL"`
	FlagFileStoragePath string `env:"FILE_STORAGE_PATH"`
	FlagRestore         bool   `env:"RESTORE"`
	FlagDatabaseDsn     string `env:"DATABASE_DSN"`
	FlagKey             string `env:"KEY"`
}

type Agent struct {
	FlagAddressAndPort string `env:"ADDRESS"`
	FlagReportInterval int64  `env:"REPORT_INTERVAL"`
	FlagPollInterval   int64  `env:"POLL_INTERVAL"`
	FlagLogLevel       string `env:"LOG_LEVEL"`
	FlagKey            string `env:"KEY"`
}

func NewConfigServer(log *logrus.Entry) *ConfigServer {
	return &ConfigServer{
		log:    log,
		Server: Server{},
	}
}

func NewConfigAgent(log *logrus.Entry) *ConfigAgent {
	return &ConfigAgent{
		log:   log,
		Agent: Agent{},
	}
}

func (c *ConfigServer) InitializeServerConfig() {
	//Парсит флаги
	c.newVarServerFlag()
}

// newVarServerFlag инициализирует структуру VarServerFlag и парсит флаги командной строки
func (c *ConfigServer) newVarServerFlag() {
	c.Server = Server{
		FlagAddressAndPort:  "localhost:8080",
		FlagLogLevel:        "info",
		FlagStoreInterval:   30,
		FlagFileStoragePath: "C:\\Users\\Сеня\\Desktop\\KursMetrics\\cmd\\server\\metrics-db.json",
		FlagRestore:         false,
		FlagDatabaseDsn:     "host=localhost dbname=postgres user=Senya password=1q2w3e4r5t sslmode=disable",
		FlagKey:             "",
	}

	err := env.Parse(&c.Server)
	if err != nil {
		c.log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}

	// Определение флагов
	flag.StringVar(&c.Server.FlagAddressAndPort, "a", c.Server.FlagAddressAndPort, "Указываем адрес и порт по которому будем подключаться")
	flag.StringVar(&c.Server.FlagLogLevel, "l", c.Server.FlagLogLevel, "Уровень логирования")
	flag.Int64Var(&c.Server.FlagStoreInterval, "i", c.Server.FlagStoreInterval, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&c.Server.FlagFileStoragePath, "f", c.Server.FlagFileStoragePath, "Полное имя файла, куда сохраняются текущие значения")
	flag.BoolVar(&c.Server.FlagRestore, "r", c.Server.FlagRestore, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&c.Server.FlagDatabaseDsn, "d", c.Server.FlagDatabaseDsn, "Строка подключения к базе данных")
	flag.StringVar(&c.Server.FlagKey, "k", c.Server.FlagKey, "Строка подключения к базе данных")
	// Парсинг флагов
	flag.Parse()
}

func (c *ConfigAgent) InitializeAgentConfig() Agent {
	//Парсит флаги и переменные окружения
	return c.newAgentFlag()
}

func (c *ConfigAgent) newAgentFlag() Agent {
	cfgAgentFlags := Agent{
		FlagAddressAndPort: "localhost:8080",
		FlagReportInterval: 10,
		FlagPollInterval:   2,
		FlagLogLevel:       "info",
		FlagKey:            "",
	}

	//Парсит переменные окружения для агента
	err := env.Parse(&cfgAgentFlags)
	if err != nil {
		c.log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}

	// Определение флагов
	flag.StringVar(&cfgAgentFlags.FlagAddressAndPort, "a", cfgAgentFlags.FlagAddressAndPort, "Указываем адрес и порт по которому будем подключаться")
	flag.Int64Var(&cfgAgentFlags.FlagReportInterval, "r", cfgAgentFlags.FlagReportInterval, "Время ожидания перед отправкой в секундах, по умолчанию 10 сек")
	flag.Int64Var(&cfgAgentFlags.FlagPollInterval, "p", cfgAgentFlags.FlagPollInterval, "Частота опроса метрик из пакета runtime в секундах, по умолчанию 2 сек")
	flag.StringVar(&cfgAgentFlags.FlagLogLevel, "l", cfgAgentFlags.FlagLogLevel, "Уровень логирования")
	flag.StringVar(&cfgAgentFlags.FlagKey, "k", cfgAgentFlags.FlagKey, "Строка подключения к базе данных")
	// Парсинг флагов
	flag.Parse()

	return cfgAgentFlags
}
