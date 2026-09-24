package golog

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func TestLevelsWriteTheirPrefixAndMessage(t *testing.T) {
	tests := []struct {
		name   string
		log    func(*Logger, string, ...any)
		prefix string
	}{
		{name: "info", log: (*Logger).Info, prefix: "[I]"},
		{name: "warn", log: (*Logger).Warn, prefix: "[W]"},
		{name: "error", log: (*Logger).Error, prefix: "[E]"},
		{name: "debug", log: (*Logger).Debug, prefix: "[D]"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			logger := NewLoggerWithWriter(&output)
			test.log(logger, "count=%d", 7)

			message := output.String()
			if !strings.Contains(message, test.prefix) {
				t.Fatalf("level prefix %q missing: %q", test.prefix, message)
			}
			if !strings.Contains(message, "count=7") {
				t.Fatalf("formatted arguments missing: %q", message)
			}
		})
	}
}

// TestLoggerReportsCallerFileAndLine pins the diagnostic context this package
// exists to restore: every entry must point at the line that logged it, not at
// a line inside golog.
func TestLoggerReportsCallerFileAndLine(t *testing.T) {
	var output bytes.Buffer
	logger := NewLoggerWithWriter(&output)

	// runtime.Caller and the logging call must stay adjacent: the expected
	// location is the line immediately after this one.
	_, _, line, ok := runtime.Caller(0)
	logger.Info("caller check")
	if !ok {
		t.Fatal("cannot resolve the caller of runtime.Caller")
	}

	want := fmt.Sprintf("logger_test.go:%d", line+1)
	if !strings.Contains(output.String(), want) {
		t.Fatalf("caller location %q missing: %q", want, output.String())
	}
}

func TestLoggerPreservesErrorMessages(t *testing.T) {
	var output bytes.Buffer
	logger := NewLoggerWithWriter(&output)
	logger.Error("request failed: %v", fmt.Errorf("dial tcp 10.0.0.5:9091: connection refused"))

	message := output.String()
	if !strings.Contains(message, "connection refused") {
		t.Fatalf("error message was dropped: %q", message)
	}
	if !strings.Contains(message, "10.0.0.5:9091") {
		t.Fatalf("error context was dropped: %q", message)
	}
}

func TestLoggerIsSafeForConcurrentUse(t *testing.T) {
	var output bytes.Buffer
	logger := NewLoggerWithWriter(&output)

	const writers = 16
	const perWriter = 32
	var group sync.WaitGroup
	for index := 0; index < writers; index++ {
		group.Add(1)
		go func(index int) {
			defer group.Done()
			for entry := 0; entry < perWriter; entry++ {
				logger.Info("writer=%d entry=%d", index, entry)
			}
		}(index)
	}
	group.Wait()

	lines := strings.Split(strings.TrimRight(output.String(), "\n"), "\n")
	if len(lines) != writers*perWriter {
		t.Fatalf("expected %d log lines, got %d", writers*perWriter, len(lines))
	}
	for _, line := range lines {
		if !strings.Contains(line, "writer=") || !strings.Contains(line, "entry=") {
			t.Fatalf("interleaved log line: %q", line)
		}
	}
}

func TestNewLoggerWritesToStandardOutput(t *testing.T) {
	if NewLogger().logger.Writer() != os.Stdout {
		t.Fatal("NewLogger no longer writes to stdout")
	}
}
