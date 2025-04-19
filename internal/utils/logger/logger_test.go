package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestEventLogger_AddEventAndLogEvents(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	logger := New(nil, Cluster, "create", "123")
	logger.WithFields(LogFields{"foo": "bar"}).WithLevel(LevelInfo)
	logger.AddEvent(Success)
	logger.LogEvents()

	output := buf.String()
	if !strings.Contains(output, "Cluster") || !strings.Contains(output, "foo=bar") {
		t.Errorf("log output missing expected content: %q", output)
	}
}

func TestEventLogger_ErrorEvent(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	logger := New(nil, Cluster, "fail", "errid")
	logger.WithLevel(LevelWarning)
	errMsg := "something went wrong"
	logger.AddEvent(Failure, errMsg)
	logger.LogEvents()

	output := buf.String()
	if !strings.Contains(output, "Failure") || !strings.Contains(output, errMsg) {
		t.Errorf("log output missing error content: %q", output)
	}
}

func TestFormatLogMessage_EdgeCases(t *testing.T) {
	fields := LogFields{
		"resource":  Cluster,
		"action":    "act",
		"status":    Success,
		"level":     string(LevelInfo),
		"timestamp": "2025-01-01T00:00:00Z",
	}
	msg := formatLogMessage(fields)
	if !strings.Contains(msg, "Cluster act") {
		t.Errorf("formatLogMessage missing resource/action: %q", msg)
	}
}

func TestGetActionString(t *testing.T) {
	if s, ok := getActionString(LogCrudMethod("foo")); !ok || s != "foo" {
		t.Errorf("getActionString(LogCrudMethod) failed: %v, %v", s, ok)
	}
	if s, ok := getActionString("bar"); !ok || s != "bar" {
		t.Errorf("getActionString(string) failed: %v, %v", s, ok)
	}
	if s, ok := getActionString(123); ok || s != "" {
		t.Errorf("getActionString(non-string) should fail, got: %v, %v", s, ok)
	}
}
