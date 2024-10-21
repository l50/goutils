package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdlog "log"
	"log/slog"
	"os"
	"time"

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
// Level: The log level for the handler.
// OutputType: The type of output for the handler.
type PrettyHandler struct {
	slog.Handler
	l          *stdlog.Logger
	Level      slog.Level
	OutputType OutputType
}

// NewPrettyHandler creates a new PrettyHandler with specified output
// writer and options. It configures the PrettyHandler for handling
// log messages with optional colorization and structured formatting.
//
// **Parameters:**
//
// out: Output writer where log messages will be written.
// opts: PrettyHandlerOptions for configuring the handler.
// outputType: Type of output for the handler.
//
// **Returns:**
//
// *PrettyHandler: A new instance of PrettyHandler.
func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions, outputType OutputType) *PrettyHandler {
	h := &PrettyHandler{
		Handler:    slog.NewJSONHandler(out, &opts.SlogOpts),
		l:          stdlog.New(out, "", 0), // Use stdlog.New
		Level:      opts.SlogOpts.Level.Level(),
		OutputType: outputType,
	}
	return h
}

// Handle processes and outputs a log record using the PrettyHandler.
// It supports both colorized and non-colorized log messages and can
// output in JSON format if not writing to a terminal.
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
	if r.Level < h.Level {
		return nil // Skip log records below the minimum level
	}

	fields, err := h.parseLogRecord(r)
	if err != nil {
		return err
	}

	if h.OutputType == JSONOutput {
		return h.outputJSON(fields)
	}
	return h.outputFormatted(fields, r.Level)
}

// outputToFile determines if the output is being written to a file
// rather than a terminal, in which case it returns true.
//
// **Returns:**
//
// bool: True if output is to a file, false otherwise.
func (h *PrettyHandler) outputToFile() bool {
	_, isFile := h.l.Writer().(*os.File)
	isTerminal := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	return isFile && !isTerminal
}

// outputJSON marshals the log fields into JSON format and outputs
// them using the logger. This is used when logging to a file or
// non-terminal output.
//
// **Parameters:**
//
// fields: Log fields to be marshaled and output.
//
// **Returns:**
//
// error: An error if JSON marshaling or output fails.
func (h *PrettyHandler) outputJSON(fields map[string]interface{}) error {
	jsonData, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}
	h.l.Println(string(jsonData))
	return nil
}

// outputFormatted formats the log fields into a colorized string
// based on the log level and outputs it. This is used for terminal
// outputs.
//
// **Parameters:**
//
// fields: Log fields to be formatted and outputted.
// level: Log level used for determining the color.
//
// **Returns:**
//
// error: An error if formatting or output fails.
func (h *PrettyHandler) outputFormatted(fields map[string]interface{}, level slog.Level) error {
	timeStr := fields["time"].(string)
	msg := fields["msg"].(string)

	var levelColor string
	if h.OutputType == ColorOutput {
		switch level {
		case slog.LevelDebug:
			levelColor = "\033[36m" // Cyan
		case slog.LevelInfo:
			levelColor = "\033[32m" // Green
		case slog.LevelWarn:
			levelColor = "\033[33m" // Yellow
		case slog.LevelError:
			levelColor = "\033[31m" // Red
		default:
			levelColor = "\033[0m" // Default
		}
	}

	resetColor := "\033[0m"
	if h.OutputType != ColorOutput {
		levelColor = ""
		resetColor = ""
	}

	formattedMsg := fmt.Sprintf("%s %s%s%s: %s\n", timeStr, levelColor, level.String(), resetColor, msg)
	_, err := h.l.Writer().Write([]byte(formattedMsg))
	return err
}

// parseLogRecord parses the given slog.Record into a map of log fields.
// It handles both JSON and non-JSON log messages.
//
// **Parameters:**
//
// r: The slog.Record to be parsed.
//
// **Returns:**
//
// map[string]interface{}: Parsed log fields.
// error: An error if parsing fails.
func (h *PrettyHandler) parseLogRecord(r slog.Record) (map[string]interface{}, error) {
	var fields map[string]interface{}

	if json.Valid([]byte(r.Message)) {
		if err := json.Unmarshal([]byte(r.Message), &fields); err != nil {
			return nil, err
		}
	} else {
		// Consider non-JSON messages as valid and create a field map
		fields = map[string]interface{}{
			"time":  r.Time.Format(time.DateTime),
			"level": r.Level.String(),
			"msg":   r.Message,
		}
	}

	return fields, nil
}

// colorizeBasedOnLevel applies color to the given log level string
// based on its severity.
//
// **Parameters:**
//
// level: Log level to be colorized.
//
// **Returns:**
//
// string: The colorized log level string.
func (h *PrettyHandler) colorizeBasedOnLevel(level slog.Level) string {
	// Create a new color object based on the log level
	colorAttr := determineColorAttribute(level)
	c := color.New(colorAttr)

	// Apply color only to the level part
	coloredLevel := c.Sprint(level.String())

	return coloredLevel
}

// determineColorAttribute returns the color attribute corresponding
// to the given log level.
//
// **Parameters:**
//
// level: Log level for which to determine the color.
//
// **Returns:**
//
// color.Attribute: The color attribute for the given log level.
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
