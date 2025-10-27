package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

type Snapshot struct {
	InFlight   int64     `json:"in_flight"`
	Processed  uint64    `json:"processed"`
	Failed     uint64    `json:"failed"`
	StartedAt  time.Time `json:"started_at"`
	Processing []string  `json:"processing"`
}

type Metrics struct {
	inFlight  atomic.Int64
	processed atomic.Uint64
	failed    atomic.Uint64
	started   time.Time
	requests  sync.Map
}

func New() *Metrics {
	return &Metrics{
		started: time.Now(),
	}
}

func (m *Metrics) IncInFlight() {
	if m == nil {
		return
	}
	m.inFlight.Add(1)
}

func (m *Metrics) DecInFlight() {
	if m == nil {
		return
	}
	m.inFlight.Add(-1)
}

func (m *Metrics) IncProcessed() {
	if m == nil {
		return
	}
	m.processed.Add(1)
}

func (m *Metrics) IncFailed() {
	if m == nil {
		return
	}
	m.failed.Add(1)
}

func (m *Metrics) Snapshot() Snapshot {
	if m == nil {
		return Snapshot{}
	}
	processing := make([]string, 0)
	m.requests.Range(func(key, _ any) bool {
		if id, ok := key.(string); ok {
			processing = append(processing, id)
		}
		return true
	})
	return Snapshot{
		InFlight:   m.inFlight.Load(),
		Processed:  m.processed.Load(),
		Failed:     m.failed.Load(),
		StartedAt:  m.started,
		Processing: processing,
	}
}

func (m *Metrics) TrackStart(id string) {
	if m == nil || id == "" {
		return
	}
	m.requests.Store(id, struct{}{})
}

func (m *Metrics) TrackDone(id string) {
	if m == nil || id == "" {
		return
	}
	m.requests.Delete(id)
}
