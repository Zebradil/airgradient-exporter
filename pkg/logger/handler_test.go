package logger

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewColoredTextHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)
	if handler == nil {
		t.Fatal("NewColoredTextHandler() returned nil")
	}

	handler = NewColoredTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	if handler == nil {
		t.Fatal("NewColoredTextHandler() returned nil with options")
	}
}

func TestColoredTextHandler_Enabled(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)

	ctx := context.Background()
	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("Enabled() should return true for Info level")
	}
	if !handler.Enabled(ctx, slog.LevelError) {
		t.Error("Enabled() should return true for Error level")
	}
	if handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("Enabled() should return false for Debug level (default is Info)")
	}

	// Test with custom level
	handler = NewColoredTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	if !handler.Enabled(ctx, slog.LevelDebug) {
		t.Error("Enabled() should return true for Debug level when set")
	}
}

func TestColoredTextHandler_Handle(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)

	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.AddAttrs(slog.String("key", "value"), slog.Int("num", 42))

	err := handler.Handle(ctx, record)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("output should contain message, got: %s", output)
	}
	if !strings.Contains(output, "key=") {
		t.Errorf("output should contain key, got: %s", output)
	}
	if !strings.Contains(output, "value") {
		t.Errorf("output should contain value, got: %s", output)
	}
	if !strings.Contains(output, "\033[") {
		t.Error("output should contain ANSI color codes")
	}
}

func TestColoredTextHandler_Handle_AllLevels(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)

	ctx := context.Background()
	levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}

	for _, level := range levels {
		buf.Reset()
		record := slog.NewRecord(time.Now(), level, "test", 0)
		err := handler.Handle(ctx, record)
		if err != nil {
			t.Errorf("Handle() error for level %v = %v", level, err)
		}
		output := buf.String()
		if !strings.Contains(output, "test") {
			t.Errorf("output should contain message for level %v", level)
		}
	}
}

func TestColoredTextHandler_Handle_WithAttributes(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)

	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	record.AddAttrs(
		slog.String("str", "value"),
		slog.Int("int", 42),
		slog.Float64("float", 3.14),
		slog.Bool("bool", true),
		slog.Duration("duration", time.Second),
	)

	err := handler.Handle(ctx, record)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "str=") || !strings.Contains(output, `"value"`) {
		t.Error("output should contain string attribute")
	}
	if !strings.Contains(output, "int=42") {
		t.Error("output should contain int attribute")
	}
	if !strings.Contains(output, "float=") {
		t.Error("output should contain float attribute")
	}
	if !strings.Contains(output, "bool=true") {
		t.Error("output should contain bool attribute")
	}
	if !strings.Contains(output, "duration=1s") {
		t.Error("output should contain duration attribute")
	}
}

func TestColoredTextHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)
	handler = handler.WithAttrs([]slog.Attr{slog.String("app", "test")}).(*ColoredTextHandler)

	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	err := handler.Handle(ctx, record)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "app=") || !strings.Contains(output, `"test"`) {
		t.Error("output should contain handler attribute")
	}
}

func TestColoredTextHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)
	handler = handler.WithGroup("group").(*ColoredTextHandler)

	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	record.AddAttrs(slog.String("key", "value"))
	err := handler.Handle(ctx, record)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "key=") {
		t.Error("output should contain grouped attribute")
	}
}

func TestColoredTextHandler_formatLevel(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)

	levels := map[slog.Level]string{
		slog.LevelDebug: "DBG",
		slog.LevelInfo:  "INF",
		slog.LevelWarn:  "WRN",
		slog.LevelError: "ERR",
	}

	for level, expected := range levels {
		formatted := handler.formatLevel(level)
		if formatted != expected {
			t.Errorf("formatLevel(%v) = %s, want %s", level, formatted, expected)
		}
	}
}

func TestColoredTextHandler_getLevelColor(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)

	colors := map[slog.Level]string{
		slog.LevelDebug: "\033[36m", // Cyan
		slog.LevelInfo:  "\033[32m", // Green
		slog.LevelWarn:  "\033[33m", // Yellow
		slog.LevelError: "\033[31m", // Red
	}

	for level, expected := range colors {
		color := handler.getLevelColor(level)
		if color != expected {
			t.Errorf("getLevelColor(%v) = %s, want %s", level, color, expected)
		}
	}
}

func TestNewPlainTextHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := NewPlainTextHandler(&buf, nil)
	if handler == nil {
		t.Fatal("NewPlainTextHandler() returned nil")
	}
}

func TestPlainTextHandler_Handle(t *testing.T) {
	var buf bytes.Buffer
	handler := NewPlainTextHandler(&buf, nil)

	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.AddAttrs(slog.String("key", "value"))

	err := handler.Handle(ctx, record)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("output should contain message, got: %s", output)
	}
	if strings.Contains(output, "\033[") {
		t.Error("output should not contain ANSI color codes")
	}
}

func TestNewUnbufferedWriter(t *testing.T) {
	writer := newUnbufferedWriter(os.Stdout)
	if writer == nil {
		t.Fatal("newUnbufferedWriter() returned nil")
	}

	// Test with non-file writer
	var buf bytes.Buffer
	writer = newUnbufferedWriter(&buf)
	if writer == nil {
		t.Fatal("newUnbufferedWriter() returned nil for bytes.Buffer")
	}
}

func TestUnbufferedWriter_Write(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Errorf("Failed to remove temp file: %v", err)
		}
	}()
	defer func() {
		if err := tmpfile.Close(); err != nil {
			t.Errorf("Failed to close temp file: %v", err)
		}
	}()

	writer := newUnbufferedWriter(tmpfile)
	if writer == nil {
		t.Fatal("newUnbufferedWriter() returned nil")
	}

	data := []byte("test data\n")
	n, err := writer.Write(data)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != len(data) {
		t.Errorf("Write() wrote %d bytes, want %d", n, len(data))
	}
}

func TestNewLogger_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("text", &buf)
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	logger.Info("test message", "key", "value")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("output should contain message, got: %s", output)
	}
}

func TestNewLogger_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("json", &buf)
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	logger.Info("test message", "key", "value")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("output should contain message, got: %s", output)
	}
	if !strings.Contains(output, "\"msg\"") {
		t.Error("output should be JSON format")
	}
}

func TestNewLogger_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("", &buf)
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	logger.Info("test message")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("output should contain message, got: %s", output)
	}
}

func TestNewLogger_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("invalid", &buf)
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}

	logger.Info("test message")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("output should contain message, got: %s", output)
	}
}

func TestColoredTextHandler_appendValue_AllTypes(t *testing.T) {
	var buf bytes.Buffer
	handler := NewColoredTextHandler(&buf, nil)

	testCases := []struct {
		name  string
		value slog.Value
		check func(output string) bool
	}{
		{
			name:  "String",
			value: slog.StringValue("test"),
			check: func(output string) bool { return strings.Contains(output, `"test"`) },
		},
		{
			name:  "Int64",
			value: slog.Int64Value(42),
			check: func(output string) bool { return strings.Contains(output, "42") },
		},
		{
			name:  "Uint64",
			value: slog.Uint64Value(100),
			check: func(output string) bool { return strings.Contains(output, "100") },
		},
		{
			name:  "Float64",
			value: slog.Float64Value(3.14),
			check: func(output string) bool { return strings.Contains(output, "3.14") },
		},
		{
			name:  "Bool",
			value: slog.BoolValue(true),
			check: func(output string) bool { return strings.Contains(output, "true") },
		},
		{
			name:  "Duration",
			value: slog.DurationValue(time.Second),
			check: func(output string) bool { return strings.Contains(output, "1s") },
		},
		{
			name:  "Time",
			value: slog.TimeValue(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)),
			check: func(output string) bool { return strings.Contains(output, "2024-01-01T12:00:00") },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			ctx := context.Background()
			record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
			record.AddAttrs(slog.Attr{Key: "key", Value: tc.value})
			err := handler.Handle(ctx, record)
			if err != nil {
				t.Fatalf("Handle() error = %v", err)
			}
			if !tc.check(buf.String()) {
				t.Errorf("output should contain expected value for %s, got: %s", tc.name, buf.String())
			}
		})
	}
}
