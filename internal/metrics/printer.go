package metrics

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintSummary writes a human-readable summary table to w.
func PrintSummary(w io.Writer, s Summary) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "METRIC\tVALUE")
	fmt.Fprintf(tw, "Total runs\t%d\n", s.Total)
	fmt.Fprintf(tw, "Succeeded\t%d\n", s.Succeeded)
	fmt.Fprintf(tw, "Failed\t%d\n", s.Failed)
	fmt.Fprintf(tw, "Avg duration\t%s\n", s.AvgDuration)
	_ = tw.Flush()
}

// PrintResults writes per-run details to w.
func PrintResults(w io.Writer, results []JobResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no results recorded")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "JOB\tSTATUS\tATTEMPTS\tDURATION\tSTARTED")
	for _, r := range results {
		status := "ok"
		if !r.Success {
			status = "fail"
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n",
			r.JobName,
			status,
			r.Attempts,
			r.Duration,
			r.StartedAt.Format("2006-01-02 15:04:05"),
		)
	}
	_ = tw.Flush()
}
