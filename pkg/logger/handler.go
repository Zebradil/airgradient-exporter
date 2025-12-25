package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"golang.org/x/term"
)

// ColoredTextHandler wraps a text handler and adds ANSI color codes
type ColoredTextHandler struct {
	handler slog.Handler
	w       io.Writer
}

func NewColoredTextHandler(w io.Writer, opts *slog.HandlerOptions) *ColoredTextHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColoredTextHandler{
		handler: slog.NewTextHandler(w, opts),
		w:       w,
	}
}

func (h *ColoredTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ColoredTextHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add color codes based on level
	var colorCode string
	switch r.Level {
	case slog.LevelDebug:
		colorCode = "\033[36m" // Cyan
	case slog.LevelInfo:
		colorCode = "\033[32m" // Green
	case slog.LevelWarn:
		colorCode = "\033[33m" // Yellow
	case slog.LevelError:
		colorCode = "\033[31m" // Red
	default:
		colorCode = "\033[0m" // Reset
	}
	resetCode := "\033[0m"

	// Write color code before the record
	if _, err := h.w.Write([]byte(colorCode)); err != nil {
		return err
	}

	// Handle the record with the underlying handler
	if err := h.handler.Handle(ctx, r); err != nil {
		return err
	}

	// Write reset code after the record
	_, err := h.w.Write([]byte(resetCode))
	return err
}

func (h *ColoredTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColoredTextHandler{
		handler: h.handler.WithAttrs(attrs),
		w:       h.w,
	}
}

func (h *ColoredTextHandler) WithGroup(name string) slog.Handler {
	return &ColoredTextHandler{
		handler: h.handler.WithGroup(name),
		w:       h.w,
	}
}

// NewLogger creates a logger based on the format and output
func NewLogger(format string, output io.Writer) *slog.Logger {
	var handler slog.Handler

	switch format {
	case "json":
		handler = slog.NewJSONHandler(output, nil)
	case "text":
		// Check if output is a terminal
		if file, ok := output.(*os.File); ok {
			if term.IsTerminal(int(file.Fd())) {
				// Use colored handler for terminal
				handler = NewColoredTextHandler(output, nil)
			} else {
				// Use plain text handler for non-terminal
				handler = slog.NewTextHandler(output, nil)
			}
		} else {
			// For non-file writers, use plain text
			handler = slog.NewTextHandler(output, nil)
		}
	default:
		// Default to text format
		if file, ok := output.(*os.File); ok {
			if term.IsTerminal(int(file.Fd())) {
				handler = NewColoredTextHandler(output, nil)
			} else {
				handler = slog.NewTextHandler(output, nil)
			}
		} else {
			handler = slog.NewTextHandler(output, nil)
		}
	}

	return slog.New(handler)
}
