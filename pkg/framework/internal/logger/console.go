package logger

import (
	"encoding/json"
	"fmt"
	"maps"
	"time"
)

type ConsoleLogger struct {
	pretty bool
}

func NewConsoleLogger(pretty bool) *ConsoleLogger {
	return &ConsoleLogger{pretty: pretty}
}

func (l *ConsoleLogger) log(level string, msg string, fields Fields) {
	entry := map[string]any{
		"level":   level,
		"message": msg,
		"time":    time.Now().UTC().Format(time.RFC3339),
	}

	maps.Copy(entry, fields)

	var b []byte
	var err error

	if l.pretty {
		b, err = json.MarshalIndent(entry, "", "  ")
	} else {
		b, err = json.Marshal(entry)
	}

	if err != nil {
		fmt.Printf(`{"level":"ERROR","message":"failed to marshal log","error":"%v"}`+"\n", err)
		return
	}
	fmt.Println(string(b))
}

func (l *ConsoleLogger) Info(msg string, fields Fields) {
	l.log("INFO", msg, fields)
}

func (l *ConsoleLogger) Error(msg string, fields Fields) {
	l.log("ERROR", msg, fields)
}

func (l *ConsoleLogger) Debug(msg string, fields Fields) {
	l.log("DEBUG", msg, fields)
}

func (l *ConsoleLogger) Warn(msg string, fields Fields) {
	l.log("WARN", msg, fields)
}
