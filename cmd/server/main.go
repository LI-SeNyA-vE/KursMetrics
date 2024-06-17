package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	metricHandlers "github.com/LI-SeNyA-vE/KursMetrics/internal/incriment1/handlers"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

var (
	addressAndPort = flag.String("a", "localhost:8080", "address and port")
)

type Config struct {
	address string `env:"ADDRESS"`
}

func main() {
	var cfg Config
	r := chi.NewRouter()
	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.address != "" {
		addressAndPort = &cfg.address
	}

	r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", metricHandlers.CorrectPostRequest)
	r.Get("/value/{typeMetric}/{nameMetric}", metricHandlers.CorrectGetRequest)
	r.Get("/", metricHandlers.AllValue)
	fmt.Println("Открыт сервер ", *addressAndPort)
	err = http.ListenAndServe(*addressAndPort, r)
	if err != nil {
		panic(err)
	}
}
