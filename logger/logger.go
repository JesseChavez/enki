package logger

import (
	"log"
	"log/slog"
	"os"
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

    sl := slog.New(slog.NewTextHandler(os.Stdout, options))

	nlog.log = sl

    return &nlog
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
