package ilse

import (
	"fmt"
	"log"
	"os"
)

const (
	PREFIX_INFO  = "\033[39;49m\033[32mINFO\t"
	PREFIX_WARN  = "\033[39;49m\033[33mWARN\t"
	PREFIX_ERROR = "\033[39;49m\033[31mERROR\t"
	PREFIX_DEBUG = "\033[39;49m\033[36mDEBUG\t"
)

type PrintLevel int

const (
	Info PrintLevel = iota
	Warn
	Error
	Debug
)

type Logger struct {
	*log.Logger
}

func newLogger() (*Logger, error) {
	path := fmt.Sprintf("%s/log.txt", conf.homeDir)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		logger,
	}, nil
}

func (l *Logger) Info(v ...interface{}) {
	l.Print(PREFIX_INFO, v, "\033[39;49m\n")
}

func (l *Logger) Warn(v ...interface{}) {
	l.Print(PREFIX_WARN, v, "\033[39;49m\n")
}

func (l *Logger) Error(v ...interface{}) {
	l.Print(PREFIX_ERROR, v, "\033[39;49m\n")
	os.Exit(0)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Print(PREFIX_DEBUG, v, "\033[39;49m\n")
}
