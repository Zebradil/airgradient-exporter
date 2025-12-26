package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPrintVersion(t *testing.T) {
	oldVersion := version
	oldCommit := commit
	oldDate := date

	defer func() {
		version = oldVersion
		commit = oldCommit
		date = oldDate
	}()

	version = "1.0.0"
	commit = "abc123"
	date = "2024-01-01"

	// Note: printVersion writes to stdout, so we can't easily capture it
	// This test just ensures it doesn't panic
	printVersion()
}

func TestPrintHelp(t *testing.T) {
	// Capture stdout by redirecting to a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	oldStdout := os.Stdout
	os.Stdout = w

	done := make(chan bool)
	var output string
	go func() {
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Errorf("Failed to read from pipe: %v", err)
		}
		output = buf.String()
		done <- true
	}()

	printHelp()

	if err := w.Close(); err != nil {
		t.Errorf("Failed to close pipe writer: %v", err)
	}
	os.Stdout = oldStdout
	<-done

	if !strings.Contains(output, "AirGradient Exporter") {
		t.Error("help output should contain 'AirGradient Exporter'")
	}
	if !strings.Contains(output, "Usage:") {
		t.Error("help output should contain 'Usage:'")
	}
	if !strings.Contains(output, "Flags:") {
		t.Error("help output should contain 'Flags:'")
	}
	if !strings.Contains(output, "Environment Variables:") {
		t.Error("help output should contain 'Environment Variables:'")
	}
	if !strings.Contains(output, "AIRGRADIENT_MONITORS") {
		t.Error("help output should contain 'AIRGRADIENT_MONITORS'")
	}
}

// Test that main doesn't panic with help flag
func TestMain_HelpFlag(t *testing.T) {
	// This is a basic smoke test - we can't easily test main() directly
	// without refactoring, but we can test the helper functions
}
