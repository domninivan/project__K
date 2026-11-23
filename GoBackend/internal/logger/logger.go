package logger

import (
	"fmt"
	"io"
	"log/slog"
)

func New(level, format string, w io.Writer) (*slog.Logger, error) {
	var lvl slog.LevelVar
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("logger: invalid level %q: %w", level, err)
	}

	opts := &slog.HandlerOptions{
		Level: &lvl,
	}

	var h slog.Handler
	switch format {
	case "text":
		h = slog.NewTextHandler(w, opts)
	case "json":
		h = slog.NewJSONHandler(w, opts)
	default:
		return nil, fmt.Errorf("logger: unknown format %q (want \"text\" or \"json\")", format)
	}
	return slog.New(h), nil
}
