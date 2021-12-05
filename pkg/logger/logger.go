package logger

import (
	"log"
	"os"
)

const defaultLogPath = "/var/log/camera-status.log"

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type MyLogger struct {
	logger *log.Logger
}

func NewMyLogger() *MyLogger {
	f, err := os.OpenFile(defaultLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
		os.Exit(42)
	}
	return &MyLogger{logger: log.New(f, "", 5)}
}

func (ml *MyLogger) Info(args ...interface{}) {
	var i []interface{}
	i = append(i, "Info:")
	i = append(i, args...)
	ml.logger.Println(i...)
}

func (ml *MyLogger) Error(args ...interface{}) {
	var i []interface{}
	i = append(i, "Error:")
	i = append(i, args...)
	ml.logger.Println(i...)
}
