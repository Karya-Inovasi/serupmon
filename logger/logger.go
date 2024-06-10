package logger

type Logger struct{}

var logger *Logger

func New() *Logger {
	if logger == nil {
		logger = &Logger{}
	}

	return logger
}
