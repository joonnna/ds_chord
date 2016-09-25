package logger

import (
	"log"
	"io"
	"fmt"
)

type Logger struct {
	log *log.Logger
	prefix string
}

var (
	err = "ERROR"
	info = "INFO"
	debug = "DEBUG"
	yellow = "\x1b[33m"
	green = "\x1b[32m"
	red = "\x1b[31m"
	reset = "\x1b[0m"
)

func (l *Logger) formatString(level string, str string, color string) string {
	return fmt.Sprintln(color + level + "[" + l.prefix + "]" + ": " + reset + str)
}

func (l *Logger) Init(output io.Writer, preFix string, flag int) {
	l.log = log.New(output, "", flag)
	l.prefix = preFix
}

func (l *Logger) Error(msg string) {
	l.log.Println(l.formatString(err, msg, red))
}

func (l *Logger) Info(msg string) {
	l.log.Println(l.formatString(info, msg, green))
}

func (l *Logger) Debug(msg string) {
	l.log.Println(l.formatString(debug, msg, yellow))
}


