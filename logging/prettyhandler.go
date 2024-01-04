package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/fatih/color"
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
	fields := extractFields(r)
	fields["time"] = r.Time.Format(time.RFC3339Nano)
	fields["level"] = colorizeLevel(r.Level, r.Level.String())
	fields["msg"] = r.Message
	var fieldsStr string
	if len(fields) > 0 {
		fields, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
		fieldsStr = string(fields)
	}

	fmt.Println(fieldsStr)

	jsonData, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	// // Check if output is a file and handle accordingly
	// if _, isFile := h.l.Writer().(*os.File); isFile {
	// 	// Strip color codes from file output
	// 	strippedMsg := str.StripANSI(fieldsStr)
	// 	fmt.Sprintf("%s %s %s", levelColor, messageColor, fieldsStr)
	// 	h.l.Println(fields["time"], strippedMsg)
	// } else {
	// 	// Output to STDOUT with color
	// 	h.l.Println(timeStr, levelColor, messageColor, fieldsStr)
	// }

	coloredOutput := h.colorizeBasedOnLevel(r.Level, string(jsonData))
	h.l.Println(coloredOutput)

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

func (h *PrettyHandler) colorizeBasedOnLevel(level slog.Level, message string) string {
	var colorMessage string
	switch level {
	case slog.LevelInfo:
		colorMessage = color.New(color.FgBlue).Sprint(message)
	case slog.LevelError:
		colorMessage = color.New(color.FgRed).Sprint(message)
	case slog.LevelDebug:
		colorMessage = color.New(color.FgMagenta).Sprint(message)
	default:
		colorMessage = message
	}
	return colorMessage
}
