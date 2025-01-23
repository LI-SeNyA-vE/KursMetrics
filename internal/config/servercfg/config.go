package servercfg

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type ConfigServer struct {
	log *logrus.Entry
	Server
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

// newVarServerFlag инициализирует структуру VarServerFlag и парсит флаги командной строки
func (c *ConfigServer) newVarServerFlag() {
	c.Server = Server{
		FlagAddressAndPort:  "localhost:8080",
		FlagLogLevel:        "debug",
		FlagStoreInterval:   30,
		FlagFileStoragePath: "/Users/senya/GolandProjects/KursMetrics/cmd/config/metrics-db.json",
		FlagRestore:         false,
		FlagDatabaseDsn:     "host=localhost dbname=postgres user=Senya password=1q2w3e4r5t sslmode=disable",
		FlagKey:             "",
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

	err := env.Parse(&c.Server)
	if err != nil {
		c.log.Info("Ошибка на этапе парсинга переменных окружения", err)
	}
}
