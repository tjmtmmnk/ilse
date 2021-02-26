package util

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

type logger struct {
	*log.Logger
}

var Logger = newLogger()

func newLogger() *logger {
	path := fmt.Sprintf("%s/log.txt", GetHomeDir())
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	l := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &logger{
		l,
	}
}

func (l *logger) Info(v ...interface{}) {
	l.Print(PREFIX_INFO, v, "\033[39;49m\n")
}

func (l *logger) Warn(v ...interface{}) {
	l.Print(PREFIX_WARN, v, "\033[39;49m\n")
}

func (l *logger) Error(v ...interface{}) {
	l.Print(PREFIX_ERROR, v, "\033[39;49m\n")
	os.Exit(0)
}

func (l *logger) Debug(v ...interface{}) {
	l.Print(PREFIX_DEBUG, v, "\033[39;49m\n")
}
