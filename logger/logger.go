package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/JesseChavez/spt"
)

type ILogger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

type Logger struct {
	log *slog.Logger
}

func New(instance string, appName string,appLogLevel string) *Logger {
	nlog := Logger{}
    // sl := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logLevel := slog.LevelDebug

	switch appLogLevel {
	case "info":
		logLevel = slog.LevelInfo
	case "debug":
		logLevel = slog.LevelDebug
	}

	log.Println("log level:", logLevel)

	options := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler *slog.TextHandler

	fileLogging := spt.FetchEnv("LOG_TO_FILE", "")

	if fileLogging == "" {
		handler = slog.NewTextHandler(os.Stdout, options)
	} else {
		output, err := OutputStream(instance, appName)

		if err != nil {
			log.Println("Error opening log file", err)
			log.Println("Fallback to stdout")
			handler = slog.NewTextHandler(os.Stdout, options)
		} else {
			handler = slog.NewTextHandler(output, options)
		}

	}

	sl := slog.New(handler)

	nlog.log = sl

    return &nlog
}

func OutputStream(instance string, appName string) (*File, error) {
	fileDir := "/var/local/log"

	fileName := instance + "-output-" + appName + ".log"

	file, err := NewDaily(fileDir, fileName, nil)

	// defer file.Close()

	return file, err
}

func (sl *Logger) Debug(msg string, keysAndValues ...interface{}) {
    sl.log.Debug(msg, keysAndValues...)
}

func (sl *Logger) Info(msg string, keysAndValues ...interface{}) {
    sl.log.Info(msg, keysAndValues...)
}

func (sl *Logger) Warn(msg string, keysAndValues ...interface{}) {
    sl.log.Warn(msg, keysAndValues...)
}

func (sl *Logger) Error(msg string, keysAndValues ...interface{}) {
    sl.log.Error(msg, keysAndValues...)
}

func (sl *Logger) Fatal(msg string, keysAndValues ...interface{}) {
    sl.log.Error(msg, keysAndValues...)
    log.Fatal(msg)
}
