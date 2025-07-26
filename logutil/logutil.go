package logutil

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

// customSLoggerHandler implements slog.Handler for custom formatting
// [time][LEVEL] message | key=value
type customSLoggerHandler struct {
	level slog.Level
}

func (h *customSLoggerHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level
}

func (h *customSLoggerHandler) Handle(_ context.Context, r slog.Record) error {
	var w io.Writer
	if r.Level >= slog.LevelError {
		w = os.Stderr
	} else {
		w = os.Stdout
	}
	ts := r.Time.Format("2006-01-02 15:04:05")
	msg := r.Message
	level := strings.ToUpper(r.Level.String())
	attrStr := ""
	r.Attrs(func(a slog.Attr) bool {
		attrStr += a.Key + "=" + attrValueToString(a.Value) + " "
		return true
	})
	if len(attrStr) > 0 {
		attrStr = strings.TrimSpace(attrStr)
		msg += " | " + attrStr
	}
	_, err := io.WriteString(w, "["+ts+"]["+level+"] "+msg+"\n")
	return err
}

// attrValueToString converts slog.Value to string for logging
func attrValueToString(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		return v.String()
	case slog.KindInt64:
		return fmt.Sprintf("%d", v.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%d", v.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%f", v.Float64())
	case slog.KindBool:
		return fmt.Sprintf("%t", v.Bool())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format(time.RFC3339)
	default:
		return v.String()
	}
}

func (h *customSLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *customSLoggerHandler) WithGroup(name string) slog.Handler {
	return h
}

func Setup(level string) {
	var lvl slog.Level
	switch strings.ToUpper(level) {
	case "DEBUG":
		lvl = slog.LevelDebug
	case "INFO":
		lvl = slog.LevelInfo
	case "WARN":
		lvl = slog.LevelWarn
	case "ERROR":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelDebug
	}
	h := &customSLoggerHandler{level: lvl}
	slog.SetDefault(slog.New(h))
}
