package logger

type Logger interface {
	RotateLogFile() error
	CleanupOldLogFiles() error

	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}
