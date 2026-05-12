package framework

import "time"

type PanicType string

const (
	PanicTypeOperational PanicType = "operational"
	PanicTypeProgrammer  PanicType = "programmer"
	PanicTypeFramework   PanicType = "framework"
	PanicTypeUnknown     PanicType = "unknown"
)

type PanicSeverity string

const (
	PanicSeverityLow      PanicSeverity = "low"
	PanicSeverityMedium   PanicSeverity = "medium"
	PanicSeverityHigh     PanicSeverity = "high"
	PanicSeverityCritical PanicSeverity = "critical"
)

type PanicInfo struct {
	Value      any
	Type       PanicType
	Severity   PanicSeverity
	StackTrace []byte
	RequestID  string
	Path       string
	Method     string
	Timestamp  time.Time
}
