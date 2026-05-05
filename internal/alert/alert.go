package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Alert holds metadata about a job failure or warning event.
type Alert struct {
	JobName   string
	Level     Level
	Message   string
	Err       error
	OccuredAt time.Time
}

// Notifier is the interface that wraps the Notify method.
type Notifier interface {
	Notify(a Alert) error
}

// StderrNotifier writes alerts to stderr (or a configurable writer).
type StderrNotifier struct {
	Writer io.Writer
}

// NewStderrNotifier returns a StderrNotifier writing to os.Stderr.
func NewStderrNotifier() *StderrNotifier {
	return &StderrNotifier{Writer: os.Stderr}
}

// Notify formats and writes the alert to the configured writer.
func (n *StderrNotifier) Notify(a Alert) error {
	errStr := ""
	if a.Err != nil {
		errStr = fmt.Sprintf(" error=%q", a.Err.Error())
	}
	_, err := fmt.Fprintf(
		n.Writer,
		"[%s] level=%s job=%q message=%q%s\n",
		a.OccuredAt.UTC().Format(time.RFC3339),
		a.Level,
		a.JobName,
		a.Message,
		errStr,
	)
	return err
}

// New constructs an Alert with the current timestamp.
func New(jobName string, level Level, message string, err error) Alert {
	return Alert{
		JobName:   jobName,
		Level:     level,
		Message:   message,
		Err:       err,
		OccuredAt: time.Now(),
	}
}
