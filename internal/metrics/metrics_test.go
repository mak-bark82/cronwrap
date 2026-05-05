package metrics_test

import (
	"errors"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/metrics"
)

func TestCollector_RecordAndAll(t *testing.T) {
	c := metrics.NewCollector()

	r := metrics.JobResult{
		JobName:  "test-job",
		Success:  true,
		Attempts: 1,
		Duration: 2 * time.Second,
	}
	c.Record(r)

	all := c.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 result, got %d", len(all))
	}
	if all[0].JobName != "test-job" {
		t.Errorf("expected job name 'test-job', got %q", all[0].JobName)
	}
}

func TestCollector_Summary_AllSucceeded(t *testing.T) {
	c := metrics.NewCollector()
	c.Record(metrics.JobResult{Success: true, Duration: 1 * time.Second})
	c.Record(metrics.JobResult{Success: true, Duration: 3 * time.Second})

	s := c.Summary()
	if s.Total != 2 {
		t.Errorf("expected Total=2, got %d", s.Total)
	}
	if s.Succeeded != 2 {
		t.Errorf("expected Succeeded=2, got %d", s.Succeeded)
	}
	if s.Failed != 0 {
		t.Errorf("expected Failed=0, got %d", s.Failed)
	}
	if s.AvgDuration != 2*time.Second {
		t.Errorf("expected AvgDuration=2s, got %v", s.AvgDuration)
	}
}

func TestCollector_Summary_Mixed(t *testing.T) {
	c := metrics.NewCollector()
	c.Record(metrics.JobResult{Success: true, Duration: 4 * time.Second})
	c.Record(metrics.JobResult{Success: false, Duration: 2 * time.Second, Error: errors.New("boom")})

	s := c.Summary()
	if s.Succeeded != 1 || s.Failed != 1 {
		t.Errorf("unexpected succeeded/failed counts: %d/%d", s.Succeeded, s.Failed)
	}
	if s.AvgDuration != 3*time.Second {
		t.Errorf("expected AvgDuration=3s, got %v", s.AvgDuration)
	}
}

func TestCollector_Summary_Empty(t *testing.T) {
	c := metrics.NewCollector()
	s := c.Summary()
	if s.Total != 0 || s.AvgDuration != 0 {
		t.Errorf("expected zero summary for empty collector, got %+v", s)
	}
}

func TestCollector_AllReturnsCopy(t *testing.T) {
	c := metrics.NewCollector()
	c.Record(metrics.JobResult{JobName: "job-a"})

	all := c.All()
	all[0].JobName = "mutated"

	original := c.All()
	if original[0].JobName == "mutated" {
		t.Error("All() should return a copy, not a reference")
	}
}
