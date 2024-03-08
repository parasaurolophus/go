// Copright Kirk Rader 2024

package logging

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"parasaurolophus/go/stacktraces"
	"parasaurolophus/go/stacktraces_test"
)

func TestSyncTrace(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	logger.SetVerbosity(TRACE)
	logger.Trace(
		func() string {
			return "trace"
		})
	writer.Flush()
	b := buffer.Bytes()
	type logEntry struct {
		Time      string   `json:"time"`
		Verbosity string   `json:"verbosity"`
		Msg       string   `json:"msg"`
		Tags      []string `json:"tags,omitempty"`
	}
	entry := logEntry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry; %s", err.Error())
	}
	if entry.Verbosity != "TRACE" {
		t.Fatalf(
			"expected verbosity 'TRACE', got '%s'",
			entry.Verbosity)
	}
	if entry.Msg != "trace" {
		t.Fatalf(
			"expected msg 'trace', got '%s'",
			entry.Msg)
	}
	if len(entry.Tags) != 0 {
		t.Fatalf("expected no tags, got %#v", entry.Tags)
	}
}

func TestSyncTraceContext(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	logger.SetVerbosity(TRACE)
	ctx := context.Background()
	logger.TraceContext(
		ctx,
		func() string {
			return "trace"
		})
	writer.Flush()
	b := buffer.Bytes()
	type logEntry struct {
		Time      string   `json:"time"`
		Verbosity string   `json:"verbosity"`
		Msg       string   `json:"msg"`
		Tags      []string `json:"tags,omitempty"`
	}
	entry := logEntry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry; %s", err.Error())
	}
	if entry.Verbosity != "TRACE" {
		t.Fatalf(
			"expected verbosity 'TRACE', got '%s'",
			entry.Verbosity)
	}
	if entry.Msg != "trace" {
		t.Fatalf(
			"expected msg 'trace', got '%s'",
			entry.Msg)
	}
	if len(entry.Tags) != 0 {
		t.Fatalf("expected no tags, got %#v", entry.Tags)
	}
}

func TestSyncFine(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	logger.Fine(
		func() string {
			return "fine"
		})
	writer.Flush()
	b := buffer.Bytes()
	type logEntry struct {
		Time      string `json:"time"`
		Verbosity string `json:"verbosity"`
		Msg       string `json:"msg"`
	}
	entry := logEntry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry; %s", err.Error())
	}
	if entry.Verbosity != "FINE" {
		t.Fatalf(
			"expected verbosity 'FINE', got '%s'",
			entry.Verbosity)
	}
	if entry.Msg != "fine" {
		t.Fatalf(
			"expected msg 'fine', got '%s'",
			entry.Msg)
	}
}

func TestSyncFineContext(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	ctx := context.Background()
	logger.FineContext(
		ctx,
		func() string {
			return "fine"
		})
	writer.Flush()
	b := buffer.Bytes()
	type logEntry struct {
		Time      string `json:"time"`
		Verbosity string `json:"verbosity"`
		Msg       string `json:"msg"`
	}
	entry := logEntry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry; %s", err.Error())
	}
	if entry.Verbosity != "FINE" {
		t.Fatalf(
			"expected verbosity 'FINE', got '%s'",
			entry.Verbosity)
	}
	if entry.Msg != "fine" {
		t.Fatalf(
			"expected msg 'fine', got '%s'",
			entry.Msg)
	}
}

func TestSyncOptional(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	logger.Optional(
		func() string {
			return "optional"
		})
	writer.Flush()
	b := buffer.Bytes()
	type logEntry struct {
		Time      string `json:"time"`
		Verbosity string `json:"verbosity"`
		Msg       string `json:"msg"`
	}
	entry := logEntry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry; %s", err.Error())
	}
	if entry.Verbosity != "OPTIONAL" {
		t.Fatalf(
			"expected verbosity 'OPTIONAL', got '%s'",
			entry.Verbosity)
	}
	if entry.Msg != "optional" {
		t.Fatalf(
			"expected msg 'optional', got '%s'",
			entry.Msg)
	}
}

func TestSyncOptionalContext(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	ctx := context.Background()
	logger.OptionalContext(
		ctx,
		func() string {
			return "optional"
		})
	writer.Flush()
	b := buffer.Bytes()
	type logEntry struct {
		Time      string `json:"time"`
		Verbosity string `json:"verbosity"`
		Msg       string `json:"msg"`
	}
	entry := logEntry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry; %s", err.Error())
	}
	if entry.Verbosity != "OPTIONAL" {
		t.Fatalf(
			"expected verbosity 'OPTIONAL', got '%s'",
			entry.Verbosity)
	}
	if entry.Msg != "optional" {
		t.Fatalf(
			"expected msg 'optional', got '%s'",
			entry.Msg)
	}
}

func TestSyncAlways(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	type Counters struct {
		Error1 uint `json:"error1"`
		Error2 uint `json:"error2"`
	}
	counters := Counters{
		Error1: 0,
		Error2: 0,
	}
	replacerCalled := false
	builderCalled := false
	baseTags := []string{"base"}
	additionalTags := []string{"additional"}
	options := LoggerOptions{
		BaseAttributes: []any{"counters", &counters},
		BaseTags:       baseTags,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			replacerCalled = true
			return a
		},
	}
	logger := New(writer, &options)
	counters.Error1 += 1
	logger.Always(
		func() string {
			builderCalled = true
			return "always"
		},
		STACKTRACE, nil,
		TAGS, additionalTags,
		"foo", "bar")
	writer.Flush()
	bytes := buffer.Bytes()
	type logEntry struct {
		Time       string   `json:"time"`
		Verbosity  string   `json:"verbosity"`
		Msg        string   `json:"msg"`
		Counters   Counters `json:"counters"`
		StackTrace string   `json:"stacktrace"`
		Tags       []string `json:"tags"`
		Foo        string   `json:"foo"`
	}
	entry := logEntry{}
	err := json.Unmarshal(bytes, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry; %s", err.Error())
	}
	if !replacerCalled {
		t.Fatalf("expected attribute replacer to have been called")
	}
	if !builderCalled {
		t.Fatalf("expected message builder to havve been called")
	}
	_, err = time.Parse(time.RFC3339Nano, entry.Time)
	if err != nil {
		t.Fatalf(
			"error parsing time '%s'; %s",
			entry.Time,
			err.Error())
	}
	if entry.Verbosity != "ALWAYS" {
		t.Fatalf(
			"expected verbosity to be 'ALWAYS', got '%s'",
			entry.Verbosity)
	}
	if entry.Msg != "always" {
		t.Fatalf(
			"expected msg to be 'always', got '%s'",
			entry.Msg)
	}
	if entry.Counters.Error1 != 1 {
		t.Fatalf(
			"expected Error1 to be 1, got %d",
			entry.Counters.Error1)
	}
	if entry.Counters.Error2 != 0 {
		t.Fatalf(
			"expected Error2 to be 0, got %d",
			entry.Counters.Error2)
	}
	name, _, err := stacktraces_test.FirstFunctionShort(entry.StackTrace)
	if err != nil {
		t.Fatalf("error parsing stack frames; %s", err.Error())
	}
	functionName := stacktraces.FunctionName()
	if name != functionName {
		t.Fatalf("expected first stack frame to be for '%s', got '%s'", functionName, name)
	}
	if entry.StackTrace == "" {
		t.Fatalf("expected stack trace not to be empty")
	}
	combinedTags := append(baseTags, additionalTags...)
	if len(entry.Tags) != len(combinedTags) {
		t.Fatalf("expected length of %#v to be 2, got %d", entry.Tags, len(entry.Tags))
	}
	for _, val := range combinedTags {
		if !slices.Contains[[]string](entry.Tags, val) {
			t.Fatalf("expected %#v to contain '%s'", entry.Tags, val)
		}
	}
	if entry.Foo != "bar" {
		t.Fatalf("expected foo to be 'bar', got '%s'", entry.Foo)
	}
}

func TestSyncAlwaysContext(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	type Counters struct {
		Error1 uint `json:"error1"`
		Error2 uint `json:"error2"`
	}
	counters := Counters{
		Error1: 0,
		Error2: 0,
	}
	replacerCalled := false
	builderCalled := false
	baseTags := []string{"base"}
	additionalTags := []string{"additional"}
	options := LoggerOptions{
		BaseAttributes: []any{"counters", &counters},
		BaseTags:       baseTags,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			replacerCalled = true
			return a
		},
	}
	logger := New(writer, &options)
	ctx := context.Background()
	counters.Error1 += 1
	logger.AlwaysContext(
		ctx,
		func() string {
			builderCalled = true
			return "always"
		},
		STACKTRACE, nil,
		TAGS, additionalTags,
		"foo", "bar")
	writer.Flush()
	bytes := buffer.Bytes()
	type logEntry struct {
		Time       string   `json:"time"`
		Verbosity  string   `json:"verbosity"`
		Msg        string   `json:"msg"`
		Counters   Counters `json:"counters"`
		StackTrace string   `json:"stacktrace"`
		Tags       []string `json:"tags"`
		Foo        string   `json:"foo"`
	}
	entry := logEntry{}
	err := json.Unmarshal(bytes, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry; %s", err.Error())
	}
	if !replacerCalled {
		t.Fatalf("expected attribute replacer to have been called")
	}
	if !builderCalled {
		t.Fatalf("expected message builder to havve been called")
	}
	_, err = time.Parse(time.RFC3339Nano, entry.Time)
	if err != nil {
		t.Fatalf(
			"error parsing time '%s'; %s",
			entry.Time,
			err.Error())
	}
	if entry.Verbosity != "ALWAYS" {
		t.Fatalf(
			"expected verbosity to be 'ALWAYS', got '%s'",
			entry.Verbosity)
	}
	if entry.Msg != "always" {
		t.Fatalf(
			"expected msg to be 'always', got '%s'",
			entry.Msg)
	}
	if entry.Counters.Error1 != 1 {
		t.Fatalf(
			"expected Error1 to be 1, got %d",
			entry.Counters.Error1)
	}
	if entry.Counters.Error2 != 0 {
		t.Fatalf(
			"expected Error2 to be 0, got %d",
			entry.Counters.Error2)
	}
	name, _, err := stacktraces_test.FirstFunctionShort(entry.StackTrace)
	if err != nil {
		t.Fatalf("error parsing stack frames; %s", err.Error())
	}
	functionName := stacktraces.FunctionName()
	if name != functionName {
		t.Fatalf("expected first stack frame to be for '%s', got '%s'", functionName, name)
	}
	if entry.StackTrace == "" {
		t.Fatalf("expected stack trace not to be empty")
	}
	combinedTags := append(baseTags, additionalTags...)
	if len(entry.Tags) != len(combinedTags) {
		t.Fatalf("expected length of %#v to be 2, got %d", entry.Tags, len(entry.Tags))
	}
	for _, val := range combinedTags {
		if !slices.Contains[[]string](entry.Tags, val) {
			t.Fatalf("expected %#v to contain '%s'", entry.Tags, val)
		}
	}
	if entry.Foo != "bar" {
		t.Fatalf("expected foo to be 'bar', got '%s'", entry.Foo)
	}
}

func TestSyncNilBuilder(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	type errorCounters struct {
		Error1 uint `json:"error1"`
		Error2 uint `json:"error2"`
	}
	counters := errorCounters{
		Error1: 0,
		Error2: 0,
	}
	options := LoggerOptions{
		BaseTags: []string{"test"},
	}
	logger := New(writer, &options)
	logger.SetBaseAttributes("counters", &counters)
	counters.Error1 += 1
	logger.Always(nil, STACKTRACE, 0)
	writer.Flush()
	bytes := buffer.Bytes()
	type logEntry struct {
		Time       string        `json:"time"`
		Verbosity  string        `json:"verbosity"`
		Msg        string        `json:"msg"`
		Counters   errorCounters `json:"counters"`
		StackTrace string        `json:"stacktrace"`
		Tags       []string      `json:"tags"`
	}
	entry := logEntry{}
	err := json.Unmarshal(bytes, &entry)
	if err != nil {
		t.Fatalf(
			"error unmarshaling log entry; %s",
			err.Error())
	}
	_, err = time.Parse(time.RFC3339Nano, entry.Time)
	if err != nil {
		t.Fatalf(
			"error parsing time '%s'; %s",
			entry.Time,
			err.Error())
	}
	if entry.Verbosity != "ALWAYS" {
		t.Fatalf(
			"expected verbosity to be 'ALWAYS', got '%s'",
			entry.Verbosity)
	}
	if entry.Msg != "" {
		t.Fatalf(
			"expected msg to be empty, got '%s'",
			entry.Msg)
	}
	if entry.Counters.Error1 != 1 {
		t.Fatalf(
			"expected Error1 to be 1, got %d",
			entry.Counters.Error1)
	}
	if entry.Counters.Error2 != 0 {
		t.Fatalf(
			"expected Error2 to be 0, got %d",
			entry.Counters.Error2)
	}
	name, _, err := stacktraces_test.FirstFunctionShort(entry.StackTrace)
	if err != nil {
		t.Fatalf("error parsing stack trace: %s", err.Error())
	}
	if name != "runtime.Callers" {
		t.Fatalf("expected stack trace to start with 'runtime.Callers', got '%s'", name)
	}
	if len(entry.Tags) != 1 {
		t.Fatalf("expected length of %#v to be 1", entry.Tags)
	}
	if !slices.Contains[[]string](entry.Tags, "test") {
		t.Fatalf("expected %#v to contain 'test'", entry.Tags)
	}
}

func TestSyncLazyEvaluation(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	replacerCalled := false
	options := LoggerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			replacerCalled = true
			return a
		},
	}
	logger := New(writer, &options)
	if logger.Enabled(TRACE) {
		t.Fatalf("expected TRACE to be disabled by default")
	}
	logger.Trace(
		func() string {
			t.Errorf("msg builder should not be called")
			return "error"
		})
	writer.Flush()
	b := buffer.Bytes()
	if len(b) > 0 {
		t.Errorf("no output should be written")
	}
	if replacerCalled {
		t.Errorf("replacer should not be called")
	}
}

func TestSyncIntTag(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	options := LoggerOptions{
		BaseTags: []string{"test"},
	}
	logger := New(writer, &options)
	logger.Always(nil, TAGS, 1)
	writer.Flush()
	b := buffer.Bytes()
	type Entry struct {
		Time      string   `json:"time"`
		Verbosity string   `json:"verbosity"`
		Msg       string   `json:"msg"`
		Tags      []string `json:"tags"`
	}
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry: %s", err.Error())
	}
	if len(entry.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(entry.Tags))
	}
	if !slices.Contains[[]string](entry.Tags, "test") {
		t.Fatalf("expected %#v to contain 'test'", entry.Tags)
	}
	if !slices.Contains[[]string](entry.Tags, "1") {
		t.Fatalf("expected %#v to contain 'test'", entry.Tags)
	}
}

func TestSyncStringTag(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	options := LoggerOptions{
		BaseTags: []string{"test"},
	}
	logger := New(writer, &options)
	logger.Always(nil, TAGS, "foo")
	writer.Flush()
	b := buffer.Bytes()
	type Entry struct {
		Time      string   `json:"time"`
		Verbosity string   `json:"verbosity"`
		Msg       string   `json:"msg"`
		Tags      []string `json:"tags"`
	}
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry: %s", err.Error())
	}
	if len(entry.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(entry.Tags))
	}
	if !slices.Contains[[]string](entry.Tags, "test") {
		t.Fatalf("expected %#v to contain 'test'", entry.Tags)
	}
	if !slices.Contains[[]string](entry.Tags, "foo") {
		t.Fatalf("expected %#v to contain 'test'", entry.Tags)
	}
}

func TestSyncStringerTag(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	options := LoggerOptions{
		BaseTags: []string{"test"},
	}
	logger := New(writer, &options)
	logger.Always(nil, TAGS, TRACE)
	writer.Flush()
	b := buffer.Bytes()
	type Entry struct {
		Time      string   `json:"time"`
		Verbosity string   `json:"verbosity"`
		Msg       string   `json:"msg"`
		Tags      []string `json:"tags"`
	}
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling log entry: %s", err.Error())
	}
	if len(entry.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(entry.Tags))
	}
	if !slices.Contains[[]string](entry.Tags, "test") {
		t.Fatalf("expected %#v to contain 'test'", entry.Tags)
	}
	if !slices.Contains[[]string](entry.Tags, TRACE.String()) {
		t.Fatalf("expected %#v to contain '%s'", entry.Tags, TRACE.String())
	}
}

func TestSyncUnrecognizedLevel(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	ctx := context.Background()
	lgr := New(writer, nil)
	wrapped := lgr.(*syncLogger).wrapped
	wrapped.Log(ctx, slog.Level(100), "slog.Level(100)")
	writer.Flush()
	b := buffer.Bytes()
	type Entry struct {
		Time      string `json:"time"`
		Verbosity string `json:"verbosity"`
		Msg       string `json:"msg"`
	}
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error unmarshaling entry: %s", err.Error())
	}
	if !strings.HasPrefix(entry.Verbosity, "ERROR+") {
		t.Fatalf("expected verbosity to start with 'ERROR+', got '%s'", entry.Verbosity)
	}
}

func TestSyncBaseAttributes(t *testing.T) {
	logger := New(os.Stdout, nil)
	actual := logger.BaseAttributes()
	if len(actual) != 0 {
		t.Fatalf("expected base attributes to be empty, got %v", actual)
	}
	expected := []any{"key", "value"}
	logger.SetBaseAttributes(expected...)
	actual = logger.BaseAttributes()
	if len(actual) != len(expected) {
		t.Fatalf("expected %v to be the same as %v", actual, expected)
	}
	for _, v := range expected {
		if !slices.Contains[[]any](actual, v) {
			t.Fatalf("expected %v to contain %v", actual, v)
		}
	}
}

func TestSyncBaseTags(t *testing.T) {
	logger := New(os.Stdout, nil)
	actual := logger.BaseTags()
	if len(actual) != 0 {
		t.Fatalf("expected base tags to be empty, got %v", actual)
	}
	expected := []string{"tag1", "tag2"}
	logger.SetBaseTags(expected...)
	actual = logger.BaseTags()
	if len(actual) != len(expected) {
		t.Fatalf("expected %v to be the same as %v", actual, expected)
	}
	for _, v := range expected {
		if !slices.Contains[[]string](actual, v) {
			t.Fatalf("expected %v to contain %v", actual, v)
		}
	}
}

func TestSyncVerbosity(t *testing.T) {
	logger := New(os.Stdout, nil)
	actual := logger.Verbosity()
	if actual != FINE {
		t.Fatalf("expected verbosity to be %d, got %d", FINE, actual)
	}
	logger.SetVerbosity(TRACE)
	actual = logger.Verbosity()
	if actual != TRACE {
		t.Fatalf("expected verbosity to be %d, got %d", TRACE, actual)
	}
}

func TestSyncOddAttributes(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	options := LoggerOptions{
		BaseAttributes: []any{"foo", "bar", "stacktrace", nil},
		BaseTags:       []string{"test"},
	}
	logger := New(writer, &options)
	logger.Always(nil, "baz", "waka", "hoo")
	writer.Flush()
	b := buffer.Bytes()
	type Entry struct {
		Time       string   `json:"time"`
		Verbosity  string   `json:"verbosity"`
		Msg        string   `json:"msg"`
		Foo        string   `json:"foo"`
		Baz        string   `json:"baz"`
		StackTrace string   `json:"stacktrace"`
		Tags       []string `json:"tags"`
	}
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error parsing log entry: %s", err.Error())
	}
	if entry.Foo != "bar" {
		t.Fatalf("expected value of 'foo' to be 'bar', got '%s'", entry.Foo)
	}
	if entry.Baz != "waka" {
		t.Fatalf("expected value of 'baz' to be 'waka', got '%s'", entry.Baz)
	}
	expected := stacktraces.FunctionName()
	actual, _, err := stacktraces_test.FirstFunctionShort(entry.StackTrace)
	if err != nil {
		t.Fatalf("error parsing stacktrace: %s", err.Error())
	}
	if actual != expected {
		t.Fatalf("expected function name to be '%s', got '%s'", expected, actual)
	}
}

func TestSyncBadKey(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	logger.Always(nil, "good1", 1, 10, "bad", "good2", 2)
	writer.Flush()
	b := buffer.Bytes()
	type Entry struct {
		Time      string `json:"time"`
		Verbosity string `json:"verbosity"`
		Msg       string `json:"msg"`
		Good1     int    `json:"good1"`
		Good2     int    `json:"good2"`
	}
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf("error parsing log entry: %s", err.Error())
	}
	if entry.Good1 != 1 {
		t.Fatalf("expected value of 'good1' to be 1, got %d", entry.Good1)
	}
	if entry.Good2 != 2 {
		t.Fatalf("expected value of 'good2' to be 2, got %d", entry.Good2)
	}
}

func TestSyncMessageBuilderPanic(t *testing.T) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	logger.Always(
		func() string { panic("deliberate") },
	)
	writer.Flush()
	s := buffer.String()
	parts := strings.Split(s, "\n")
	fmt.Println(parts[0])
	if len(parts) < 2 {
		t.Fatalf("expected 2 entries, got \"%s\"", s)
	}
	type Entry struct {
		Time       string   `json:"time"`
		Verbosity  string   `json:"verbosity"`
		Msg        string   `json:"msg"`
		StackTrace string   `json:"stacktrace"`
		Recovered  string   `json:"recovered"`
		Tags       []string `json:"tags"`
	}
	entry := Entry{}
	err := json.Unmarshal([]byte(parts[0]), &entry)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(entry)
	if len(entry.Tags) < 2 {
		t.Fatalf("expected 2 tags, git %#v", entry.Tags)
	}
	if !slices.Contains(entry.Tags, PANIC) {
		t.Errorf("expected %#v to contain \"%s\"", PANIC, entry.Tags)
	}
	if !slices.Contains(entry.Tags, INJECTED) {
		t.Errorf("expected %#v to contain \"%s\"", INJECTED, entry.Tags)
	}
	if entry.Recovered != "deliberate" {
		t.Errorf("expected \"deliberate\", got \"%s\"", entry.Recovered)
	}
}

func TestSyncNegativeStackTraceParam(t *testing.T) {
	_, expected, _, _, _ := stacktraces.FunctionInfo(-2)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	logger.Always(nil, STACKTRACE, -2)
	writer.Flush()
	b := buffer.Bytes()
	type Entry struct {
		Time       string `json:"time"`
		Verbosity  string `json:"verbosity"`
		Msg        string `json:"msg"`
		StackTrace string `json:"stacktrace"`
	}
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf(err.Error())
	}
	actual, _, err := stacktraces_test.FirstFunctionShort(entry.StackTrace)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, actual)
	}
}

func TestSyncZeroStackTraceParam(t *testing.T) {
	_, expected, _, _, _ := stacktraces.FunctionInfo(0)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	logger.Always(nil, STACKTRACE, 0)
	writer.Flush()
	b := buffer.Bytes()
	type Entry struct {
		Time       string `json:"time"`
		Verbosity  string `json:"verbosity"`
		Msg        string `json:"msg"`
		StackTrace string `json:"stacktrace"`
	}
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf(err.Error())
	}
	actual, _, err := stacktraces_test.FirstFunctionShort(entry.StackTrace)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, actual)
	}
}

func TestSyncStringStackTraceParam(t *testing.T) {
	expected := stacktraces.FunctionName()
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := New(writer, nil)
	logger.Always(nil, STACKTRACE, expected)
	writer.Flush()
	b := buffer.Bytes()
	type Entry struct {
		Time       string `json:"time"`
		Verbosity  string `json:"verbosity"`
		Msg        string `json:"msg"`
		StackTrace string `json:"stacktrace"`
	}
	entry := Entry{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatalf(err.Error())
	}
	actual, _, err := stacktraces_test.FirstFunctionShort(entry.StackTrace)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if actual != expected {
		t.Fatalf("expected \"%s\", got \"%s\"", expected, actual)
	}
}

func TestSyncIsEnableContext(t *testing.T) {

	logger := New(os.Stdout, nil)
	if logger.EnabledContext(context.Background(), TRACE) {
		t.Fatalf("TRACE should be disabled by default")
	}
}

func TestSyncStop(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("recovered %v from panic", r)
		}
	}()

	logger := New(os.Stdout, nil)
	logger.Stop()
}