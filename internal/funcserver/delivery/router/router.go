/*
Package router содержит тип Router, который настраивает и регистрирует
все маршруты (эндпоинты) для HTTP-сервера, а также подключает цепочку
middleware. В частности, здесь определяются пути для приёма и вывода
метрик (как в URL-форме, так и в JSON-формате), пинг для проверки БД
и т.д.
*/
package router

import (
	"github.com/LI-SeNyA-vE/KursMetrics/internal/config/servercfg"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/delivery/handlers"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/delivery/middleware"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/funcserver/storages"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// Router оборачивает chi.Mux, хранит ссылку на конфигурацию сервера (servercfg.Server),
// логгер и абстракцию для работы с метриками (storage). Это позволяет
// настраивать маршруты и подключать необходимые middleware-функции.
type Router struct {
	log *logrus.Entry
	servercfg.Server
	storage storages.MetricsStorage
	*chi.Mux
}

// NewRouter создаёт новую структуру Router, инициализируя её логгером, конфигурацией и хранилищем метрик.
// По умолчанию Mux будет равен nil, инициализация роутов происходит в методе SetupRouter.
func NewRouter(log *logrus.Entry, cfg servercfg.Server, storages storages.MetricsStorage) *Router {
	return &Router{
		log:     log,
		Server:  cfg,
		storage: storages,
		Mux:     nil,
	}
}

// SetupRouter инициализирует chi.Router и регистрирует следующие middleware:
//
//   - LoggingMiddleware: логирование запросов (URI, статус, время);
//   - HashSHA256: проверка HMAC SHA256, если ключ задан;
//   - GunzipMiddleware: распаковка входящих gzip-запросов;
//   - GzipMiddleware: сжатие ответов при Accept-Encoding: gzip.
//
// Затем маппит пути для разных действий:
//
//   - POST /update/{typeMetric}/{nameMetric}/{countMetric} — обновление метрики по URL.
//   - POST /value/ — получение значения метрики в JSON-формате.
//   - POST /update/ — обновление метрики в JSON-формате.
//   - POST /updates/ — обновление массива метрик (batch).
//   - GET /value/{typeMetric}/{nameMetric} — получение метрики по URL.
//   - GET /ping — вывод всех метрик (в данном случае роут переопределён дважды: для вывода метрик
//     и проверки БД, но конечный хендлер тот же / см. комментарий в коде).
//   - GET / — вывод всех метрик.
//   - GET /ping — проверка состояния БД.
//
// По окончании инициализации метод присваивает созданный chi.Router
// в поле Mux и готов к использованию в http.ListenAndServe().
func (rout *Router) SetupRouter() {
	rout.Mux = chi.NewRouter()
	mw := middleware.NewMiddleware(rout.log, rout.Server)
	hl := handlers.NewHandler(rout.log, rout.Server, rout.storage)

	// Подключаем middleware в цепочку
	rout.Mux.Use(mw.LoggingMiddleware)
	if rout.Server.FlagTrustedSubnet != "" {
		rout.Mux.Use(mw.TrustedSubnet)
	}
	rout.Mux.Use(mw.RsaDecoder)
	rout.Mux.Use(mw.HashSHA256)
	rout.Mux.Use(mw.GunzipMiddleware)
	rout.Mux.Use(mw.GzipMiddleware)

	// Регистрация хендлеров
	rout.Mux.Post("/update/{typeMetric}/{nameMetric}/{countMetric}", hl.PostAddValue)
	rout.Mux.Post("/value/", hl.JSONValue)
	rout.Mux.Post("/update/", hl.JSONUpdate)
	rout.Mux.Post("/updates/", hl.PostAddArrayMetrics)
	rout.Mux.Get("/value/{typeMetric}/{nameMetric}", hl.GetReceivingMetric)
	rout.Mux.Get("/ping", hl.GetReceivingAllMetric) // Похоже, используется для вывода метрик
	rout.Mux.Get("/", hl.GetReceivingAllMetric)
	rout.Mux.Get("/ping", hl.Ping) // Похоже, второй роут /ping для проверки БД
}
