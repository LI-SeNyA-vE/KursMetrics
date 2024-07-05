package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getCustomLoggerConfig(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.TimeKey = "time"

	config := zap.Config{
		Level:            lvl,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return config.Build()
}

func Initialize(level string) error {
	var err error
	Log, err = getCustomLoggerConfig(level)
	if err != nil {
		return err
	}
	return nil
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		uri    string
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggingMiddleware(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // функция Now() возвращает текущее время
		responseData := &responseData{
			status: 0,
			uri:    r.RequestURI, // Захват URI из запроса
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)           // обслуживание оригинального запроса
		duration := time.Since(start) // Since возвращает разницу во времени между start

		Log.Sugar().Infoln(
			"uri:", r.RequestURI,
			"method:", r.Method,
			"status:", responseData.status, // получаем перехваченный код статуса ответа
			"duration:", duration,
			"size:", len(responseData.uri), // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn) // возвращаем функционально расширенный хендлер
}
