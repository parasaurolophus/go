// Copyright Kirk Rader 2024

package logging

import (
	"context"
	"io"
	"sync"
)

type (

	// Parameters for syncLogger.log() sent via logChan.
	logParams struct {
		ctx        context.Context
		verbosity  Verbosity
		message    MessageBuilder
		attributes []any
	}

	asyncLogger struct {
		synced      *syncLogger
		logChan     chan logParams
		ackChan     chan bool
		settingsMtx *sync.Mutex
	}
)

// Log at TRACE verbosity.
func (l *asyncLogger) Trace(message MessageBuilder, attributes ...any) {

	l.logChan <- logParams{
		verbosity:  TRACE,
		message:    message,
		attributes: attributes,
	}
}

// Log at TRACE verbosity using the supplied context.
func (l *asyncLogger) TraceContext(ctx context.Context, message MessageBuilder, attributes ...any) {

	l.logChan <- logParams{
		ctx:        ctx,
		verbosity:  TRACE,
		message:    message,
		attributes: attributes,
	}
}

// Log at FINE verbosity.
func (l *asyncLogger) Fine(message MessageBuilder, attributes ...any) {

	l.logChan <- logParams{
		verbosity:  FINE,
		message:    message,
		attributes: attributes,
	}
}

// Log at FINE verbosity using the supplied context.
func (l *asyncLogger) FineContext(ctx context.Context, message MessageBuilder, attributes ...any) {

	l.logChan <- logParams{
		ctx:        ctx,
		verbosity:  FINE,
		message:    message,
		attributes: attributes,
	}
}

// Log at OPTIONAL verbosity.
func (l *asyncLogger) Optional(message MessageBuilder, attributes ...any) {

	l.logChan <- logParams{
		verbosity:  OPTIONAL,
		message:    message,
		attributes: attributes,
	}
}

// Log at OPTIONAL verbosity using the supplied context.
func (l *asyncLogger) OptionalContext(ctx context.Context, message MessageBuilder, attributes ...any) {

	l.logChan <- logParams{
		ctx:        ctx,
		verbosity:  OPTIONAL,
		message:    message,
		attributes: attributes,
	}
}

// Log at ALWAYS verbosity.
func (l *asyncLogger) Always(message MessageBuilder, attributes ...any) {

	l.logChan <- logParams{
		verbosity:  ALWAYS,
		message:    message,
		attributes: attributes,
	}
}

// Log at ALWAYS verbosity using the supplied context.
func (l *asyncLogger) AlwaysContext(ctx context.Context, message MessageBuilder, attributes ...any) {

	l.logChan <- logParams{
		ctx:        ctx,
		verbosity:  ALWAYS,
		message:    message,
		attributes: attributes,
	}
}

// Return true if and only the specified level is enabled.
func (l *asyncLogger) Enabled(verbosity Verbosity) bool {

	defer l.settingsMtx.Unlock()
	l.settingsMtx.Lock()
	return l.synced.Enabled(verbosity)
}

// Return true if and only the specified level is enabled using the given context.
func (l *asyncLogger) EnabledContext(ctx context.Context, verbosity Verbosity) bool {

	defer l.settingsMtx.Unlock()
	l.settingsMtx.Lock()
	return l.synced.EnabledContext(ctx, verbosity)
}

// Return the enabled verbosity level.
func (l *asyncLogger) Verbosity() Verbosity {

	defer l.settingsMtx.Unlock()
	l.settingsMtx.Lock()
	return l.synced.Verbosity()
}

// Update the enabled verbosity level.
func (l *asyncLogger) SetVerbosity(verbosity Verbosity) {

	defer l.settingsMtx.Unlock()
	l.settingsMtx.Lock()
	l.synced.SetVerbosity(verbosity)
}

// Return the base attributes.
func (l *asyncLogger) BaseAttributes() []any {

	defer l.settingsMtx.Unlock()
	l.settingsMtx.Lock()
	return l.synced.BaseAttributes()
}

// Update the base attributes.
func (l *asyncLogger) SetBaseAttributes(atrributes ...any) {

	defer l.settingsMtx.Unlock()
	l.settingsMtx.Lock()
	l.synced.SetBaseAttributes(atrributes...)
}

// Return the base tags.
func (l *asyncLogger) BaseTags() []string {

	defer l.settingsMtx.Unlock()
	l.settingsMtx.Lock()
	return l.synced.BaseTags()
}

// Update the base tags.
func (l *asyncLogger) SetBaseTags(tags ...string) {

	defer l.settingsMtx.Unlock()
	l.settingsMtx.Lock()
	l.synced.SetBaseTags(tags...)
}

// Stop the worker goroutine.
func (l *asyncLogger) Stop() {
	close(l.logChan)
	<-l.ackChan
}

// Invoked as a goroutine by NewAsync().
func (l asyncLogger) worker() {

	defer func() {
		l.ackChan <- true
		close(l.ackChan)
	}()

	for params := range l.logChan {

		log := func() {
			defer l.settingsMtx.Unlock()
			l.settingsMtx.Lock()
			l.synced.log(params.ctx, params.verbosity, params.message, params.attributes...)
		}

		log()
	}
}

func NewAsync(writer io.Writer, options *LoggerOptions) Logger {

	logger := new(asyncLogger)
	logger.ackChan = make(chan bool)
	logger.logChan = make(chan logParams)
	logger.synced = New(writer, options).(*syncLogger)
	logger.settingsMtx = new(sync.Mutex)
	go logger.worker()
	return logger
}
