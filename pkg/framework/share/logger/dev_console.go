package logger

import (
	"fmt"
	"time"

	"github.com/Danieljosh-uduma/zen/pkg/framework/internal/utils"
)

type DevConsoleLogger struct{}

func NewDevConsoleLogger() *DevConsoleLogger {
	return &DevConsoleLogger{}
}

func (l *DevConsoleLogger) log(level string, color string, msg string, fields Fields) {
	fmt.Printf(
		"%s%s%s [%s%s%s] %s",
		// color and level
		color, level, utils.ColorReset,
		// time
		utils.ColorGray, time.Now().UTC().Format("15:04:05"), utils.ColorReset,

		msg,
	)

	for k, v := range fields {
		fmt.Printf(" %s=%v", k, v)
	}

	fmt.Println()
}

// Info logs an info message
func (l *DevConsoleLogger) Info(msg string, fields Fields) {
	l.log("INFO", utils.ColorGreen, msg, fields)
}

// Error logs an error message
func (l *DevConsoleLogger) Error(msg string, fields Fields) {
	l.log("ERROR", utils.ColorRed, msg, fields)
}

// Debug logs a debug message
func (l *DevConsoleLogger) Debug(msg string, fields Fields) {
	l.log("DEBUG", utils.ColorCyan, msg, fields)
}

// Warn logs a warning message
func (l *DevConsoleLogger) Warn(msg string, fields Fields) {
	l.log("WARN", utils.ColorYellow, msg, fields)
}
