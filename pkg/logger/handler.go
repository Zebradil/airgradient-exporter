package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"golang.org/x/term"
)

// ColoredTextHandler formats logs with custom formatting: time, colored level, bright msg, and attributes
type ColoredTextHandler struct {
	w     io.Writer
	opts  *slog.HandlerOptions
	attrs []slog.Attr
	group string
}

func NewColoredTextHandler(w io.Writer, opts *slog.HandlerOptions) *ColoredTextHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColoredTextHandler{
		w:    w,
		opts: opts,
	}
}

func (h *ColoredTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *ColoredTextHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)

	// Format time without timezone
	if !r.Time.IsZero() {
		buf = append(buf, r.Time.Format("2006-01-02T15:04:05.000")...)
		buf = append(buf, ' ')
	}

	// Format level with color (3-character abbreviation)
	levelColor := h.getLevelColor(r.Level)
	levelStr := h.formatLevel(r.Level)
	buf = append(buf, levelColor...)
	buf = append(buf, levelStr...)
	buf = append(buf, "\033[0m"...)
	buf = append(buf, ' ')

	// Format message with brightness
	brightCode := "\033[1m"
	resetCode := "\033[0m"
	buf = append(buf, brightCode...)
	buf = append(buf, r.Message...)
	buf = append(buf, resetCode...)

	// Add attributes
	r.Attrs(func(a slog.Attr) bool {
		buf = append(buf, ' ')
		buf = h.appendAttr(buf, a)
		return true
	})

	// Add handler attributes
	for _, a := range h.attrs {
		buf = append(buf, ' ')
		buf = h.appendAttr(buf, a)
	}

	buf = append(buf, '\n')

	_, err := h.w.Write(buf)
	return err
}

func (h *ColoredTextHandler) appendAttr(buf []byte, a slog.Attr) []byte {
	if a.Equal(slog.Attr{}) {
		return buf
	}

	// Format as key=value
	buf = append(buf, a.Key...)
	buf = append(buf, '=')
	buf = h.appendValue(buf, a.Value)
	return buf
}

func (h *ColoredTextHandler) appendValue(buf []byte, v slog.Value) []byte {
	switch v.Kind() {
	case slog.KindString:
		buf = append(buf, fmt.Sprintf("%q", v.String())...)
	case slog.KindInt64:
		buf = append(buf, fmt.Sprintf("%d", v.Int64())...)
	case slog.KindUint64:
		buf = append(buf, fmt.Sprintf("%d", v.Uint64())...)
	case slog.KindFloat64:
		buf = append(buf, fmt.Sprintf("%g", v.Float64())...)
	case slog.KindBool:
		buf = append(buf, fmt.Sprintf("%t", v.Bool())...)
	case slog.KindDuration:
		buf = append(buf, v.Duration().String()...)
	case slog.KindTime:
		buf = append(buf, v.Time().Format("2006-01-02T15:04:05.000")...)
	case slog.KindAny:
		buf = append(buf, fmt.Sprintf("%v", v.Any())...)
	case slog.KindLogValuer:
		return h.appendValue(buf, v.LogValuer().LogValue())
	case slog.KindGroup:
		// For groups, format as key.subkey=value
		groupAttrs := v.Group()
		for i, ga := range groupAttrs {
			if i > 0 {
				buf = append(buf, ' ')
			}
			buf = h.appendAttr(buf, ga)
		}
	}
	return buf
}

func (h *ColoredTextHandler) formatLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "DBG"
	case slog.LevelInfo:
		return "INF"
	case slog.LevelWarn:
		return "WRN"
	case slog.LevelError:
		return "ERR"
	default:
		return level.String()
	}
}

func (h *ColoredTextHandler) getLevelColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "\033[36m" // Cyan
	case slog.LevelInfo:
		return "\033[32m" // Green
	case slog.LevelWarn:
		return "\033[33m" // Yellow
	case slog.LevelError:
		return "\033[31m" // Red
	default:
		return "\033[0m" // Reset
	}
}

func (h *ColoredTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColoredTextHandler{
		w:     h.w,
		opts:  h.opts,
		attrs: append(h.attrs, attrs...),
		group: h.group,
	}
}

func (h *ColoredTextHandler) WithGroup(name string) slog.Handler {
	return &ColoredTextHandler{
		w:     h.w,
		opts:  h.opts,
		attrs: h.attrs,
		group: h.group + name + ".",
	}
}

// PlainTextHandler formats logs with the same format as ColoredTextHandler but without colors
type PlainTextHandler struct {
	*ColoredTextHandler
}

func NewPlainTextHandler(w io.Writer, opts *slog.HandlerOptions) *PlainTextHandler {
	return &PlainTextHandler{
		ColoredTextHandler: &ColoredTextHandler{
			w:    w,
			opts: opts,
		},
	}
}

func (h *PlainTextHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)

	// Format time without timezone
	if !r.Time.IsZero() {
		buf = append(buf, r.Time.Format("2006-01-02T15:04:05.000")...)
		buf = append(buf, ' ')
	}

	// Format level without color (3-character abbreviation)
	levelStr := h.formatLevel(r.Level)
	buf = append(buf, levelStr...)
	buf = append(buf, ' ')

	// Format message without brightness
	buf = append(buf, r.Message...)

	// Add attributes
	r.Attrs(func(a slog.Attr) bool {
		buf = append(buf, ' ')
		buf = h.appendAttr(buf, a)
		return true
	})

	// Add handler attributes
	for _, a := range h.attrs {
		buf = append(buf, ' ')
		buf = h.appendAttr(buf, a)
	}

	buf = append(buf, '\n')

	_, err := h.w.Write(buf)
	return err
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
				// Use plain text handler with same format for non-terminal
				handler = NewPlainTextHandler(output, nil)
			}
		} else {
			// For non-file writers, use plain text with same format
			handler = NewPlainTextHandler(output, nil)
		}
	default:
		// Default to text format
		if file, ok := output.(*os.File); ok {
			if term.IsTerminal(int(file.Fd())) {
				handler = NewColoredTextHandler(output, nil)
			} else {
				handler = NewPlainTextHandler(output, nil)
			}
		} else {
			handler = NewPlainTextHandler(output, nil)
		}
	}

	return slog.New(handler)
}
