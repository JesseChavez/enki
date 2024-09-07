package logger

import (
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	log *slog.Logger
}

func New() *Logger {
	nlog := Logger{}
    // sl := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    sl := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
