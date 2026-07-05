package logger

import (
	"fmt"
	"sort"
	"time"

	"github.com/bukasin1/zen/internal/zencore/internal/utils"
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
		utils.ColorGray, time.Now().UTC().Format(time.RFC3339), utils.ColorReset,

		msg,
	)

	fieldsKeys := make([]string, 0, len(fields))
	for k := range fields {
		fieldsKeys = append(fieldsKeys, k)
	}
	sort.Strings(fieldsKeys)

	for _, k := range fieldsKeys {
		v := fields[k]
		switch k {
		case "addr":
			fmt.Printf(" on %shttp://localhost%v%s", utils.ColorYellow, v, utils.ColorReset)
			continue
		case "status":
			statusCol := utils.StatusColor(v.(int))
			v = fmt.Sprintf("%s%-3d%s", statusCol, v, utils.ColorReset)
		case "method":
			v = fmt.Sprintf("%s%-6s%s", utils.ColorBlue, v, utils.ColorReset)
		case "duration":
			v = fmt.Sprintf("%s%-6s%s", utils.ColorGreen, v, utils.ColorReset)
		case "path":
			v = fmt.Sprintf("%s%s%s", utils.ColorYellow, v, utils.ColorReset)
		case "size":
			v = fmt.Sprintf("%s%d%s", utils.ColorGray, v, utils.ColorReset)
		// case "requestID":
		// 	v = fmt.Sprintf("%s%-16s%s", utils.ColorGray, v, utils.ColorReset)
		// case "ip":
		// 	v = fmt.Sprintf("%s%-16s%s", utils.ColorGray, v, utils.ColorReset)
		default:
		}

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
