package logging

import (
	"fmt"
	"log"
	"strings"
)


type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

func LevelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func StringToLevel(level string) LogLevel {
	switch level {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

type Logger interface {
	Info(message ...string)
	Error(message ...string)
	Debug(message ...string)
	Warn(message ...string)
	SetLevel(level LogLevel)
}

type CoreLogger struct {
	level LogLevel
}

func NewLogger(logLevel LogLevel) Logger {
	return &CoreLogger{level: logLevel}
}

func (l *CoreLogger) logMessage(level LogLevel, message ...string) {
	if l.level <= level {
		printMessage := strings.Join(message, " ")
		l.printLog(level, printMessage)
	}
}

func (l *CoreLogger) printLog(level LogLevel, message ...string) {
	payload := fmt.Sprintf("[%s]: %s", LevelToString(level), strings.Join(message, " "))
	log.Println(payload)
}

func (l *CoreLogger) Info(message ...string) {
	l.logMessage(INFO, message...)
}

func (l *CoreLogger) Error(message ...string) {
	l.logMessage(ERROR, message...)
}

func (l *CoreLogger) Debug(message ...string) {
	l.logMessage(DEBUG, message...)
}

func (l *CoreLogger) Warn(message ...string) {
	l.logMessage(WARNING, message...)
}

func (l *CoreLogger) SetLevel(level LogLevel) {
	l.level = level
}
