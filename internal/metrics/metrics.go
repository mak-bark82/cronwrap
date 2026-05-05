package metrics

import (
	"sync"
	"time"
)

// JobResult represents the outcome of a single job execution.
type JobResult struct {
	JobName   string
	Success   bool
	Attempts  int
	Duration  time.Duration
	StartedAt time.Time
	Error     error
}

// Collector accumulates job execution metrics in memory.
type Collector struct {
	mu      sync.Mutex
	results []JobResult
}

// NewCollector creates a new Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// Record stores a job result.
func (c *Collector) Record(r JobResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.results = append(c.results, r)
}

// All returns a copy of all recorded results.
func (c *Collector) All() []JobResult {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]JobResult, len(c.results))
	copy(out, c.results)
	return out
}

// Summary returns aggregate statistics across all recorded results.
func (c *Collector) Summary() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()

	s := Summary{Total: len(c.results)}
	var totalDuration time.Duration
	for _, r := range c.results {
		if r.Success {
			s.Succeeded++
		} else {
			s.Failed++
		}
		totalDuration += r.Duration
	}
	if s.Total > 0 {
		s.AvgDuration = totalDuration / time.Duration(s.Total)
	}
	return s
}

// Summary holds aggregate metrics.
type Summary struct {
	Total       int
	Succeeded   int
	Failed      int
	AvgDuration time.Duration
}
