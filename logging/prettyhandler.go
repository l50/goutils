package logging

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/l50/goutils/v2/str"
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
	fields := extractFields(r)
	fields["time"] = r.Time.Format(time.RFC3339Nano)
	fields["level"] = r.Level.String()
	fields["msg"] = r.Message

	// Marshal fields to JSON
	jsonData, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	// Check if the output is to a terminal or a file
	_, isFile := h.l.Writer().(*os.File)
	isTerminal := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	if isFile && !isTerminal {
		// Strip color codes from file output and print
		strippedMsg := str.StripANSI(string(jsonData))
		h.l.Println(strippedMsg)
	} else {
		// Output to STDOUT with color and print
		coloredOutput := h.colorizeBasedOnLevel(r.Level, string(jsonData))
		h.l.Println(coloredOutput)
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

func (h *PrettyHandler) colorizeBasedOnLevel(level slog.Level, message string) string {
	// Check if the output is a terminal
	isTerminal := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	if !isTerminal {
		return message // No color if not a terminal
	}

	colorAttr := determineColorAttribute(level)
	if colorAttr == color.Reset {
		return message // No color for default case
	}

	coloredOutput := color.New(colorAttr).Sprint(message)
	// fmt.Printf("Original Message: %s, Colored Message: %s\n", message, coloredOutput) // Debug print
	return coloredOutput
}

func determineColorAttribute(level slog.Level) color.Attribute {
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
