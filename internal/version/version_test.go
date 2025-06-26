package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	// Save original value
	origVersion := Version
	defer func() { Version = origVersion }()

	// Test default value
	Version = "dev"
	if got := GetVersion(); got != "dev" {
		t.Errorf("GetVersion() = %q, want %q", got, "dev")
	}

	// Test custom value
	Version = "1.0.0"
	if got := GetVersion(); got != "1.0.0" {
		t.Errorf("GetVersion() = %q, want %q", got, "1.0.0")
	}
}

func TestGetFullVersion(t *testing.T) {
	// Save original values
	origVersion := Version
	origCommit := Commit
	origDate := Date
	defer func() {
		Version = origVersion
		Commit = origCommit
		Date = origDate
	}()

	// Test with default values
	Version = "dev"
	Commit = "none"
	Date = "unknown"

	expected := "dev (commit: none, built: unknown)"
	if got := GetFullVersion(); got != expected {
		t.Errorf("GetFullVersion() = %q, want %q", got, expected)
	}

	// Test with custom values
	Version = "1.0.0"
	Commit = "abc123"
	Date = "2025-01-26T12:00:00Z"

	expected = "1.0.0 (commit: abc123, built: 2025-01-26T12:00:00Z)"
	if got := GetFullVersion(); got != expected {
		t.Errorf("GetFullVersion() = %q, want %q", got, expected)
	}
}

func TestGetUserAgent(t *testing.T) {
	// Save original value
	origVersion := Version
	defer func() { Version = origVersion }()

	Version = "1.0.0"
	userAgent := GetUserAgent()

	// Check that it contains the expected parts
	if !strings.HasPrefix(userAgent, "payjp-cli/1.0.0") {
		t.Errorf("GetUserAgent() should start with 'payjp-cli/1.0.0', got %q", userAgent)
	}

	if !strings.Contains(userAgent, "go/"+runtime.Version()) {
		t.Errorf("GetUserAgent() should contain Go version, got %q", userAgent)
	}
}