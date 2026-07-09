package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
)

func TestAddLogData(t *testing.T) {
	// Create a buffer to capture the log output
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	// Create a dummy HTTP request
	req := httptest.NewRequest("POST", "/test-path", nil)
	req.Header.Set("X-Test-Header", "test-value")
	req.RemoteAddr = "127.0.0.1"

	// Create an event and add log data
	event := logger.Log()
	addLogData(req, event, "test message")

	// Parse log output as JSON
	var logData map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logData); err != nil {
		t.Fatalf("Failed to parse log JSON: %v. Output: %s", err, buf.String())
	}

	// Verify logged values
	msgVal, exists := logData["message"]
	if !exists {
		msgVal, exists = logData["msg"]
	}
	if !exists || msgVal != "test message" {
		t.Errorf("Expected message 'test message', got log details: %v", logData)
	}

	if logData["method"] != "POST" {
		t.Errorf("Expected method 'POST', got: %v", logData["method"])
	}

	if logData["path"] != "/test-path" {
		t.Errorf("Expected path '/test-path', got: %v", logData["path"])
	}

	// Since r.Header headers are logged with loop indexing:
	// oname := fmt.Sprintf("%s-%d", name, i)
	if logData["X-Test-Header-0"] != "test-value" {
		t.Errorf("Expected X-Test-Header-0 'test-value', got: %v", logData["X-Test-Header-0"])
	}
}
