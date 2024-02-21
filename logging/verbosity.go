// Copright Kirk Rader 2024

package logging

import (
	"fmt"
	"log/slog"
	"parasaurolophus/go/stacktraces"
	"strconv"
)

// Verbosity-based nomenclature used in place of slog.Level.
type Verbosity int

// Mapping of Verbosity to slog.Level values.
//
// Generally, assume that only ALWAYS will be enabled in production environments
// and that TRACE will never be enabled outside of development environments.
const (

	// Only emit a log entry when extremely verbose output is specified.
	//
	// Intended for use in development environments for focused debugging
	// sessions. This should never be enabled outside of development
	// environments. Any logging that might potentially reveal PII, SPI or
	// critically sensitive security information must only be written at TRACE
	// level in environments where only synthetic or redacted data is in use.
	TRACE = Verbosity(slog.LevelDebug)

	// Only emit a log entry when unusually verbose output is specified.
	//
	// Intended for use in development environments for everyday testing and
	// troubleshooting prior to a release candidate being deployed.
	FINE = Verbosity(slog.LevelInfo)

	// Only emit a log entry when moderately verbose output is specified.
	//
	// Intended for use in testing and staging environments, e.g. during
	// acceptance and regression tests before release to production.
	OPTIONAL = Verbosity(slog.LevelWarn)

	// Always emit a log entry.
	//
	// Intended for production environments to drive monitoring, alerting and
	// reporting.
	ALWAYS = Verbosity(slog.LevelError)
)

// Implement fmt.Stringer interface for Verbosity.
func (v Verbosity) String() string {

	switch v {

	case TRACE:
		return "TRACE"

	case FINE:
		return "FINE"

	case OPTIONAL:
		return "OPTIONAL"

	case ALWAYS:
		return "ALWAYS"

	default:
		return strconv.Itoa(int(v))
	}
}

// Implement the fmt.Scanner interface.
func (v *Verbosity) Scan(state fmt.ScanState, verb rune) error {

	b, err := state.Token(true, nil)

	if err != nil {
		return err
	}

	token := string(b)

	switch token {

	case TRACE.String():
		*v = TRACE

	case FINE.String():
		*v = FINE

	case OPTIONAL.String():
		*v = OPTIONAL

	case ALWAYS.String():
		*v = ALWAYS

	default:
		n, err := strconv.Atoi(token)
		if err != nil {
			return stacktraces.New(fmt.Sprintf("unsupported verbosity token: '%s'", token), nil)
		}
		*v = Verbosity(n)
	}

	return nil
}
