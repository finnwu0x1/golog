// Package golog is a small levelled logger for command-line Go programs.
//
// It writes colourised, printf-style lines to a single stream and reports the
// caller's file and line on every entry. It carries no redaction rules, no
// dependency-injection hooks and no knowledge of any particular application, so
// it can be dropped into any project as-is.
package golog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Colour escape sequences used for the level prefixes.
const (
	resetColor  = "\033[0m"
	redColor    = "\033[31m"
	greenColor  = "\033[32m"
	yellowColor = "\033[33m"
	blueColor   = "\033[34m"
)

// Logger writes levelled log lines to one output stream.
//
// The zero value is not usable: construct a Logger with NewLogger or
// NewLoggerWithWriter. A Logger is safe for concurrent use.
type Logger struct {
	mu     sync.Mutex
	logger *log.Logger
}

// NewLogger returns a Logger that writes to standard output.
func NewLogger() *Logger {
	return NewLoggerWithWriter(os.Stdout)
}

// NewLoggerWithWriter returns a Logger that writes to output.
func NewLoggerWithWriter(output io.Writer) *Logger {
	return &Logger{
		// Lshortfile appends the caller's file:line to every entry. print passes
		// the matching call depth so that location is the caller of the level
		// method rather than a line inside this package.
		logger: log.New(output, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs at informational level.
func (l *Logger) Info(format string, args ...any) {
	l.print(greenColor+"[I] "+resetColor, format, args...)
}

// Warn logs at warning level.
func (l *Logger) Warn(format string, args ...any) {
	l.print(yellowColor+"[W] "+resetColor, format, args...)
}

// Error logs at error level.
func (l *Logger) Error(format string, args ...any) {
	l.print(redColor+"[E] "+resetColor, format, args...)
}

// Debug logs at debug level.
func (l *Logger) Debug(format string, args ...any) {
	l.print(blueColor+"[D] "+resetColor, format, args...)
}

// Fatal logs at fatal level and terminates the process with exit status 1.
func (l *Logger) Fatal(format string, args ...any) {
	l.print(redColor+"[F] "+resetColor, format, args...)
	os.Exit(1)
}

// print serialises the write, applies the level prefix and forwards the
// formatted message. The call depth is 3 because the chain is
// print -> level method -> caller. A failing writer is reported nowhere else,
// matching the standard library logger, which also drops write errors.
func (l *Logger) print(prefix string, format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.SetPrefix(prefix)
	_ = l.logger.Output(3, fmt.Sprintf(format, args...))
}
