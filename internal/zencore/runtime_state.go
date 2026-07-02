package zencore

import (
	"sync/atomic"
	"time"
)

type RuntimeState struct {
	startedAt time.Time

	isReady        atomic.Bool
	isShuttingDown atomic.Bool
}

func NewRuntimeState() *RuntimeState {
	state := &RuntimeState{
		startedAt: time.Now(),
	}

	state.isReady.Store(false)
	state.isShuttingDown.Store(false)

	return state
}

func (r *RuntimeState) StartedAt() time.Time {
	return r.startedAt
}

func (r *RuntimeState) Uptime() time.Duration {
	return time.Since(r.startedAt)
}

func (r *RuntimeState) IsReady() bool {
	return r.isReady.Load()
}

func (r *RuntimeState) SetReady(value bool) {
	r.isReady.Store(value)
}

func (r *RuntimeState) IsShuttingDown() bool {
	return r.isShuttingDown.Load()
}

func (r *RuntimeState) SetShuttingDown(value bool) {
	r.isShuttingDown.Store(value)
}
