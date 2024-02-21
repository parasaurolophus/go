// Copright Kirk Rader 2024

package logging

import "log/slog"

// Configuration parameters for an instance of Logger.
type LoggerOptions struct {

	// Allow panics in Logger.log() to cause abnormal process termination.
	AllowPanics bool

	// Initial set of attributes that will be added to every log entry.
	BaseAttributes []any

	// Initial set of tags that will be added to every log entry.
	BaseTags []string

	// Pass through to HandlerOptions for the wrapped slog.Logger.
	AddSource bool

	// Shared slog.LevelVar, if desired; a Leveler will be created if this is
	// nil.
	Level *slog.LevelVar

	// If not nil, an attribute replacer function that will be called in
	// addition to replacing "level" attributes with "verbosiy" and other
	// special attribute handling.
	ReplaceAttr func([]string, slog.Attr) slog.Attr
}
