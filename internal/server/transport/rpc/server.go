// Package rpc реализовывает функцию старта gRPC сервера
package rpc

import (
	"fmt"
	pb "github.com/LI-SeNyA-vE/KursMetrics/api/proto/v1/metrics"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/server/storages"
	rpchandler "github.com/LI-SeNyA-vE/KursMetrics/internal/server/transport/rpc/handler"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func StartServerRPC(cfgServer *servercfg.ConfigServer, storage storages.MetricsStorage, log *logrus.Entry) {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricsServiceServer(s, rpchandler.NewMetricsServer(storage, log))

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err = s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
