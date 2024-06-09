package handlers

import (
	"fmt"
	"net/http"
	"strings"

	serverMetric "github.com/LI-SeNyA-vE/YaGo/internal/incriment1/server"
	storageMetric "github.com/LI-SeNyA-vE/YaGo/internal/incriment1/storage"
)

type MyStruct storageMetric.MetricStorage

func (m *MyStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL: ", r.URL.Path)        //Можно удалить
	parts := strings.Split(r.URL.Path, "/") //Разделение приходящего URL по знаку '/'
	if serverMetric.ValidationLengthsURL(parts, w) {
		return
	}

	update, typeMetric, nameMetric, countMetric := parts[1], parts[2], parts[3], parts[4]
	if serverMetric.ValidationFirstElementURL(update, "update", w) {
		return
	}
	if serverMetric.ValidationTypeMetric(typeMetric, nameMetric, countMetric, w) {
		return
	}
	w.WriteHeader(http.StatusOK) //Отправляет ответ что всё ОК
}
