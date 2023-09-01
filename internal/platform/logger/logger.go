package logger

type Logger interface {
	Debugf(message string, args ...any)
	Infof(message string, args ...any)
}
