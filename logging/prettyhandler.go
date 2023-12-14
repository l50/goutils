package logging

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
)

// PrettyHandlerOptions - Options for PrettyHandler
type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

// PrettyHandler - Custom handler for color logging
type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
	return h
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	levelColor := colorizeLevel(r.Level)
	timeStr := r.Time.Format("[15:05:05.000]")
	messageColor := colorizeMessage(r.Message, r.Level)
	fields, err := json.MarshalIndent(extractFields(r), "", "  ")
	if err != nil {
		return err
	}

	h.l.Println(timeStr, levelColor, messageColor, string(fields))
	return nil
}
