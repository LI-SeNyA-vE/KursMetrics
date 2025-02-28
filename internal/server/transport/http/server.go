// Package rpc реализовывает функцию старта HTTP сервера
package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/storages"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/transport/http/router"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServerHTTP(cfgServer *servercfg.ConfigServer, storage storages.MetricsStorage, log *logrus.Entry) {
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
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Ошибка сервера: %v", err)
	}
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
