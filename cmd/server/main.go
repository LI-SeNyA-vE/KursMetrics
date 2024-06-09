package main

import (
	"fmt"
	"net/http"

	metricStor "github.com/LI-SeNyA-vE/YaGo/internal/incriment1/handlers"
)

/* func NotUpdate(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "No update", http.StatusBadRequest)
}

func NoGaugeOrCounter(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Это не 'gauge' и не 'counter' запросы", http.StatusBadRequest)
} */

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &metricStor.MyStruct{})
	fmt.Println("Открыт сервер http://localhost:8080")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
