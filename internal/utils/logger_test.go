package utils

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestInitLogger_DebugTrue(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize logger with debug=true
	InitLogger(true)

	// Write a test log message
	Logger.Debug().Msg("test message")

	// Close write end and read output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify debug mode is enabled
	if Logger.GetLevel() != zerolog.DebugLevel {
		t.Errorf("Expected debug level, got %v", Logger.GetLevel())
	}
	// Verify output contains the message
	if !strings.Contains(output, "test message") {
		t.Errorf("Output should contain 'test message', got: %s", output)
	}
	// Verify output format is raw text (not JSON)
	if strings.Contains(output, "\"level\"") {
		t.Error("Output should not be JSON format")
	}
	// Verify no timestamp in output
	if strings.Contains(output, "timestamp") {
		t.Error("Output should not contain timestamp")
	}
}

func TestInitLogger_DebugFalse(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize logger with debug=false
	InitLogger(false)

	// Write a test log message
	Logger.Debug().Msg("test message")

	// Close write end and read output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify logger is disabled (no output)
	if Logger.GetLevel() != zerolog.Disabled {
		t.Errorf("Expected disabled level, got %v", Logger.GetLevel())
	}
	// Verify no output
	if output != "" {
		t.Errorf("Logger should be silent when debug=false, got: %s", output)
	}
}

func TestInitLogger_DebugTrueVerboseTrue(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize logger with debug=true, verbose=true
	InitLogger(true)

	// Write a test log message
	Logger.Debug().Msg("test message")

	// Close write end and read output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify debug mode is enabled (debug takes precedence)
	if Logger.GetLevel() != zerolog.DebugLevel {
		t.Errorf("Expected debug level, got %v", Logger.GetLevel())
	}
	// Verify output contains the message
	if !strings.Contains(output, "test message") {
		t.Errorf("Output should contain 'test message', got: %s", output)
	}
}

func TestInitLogger_DebugFalseVerboseTrue(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize logger with debug=false, verbose=true
	InitLogger(false)

	// Write a test log message
	Logger.Debug().Msg("test message")

	// Close write end and read output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify logger is disabled (verbose is no-op without debug)
	if Logger.GetLevel() != zerolog.Disabled {
		t.Errorf("Expected disabled level, got %v", Logger.GetLevel())
	}
	// Verify no output
	if output != "" {
		t.Errorf("Logger should be silent when debug=false even with verbose=true, got: %s", output)
	}
}

func TestInitLogger_LogOutputFormat(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize logger with debug=true
	InitLogger(true)

	// Write a test log message with fields
	Logger.Debug().Str("key", "value").Str("key2", "value2").Msg("test message")

	// Close write end and read output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify output format is raw text (human-readable)
	if !strings.Contains(output, "[DEBUG]") {
		t.Errorf("Output should contain [DEBUG], got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Output should contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("Output should contain 'key=value', got: %s", output)
	}
	if !strings.Contains(output, "key2=value2") {
		t.Errorf("Output should contain 'key2=value2', got: %s", output)
	}
	// Verify not JSON format
	if strings.Contains(output, "\"level\"") {
		t.Error("Output should not be JSON format")
	}
	if strings.Contains(output, "\"message\"") {
		t.Error("Output should not be JSON format")
	}
	// Verify no timestamp
	if strings.Contains(output, "timestamp") {
		t.Error("Output should not contain timestamp")
	}
	if strings.Contains(output, "time=") {
		t.Error("Output should not contain time field")
	}
}

func TestInitLogger_LogSuppressionWhenDebugDisabled(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Initialize logger with debug=false
	InitLogger(false)

	// Write multiple test log messages
	Logger.Debug().Msg("message 1")
	Logger.Debug().Str("key", "value").Msg("message 2")
	Logger.Debug().Int("count", 42).Msg("message 3")

	// Close write end and read output
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify no output at all
	if output != "" {
		t.Errorf("Logger should suppress all output when debug=false, got: %s", output)
	}
	// Verify logger level is disabled
	if Logger.GetLevel() != zerolog.Disabled {
		t.Errorf("Expected disabled level, got %v", Logger.GetLevel())
	}
}
