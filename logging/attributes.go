// Copright Kirk Rader 2024

package logging

// Specially handled attributes.
const (

	// syncLogger.log() will include the value returned by recover() when
	// logging a panic.
	RECOVERED = "recovered"

	// Values of "stacktrace" attributes will be replaced with one-line stack
	// traces for the function that called the given logging method.
	STACKTRACE = "stacktrace"

	// Value will be merged with the currently configured
	// LoggerOptions.BaseTags.
	TAGS = "tags"
)
