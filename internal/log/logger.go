package log

import (
	"log"
	"log/slog"

	"github.com/fatih/color"
	"github.com/zmjung/jamesdb/config"
)

func getLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // Default to info if unknown level
	}
}

func getLoggerFormat(format string, opt *slog.HandlerOptions) *slog.Logger {
	switch format {
	case "text":
		return slog.New(slog.NewTextHandler(log.Writer(), opt))
	case "json":
		return slog.New(slog.NewJSONHandler(log.Writer(), opt))
	case "custom":
		color.NoColor = false
		return slog.New(&CustomJSONHandler{
			JSONHandler: *slog.NewJSONHandler(log.Writer(), opt),
			w:           log.Writer(),
			opt:         opt,
		})
	case "simple":
		fallthrough
	default:
		return slog.New(&SimpleJSONHandler{
			JSONHandler: *slog.NewJSONHandler(log.Writer(), opt),
		})
	}
}

func SetDefaultLogger(cfg *config.Config) {
	opt := &slog.HandlerOptions{
		Level:     getLogLevel(cfg.Logging.Level),
		AddSource: cfg.Logging.Source,
	}

	logger := getLoggerFormat(cfg.Logging.Format, opt)
	slog.SetDefault(logger)
}
