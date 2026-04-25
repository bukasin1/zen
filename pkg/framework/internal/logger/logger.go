package logger

type Fields map[string]any

type Logger interface {
	Info(msg string, fields Fields)
	Error(msg string, fields Fields)
	Debug(msg string, fields Fields)
	Warn(msg string, fields Fields)
}
