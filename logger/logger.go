package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/JesseChavez/spt"
)

type Logger struct {
	log *slog.Logger
}

func New(appLogLevel string) *Logger {
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

	output := OutputStream()

	handler := slog.NewTextHandler(output, options)

	sl := slog.New(handler)

	nlog.log = sl

    return &nlog
}

func OutputStream() *os.File {
	fileLogging := spt.FetchEnv("LOG_TO_FILE", "")

	if fileLogging == "" {
		return os.Stdout
	}

	filePath := "/home/jessec/bryk/head/go_fm_data_lifecycle_manager/application.log"

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Println("Failed to open log file", err)
		log.Println("Using stdout to log")
		return os.Stdout
	}

	// defer file.Close()

	return file
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
