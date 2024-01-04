package logging

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
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
	colorAttr := colorizeLevel(r.Level)
	levelColor := color.New(colorAttr).Sprint(r.Level.String())
	timeStr := r.Time.Format("[15:05:05.000]")
	messageColor := color.New(colorAttr).Sprint(r.Message)

	fields := extractFields(r)
	var fieldsStr string
	if len(fields) > 0 {
		fieldsJSON, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
		fieldsStr = string(fieldsJSON)
	}

	// Determine if output should have color
	colorOutput := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	// h.l.Println(timeStr, levelColor, messageColor, fieldsStr)
	if _, isFile := h.l.Writer().(*os.File); isFile || !colorOutput {
		// Strip color codes from file output or non-color terminals
		// strippedMsg := str.StripANSI(fmt.Sprintf("%s %s %s", levelColor, messageColor, fieldsStr))
		h.l.Println(timeStr, levelColor, messageColor, fieldsStr)
		// h.l.Println(timeStr, strippedMsg)
	} else {
		// Output to STDOUT with color
		h.l.Println(timeStr, levelColor, messageColor, fieldsStr)
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

func colorizeLevel(level slog.Level) color.Attribute {
	switch level {
	case slog.LevelDebug:
		return color.FgMagenta
	case slog.LevelInfo:
		return color.FgBlue
	case slog.LevelWarn:
		return color.FgYellow
	case slog.LevelError:
		return color.FgRed
	default:
		return color.Reset
	}
}
