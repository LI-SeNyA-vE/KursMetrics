/*
Package logger предоставляет возможности для настроенного логирования на основе logrus.
В частности, он включает в себя:
  - writerHook, записывающий логи одновременно в несколько io.Writer (в файл и stdout).
  - customFormatter, кастомизирующий формат логов (с выводом даты/времени, уровня, сообщения, функции и файла).
  - NewLogger, создающий и настраивающий готовый объект logrus.Entry с нужными хуками и форматтером.
*/
package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

// writerHook отвечает за запись логов в несколько io.Writer.
// Содержит срез Writer для вывода и набор уровней логирования, при которых Hook будет срабатывать.
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

// Fire вызывается при каждом лог-сообщении. Сериализует лог-запись в строку
// и отправляет её во все объекты из Writer.
func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, _ = w.Write([]byte(line)) // игнорируем возможную ошибку записи
	}
	return err
}

// Levels возвращает список уровней логирования, при которых hook будет вызван.
func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// customFormatter кастомизирует формат логирования:
// дополняет лог цветным временем, выводит функцию и файл, откуда был вызван лог.
type customFormatter struct{}

// Format задаёт пользовательский формат строки лога, включая цветное время (ANSI),
// уровень, сообщение, имя функции и файл/номер строки.
func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
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

	logLine := fmt.Sprintf("%s level=\"%s\" msg=\"%s\" func=\"%s\" file=\"%s\"\n",
		timestamp, level, message, funcName, file)

	return []byte(logLine), nil
}

// NewLogger создаёт и настраивает logrus.Entry со следующими параметрами:
//   - Включён вызов caller (SetReportCaller(true)) для отображения функции/файла,
//   - Установлен customFormatter,
//   - Логи пишутся в файл logs/all.log и в stdout,
//   - Уровень логирования установлен на TraceLevel.
//
// Функция также обеспечивает создание директории logs/ и файла all.log, если они ещё не существуют.
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

	// Перенаправляем весь вывод логгера в io.Discard,
	// поскольку запись будем вести через writerHook.
	l.SetOutput(io.Discard)

	// Используем hook для одновременной записи в файл и stdout.
	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	// Устанавливаем максимально детальный уровень логирования.
	l.SetLevel(logrus.TraceLevel)

	return logrus.NewEntry(l)
}
