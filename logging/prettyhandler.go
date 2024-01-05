package logging

import (
	"context"
	"encoding/json"
	"fmt"
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
	fields := make(map[string]interface{}, 0)
	fields["time"] = time.Now().Format(time.RFC3339Nano)
	fields["msg"] = str.StripANSI(r.Message)

	// Determine if output is to a terminal or file
	_, isFile := h.l.Writer().(*os.File)
	isTerminal := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	var finalLogMsg string
	if isFile && !isTerminal {
		jsonData, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
		finalLogMsg = string(jsonData)
	} else {
		coloredLevel := h.colorizeBasedOnLevel(r.Level)

		finalLogMsg = fmt.Sprintf("[%s] [%s] %s", fields["time"], coloredLevel, fields["msg"])
	}

	h.l.Println(finalLogMsg)

	return nil
}

func (h *PrettyHandler) colorizeBasedOnLevel(level slog.Level) string {
	// Create a new color object based on the log level
	colorAttr := determineColorAttribute(level)
	c := color.New(colorAttr)

	// Apply color only to the level part
	coloredLevel := c.Sprint(level.String())

	return coloredLevel
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
