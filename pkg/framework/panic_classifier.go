package framework

import (
	"runtime"
	"strings"

	frameworkErrors "github.com/Danieljosh-uduma/zen/pkg/framework/errors"
)

func classifyPanic(rec any) (PanicType, PanicSeverity) {

	switch rec := rec.(type) {

	case *frameworkErrors.AppError:
		return PanicTypeOperational, PanicSeverityLow

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
