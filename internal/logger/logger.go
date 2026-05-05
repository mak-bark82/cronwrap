package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a log message.
type Level int

const (
	LevelInfo  Level = iota
	LevelWarn
	LevelError
)

var levelNames = map[Level]string{
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

// Logger writes structured log lines for cronwrap job execution.
type Logger struct {
	out    io.Writer
	jobName string
}

// New creates a Logger that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(jobName string, w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w, jobName: jobName}
}

// Info logs an informational message.
func (l *Logger) Info(msg string) {
	l.log(LevelInfo, msg)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string) {
	l.log(LevelWarn, msg)
}

// Error logs an error message.
func (l *Logger) Error(msg string) {
	l.log(LevelError, msg)
}

// JobStarted logs the beginning of a job execution.
func (l *Logger) JobStarted(attempt int) {
	l.Info(fmt.Sprintf("job started (attempt %d)", attempt))
}

// JobSucceeded logs a successful job completion with its duration.
func (l *Logger) JobSucceeded(attempt int, duration time.Duration) {
	l.Info(fmt.Sprintf("job succeeded (attempt %d) duration=%s", attempt, duration.Round(time.Millisecond)))
}

// JobFailed logs a failed job attempt with the error and duration.
func (l *Logger) JobFailed(attempt int, duration time.Duration, err error) {
	l.Warn(fmt.Sprintf("job failed (attempt %d) duration=%s error=%q", attempt, duration.Round(time.Millisecond), err))
}

// JobExhausted logs that all retry attempts have been exhausted.
func (l *Logger) JobExhausted(totalAttempts int, lastErr error) {
	l.Error(fmt.Sprintf("job exhausted all %d attempt(s): %v", totalAttempts, lastErr))
}

func (l *Logger) log(level Level, msg string) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(l.out, "%s [%s] job=%s %s\n", timestamp, levelNames[level], l.jobName, msg)
}
