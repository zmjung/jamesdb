package log

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/internal/uuid"
)

type CustomJSONHandler struct {
	slog.JSONHandler
	w   io.Writer
	opt *slog.HandlerOptions
}

type RequestContextKey string

const requestKey RequestContextKey = "request"

func NewRequestContext(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, requestKey, value)
}

func colorize(level slog.Level, line string) string {
	switch level {
	case slog.LevelDebug:
		return color.MagentaString(line)
	case slog.LevelInfo:
		return color.BlueString(line)
	case slog.LevelWarn:
		return color.YellowString(line)
	case slog.LevelError:
		return color.RedString(line)
	}
	return color.WhiteString(line)
}

func (h *CustomJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	request := ctx.Value(requestKey)

	fields := make(map[string]interface{}, r.NumAttrs()+3)
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})
	if request != nil {
		fields["request"] = request
	}
	fields["message"] = r.Message
	fields["level"] = r.Level.String()
	fields["time"] = r.Time.Format(time.RFC3339Nano)

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	h.w.Write([]byte(colorize(r.Level, string(b)) + "\n"))

	return nil
}

func ConvertContext(c *gin.Context) context.Context {
	requestId := c.GetString("requestId")
	if requestId == "" {
		var err error
		requestId, err = uuid.GenerateShortID()
		if err != nil {
			slog.Error("Error generating thread ID", "error", err)
			requestId = "unknown"
		}
		c.Set("requestId", requestId)
	}

	return NewRequestContext(context.Background(), map[string]string{
		"id":     requestId,
		"method": c.Request.Method,
		"url":    c.Request.URL.Path,
	})
}
