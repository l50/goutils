package logging

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
)

// PrettyHandlerOptions represents options used for configuring
// the PrettyHandler.
//
// **Attributes:**
//
// SlogOpts: Options for the underlying slog.Handler.
type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

// PrettyHandler is a custom log handler that provides colorized
// logging output. It wraps around slog.Handler and adds color to
// log messages based on their level.
//
// **Attributes:**
//
// Handler: The underlying slog.Handler used for logging.
// l: Standard logger used for outputting log messages.
type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

// NewPrettyHandler creates a new PrettyHandler with specified output
// writer and options. It configures a PrettyHandler for colorized
// logging output.
//
// **Parameters:**
//
// out: Output writer where log messages will be written.
// opts: PrettyHandlerOptions for configuring the handler.
//
// **Returns:**
//
// *PrettyHandler: A new instance of PrettyHandler.
func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
	return h
}

// Handle formats and outputs a log message for PrettyHandler. It
// colorizes the log level, message, and adds structured fields
// to the log output.
//
// **Parameters:**
//
// ctx: Context for the log record.
// r: The log record containing log data.
//
// **Returns:**
//
// error: An error if any issue occurs during log handling.
func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	levelColor := colorizeLevel(r.Level)
	timeStr := r.Time.Format("[15:05:05.000]")
	messageColor := colorizeMessage(r.Message, r.Level)

	// Extract fields and marshal if they exist
	fields := extractFields(r)
	var fieldsStr string
	if len(fields) > 0 {
		marshaledFields, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
		fieldsStr = string(marshaledFields)
	}

	// Construct and print the log message
	if fieldsStr != "" {
		h.l.Println(timeStr, levelColor, messageColor, fieldsStr)
	} else {
		// Omit fields part if it's empty
		h.l.Println(timeStr, levelColor, messageColor)
	}

	return nil
}

func extractFields(r slog.Record) map[string]interface{} {
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})
	return fields
}
