package zencore

import (
	"fmt"

	frameworkErrors "github.com/bukasin1/zen/internal/zencore/errors"
	"github.com/bukasin1/zen/internal/zencore/logger"
)

type PanicHandler interface {
	Handle(*Context, *PanicInfo)
}

// DefaultPanicHandler is the default panic handler for the application.
type DefaultPanicHandler struct{}

// Handle handles the panic and returns an error.
func (h *DefaultPanicHandler) Handle(c *Context, info *PanicInfo) {

	fields := logger.Fields{
		"path":      info.Path,
		"method":    info.Method,
		"requestID": info.RequestID,
		"panic":     info.Value,
		"panicType": info.Type,
		"severity":  info.Severity,
	}

	if c.app.RecoveryConfig.IncludeStack {
		fields["stackTrace"] = string(info.StackTrace)
	}

	c.LogError("Request panicked", fields)

	// Cannot safely write another response.
	if c.responseCommitted.Load() {
		return
	}

	message := "internal server error"
	var details any

	switch err := info.Value.(type) {

	case *frameworkErrors.AppError:
		c.Error(err.Status, err.Message, err.Code, err.Details)
		return

	case *frameworkPanic:
		c.LogError("Framework invariant violated", logger.Fields{
			"panic": info.Value,
		})

		if c.app.RecoveryConfig.ExposeError {
			message = err.Error()
		}

	case error:
		if c.app.RecoveryConfig.ExposeError {
			message = err.Error()
		}

	default:
		if c.app.RecoveryConfig.ExposeError {
			message = fmt.Sprintf("panic occurred: %v", info.Value)
		}
	}

	if c.app.RecoveryConfig.IncludeStack {
		details = string(info.StackTrace)
	}

	c.Error(
		500,
		message,
		frameworkErrors.ErrInternal,
		details,
	)
}
