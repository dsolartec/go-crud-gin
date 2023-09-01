package logger

import "fmt"

type localLogger struct{}

func (logger *localLogger) Debugf(message string, args ...any) {
	fmt.Printf("[DEBUG] %s\r\n", fmt.Sprintf(message, args...))
}

func (logger *localLogger) Infof(message string, args ...any) {
	fmt.Printf("[INFO] %s\r\n", fmt.Sprintf(message, args...))
}

func NewLocalLogger() Logger {
	return &localLogger{}
}
