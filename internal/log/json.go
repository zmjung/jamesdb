package log

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/zmjung/jamesdb/internal/uuid"
)

type CustomJSONHandler struct {
	slog.JSONHandler
}

type RequestContextKey string

const requestKey RequestContextKey = "request"

func NewRequestContext(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, requestKey, value)
}

func getRequest(ctx context.Context) (map[string]string, bool) {
	value, ok := ctx.Value(requestKey).(map[string]string)
	return value, ok
}

func (h *CustomJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	if request, exists := getRequest(ctx); exists {
		r.AddAttrs(slog.Any(string(requestKey), request))
	}
	return h.JSONHandler.Handle(ctx, r)
}

func GetContext(c *gin.Context) context.Context {
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
