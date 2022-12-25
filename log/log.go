package log

import (
	"fmt"
	"log"

	"github.com/2gis-demo-app/cfg"
)

type Logger interface {
	LogErrorf(format string, v ...any)
	LogInfo(format string, v ...any)
}

type MyLogger struct {
	logger *log.Logger
}

func NewLogger(config cfg.Config) Logger {
	//TODO: configure log
	return &MyLogger{
		logger: log.Default(),
	}
}

func (l *MyLogger) LogErrorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v)
	l.logger.Printf("[Error]: %s\n", msg)
}

func (l *MyLogger) LogInfo(format string, v ...any) {
	msg := fmt.Sprintf(format, v)
	l.logger.Printf("[Info]: %s\n", msg)
}
