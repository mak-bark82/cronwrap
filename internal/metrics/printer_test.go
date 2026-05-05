package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/metrics"
)

func TestPrintSummary_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	s := metrics.Summary{
		Total:       5,
		Succeeded:   4,
		Failed:      1,
		AvgDuration: 3 * time.Second,
	}
	metrics.PrintSummary(&buf, s)
	out := buf.String()

	for _, want := range []string{"5", "4", "1", "3s", "METRIC", "VALUE"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestPrintResults_Empty(t *testing.T) {
	var buf bytes.Buffer
	metrics.PrintResults(&buf, nil)
	if !strings.Contains(buf.String(), "no results") {
		t.Errorf("expected 'no results' message, got: %s", buf.String())
	}
}

func TestPrintResults_ShowsJobData(t *testing.T) {
	var buf bytes.Buffer
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	results := []metrics.JobResult{
		{
			JobName:   "backup",
			Success:   true,
			Attempts:  2,
			Duration:  5 * time.Second,
			StartedAt: now,
		},
		{
			JobName:   "cleanup",
			Success:   false,
			Attempts:  3,
			Duration:  1 * time.Second,
			StartedAt: now,
		},
	}
	metrics.PrintResults(&buf, results)
	out := buf.String()

	for _, want := range []string{"backup", "ok", "cleanup", "fail", "2024-06-01"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
