package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

// writerHook отвечает за запись логов
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// customFormatter кастомизирует вывод логов
type customFormatter struct{}

func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Добавляем цвет времени (ANSI-коды)
	colorReset := "\033[0m"
	colorCyan := "\033[36m" // Цвет для времени (голубой)

	timestamp := fmt.Sprintf("%s • %s%s", colorCyan, entry.Time.Format(time.RFC3339), colorReset)
	level := entry.Level.String()
	message := entry.Message

	funcName := "unknown"
	file := "unknown"

	if entry.HasCaller() {
		funcName = entry.Caller.Function
		file = fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
	}

	log := fmt.Sprintf("%s level=\"%s\" msg=\"%s\" func=\"%s\" file=\"%s\"\n",
		timestamp, level, message, funcName, file)

	return []byte(log), nil
}

func NewLogger() *logrus.Entry {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&customFormatter{})

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		panic(err)
	}

	readDir, err := os.ReadDir("logs")
	if err != nil {
		panic(err)
	}

	var existFile bool
	for _, entry := range readDir {
		if entry.Name() == "all.log" {
			existFile = true
			break
		}
		existFile = false
	}

	if !existFile {
		_, err = os.Create("logs/all.log")
		if err != nil {
			panic(err)
		}
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	return logrus.NewEntry(l)
}
