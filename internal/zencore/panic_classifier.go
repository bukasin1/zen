package zencore

import (
	"runtime"
	"strings"

	frameworkErrors "github.com/bukasin1/zen/internal/zencore/errors"
)

func classifyPanic(rec any) (PanicType, PanicSeverity) {

	switch rec := rec.(type) {

	case *frameworkErrors.AppError:
		return PanicTypeOperational, PanicSeverityLow

	case *frameworkPanic:
		return PanicTypeFramework, PanicSeverityCritical

	case runtime.Error:
		return PanicTypeProgrammer, PanicSeverityHigh

	case error:
		return PanicTypeOperational, PanicSeverityMedium

	case string:
		value := strings.ToLower(rec)

		if strings.Contains(value, "framework") {
			return PanicTypeFramework, PanicSeverityCritical
		}

		return PanicTypeUnknown, PanicSeverityMedium

	default:
		return PanicTypeUnknown, PanicSeverityMedium
	}
}
