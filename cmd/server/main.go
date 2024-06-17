package main

import (
	"flag"
	"fmt"
	"net/http"

	metricHandlers "github.com/LI-SeNyA-vE/KursMetrics/internal/incriment1/handlers"
	"github.com/go-chi/chi/v5"
)

var (
	addressAndPort = flag.String("a", "localhost:8080", "address and port")
)

func main() {
	r := chi.NewRouter()

	flag.Parse()
	r.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", metricHandlers.CorrectPostRequest)
	r.Get("/value/{typeMetric}/{nameMetric}", metricHandlers.CorrectGetRequest)
	r.Get("/", metricHandlers.AllValue)
	fmt.Println("Открыт сервер ", *addressAndPort)
	err := http.ListenAndServe(*addressAndPort, r)
	if err != nil {
		panic(err)
	}
}
