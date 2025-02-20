/*
Package funcserver реализует логику сервера в проекте KursMetrics.
Основная задача сервера – принимать метрики от агента (и потенциально от других источников),
хранить их (в PostgreSQL, файле или памяти) и предоставлять API для чтения и обновления.
Функция Run является точкой входа для сервера: она инициализирует конфигурацию, логгер,
попытки подключения к базе данных и, при неудаче, файл-хранилищу, а затем запускает веб-сервер.
*/
package funcserver

import (
	"context"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/delivery/router"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages/database/postgresql"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages/filemetric"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/logger"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/rsakey"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Run инициализирует настройки сервера (флаги, окружение), логгер,
// пытается подключиться к базе данных (до 3 раз), либо использует файловое хранилище.
// При необходимости загружает сохранённые метрики и поднимает HTTP-сервер (router)
// на адресе, указанном в конфигурации. Также запускается pprof на порту :6060.
func Run() {
	var err error
	var storage storages.MetricsStorage
	var wg sync.WaitGroup

	// Запуск pprof на localhost:6060 для профилирования.
	go func() {
		log.Println("pprof запущен на :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Не удалось запустить pprof: %v", err)
		}
	}()

	// Инициализация логгера.
	log := logger.NewLogger()

	// Создаём конфиг
	cfgServer := servercfg.NewConfigServer(log)

	//Парсим флаги и переменный окружения
	cfgServer.InitializeServerConfig()

	// Проверяем на наличие приватного ключа
	if cfgServer.FlagCryptoKey != "" {
		err = rsakey.CheckKey(cfgServer.FlagCryptoKey)
		if err != nil {
			err = rsakey.GenerateAndSaveKeys(cfgServer.FlagCryptoKey)
			if err != nil {
				log.Errorf("ошибка при создании пары ключей: %v", err)
			} else {
				log.Info("Успешно созданы пары ключей")
			}
			//TODO на агенте сделать горутинку, которая будет проверять правильность открытого ключа и если он не правильный, то кидать запросы на сервере на отправку открытого ключа и не выполнять никаких других действий пока не получит ключ
		}
	}

	// Попытка подключения к БД PostgreSQL несколько раз.
	for i := 0; i < 3; i++ {
		storage, err = postgresql.NewConnectDB(log, cfgServer.Server)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		// Если БД недоступна, используем файловое хранилище.
		storage, err = filemetric.NewFileStorage(cfgServer.Server)
		if err != nil {
			log.Info(fmt.Errorf("ошибка при объявлении хранения в файле err: %s", err))
			// При необходимости можно переключиться на in-memory:
			// storage = memorymetric.NewMetricStorage()
		}
	}

	// Загрузка метрик из хранилища (если нужна реставрация из файла или иное).
	err = storage.LoadMetric()
	if err != nil {
		log.Info(err)
		// Также можно переключиться на in-memory при ошибке загрузки:
		// storage = memorymetric.NewMetricStorage()
	}

	// Создаём и настраиваем роутер (HTTP-маршруты и middleware).
	r := router.NewRouter(log, cfgServer.Server, storage)
	r.SetupRouter()

	server := &http.Server{
		Addr:    cfgServer.FlagAddressAndPort,
		Handler: r.Mux,
	}

	go func() {
		handleSignals(server)
	}()

	// Запуск HTTP-сервера на сконфигурированном адресе.
	log.Info("Открыт сервер на ", cfgServer.FlagAddressAndPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка сервера: %v", err)
	}

	wg.Wait()
}

func handleSignals(server *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-signalChan
	fmt.Println("Сервер: получен сигнал завершения, завершаем работу...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Ошибка при завершении сервера: %v", err)
	}
	fmt.Println("Сервер: завершение работы успешно.")
}
