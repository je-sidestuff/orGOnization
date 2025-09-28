package templates

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPrintWeekdays_Monday(t *testing.T) {
	loc, err := time.LoadLocation("Atlantic/Reykjavik")
	if err != nil {
		t.Fatal(err)
	}

	// Define a monday date
	date := time.Date(2024, 3, 18, 0, 0, 0, 0, loc)

	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	// Replace os.Stdout with the write end of the pipe
	origStdout := os.Stdout
	os.Stdout = w

	// Call PrintWeekdays in a goroutine so that we can immediately read from the pipe
	go func() {
		defer func() {
			if err := w.Close(); err != nil {
				t.Error(err)
			}
		}()
		if err := PrintWeekdays(date); err != nil {
			t.Error(err)
		}
	}()

	// Read the output
	out, _ := io.ReadAll(r)

	// Restore os.Stdout
	os.Stdout = origStdout

	// Validate captured output, allowing for slight time difference formatting
	expectedOutputStart := "Monday, Mar 18th:"
	expectedOutputEnd := "Sunday, Mar 24th:\n"
	if !strings.HasPrefix(string(out), expectedOutputStart) || !strings.HasSuffix(string(out), expectedOutputEnd) {
		t.Errorf("PrintWeekdays(%v) output did not match expected format.\nExpected to start with: %v\nExpected to end with: %v\nReceived: %v", date, expectedOutputStart, expectedOutputEnd, string(out))
	}
}
