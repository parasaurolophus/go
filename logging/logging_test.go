// Copright Kirk Rader 2024

package logging

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"parasaurolophus/go/stacktraces"
)

func firstFunction(strackTrace string) (string, int, error) {

	frame := strings.Split(strackTrace, "<")[0]
	parts := strings.Split(frame, ":")

	if len(parts) < 1 {
		return "", 0, errors.New("no colon in stack frame")
	}

	name := strings.Trim(strings.Split(parts[1], "[")[0], " ")
	n, err := strconv.Atoi(parts[0])
	return name, n, err
}

func TestAlways(t *testing.T) {

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
	logger.Always(
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
		t.Fatalf("TestAlways: error unmarshaling log entry; %s", err.Error())
	}

	if !replacerCalled {
		t.Fatalf("TestAlways: expected attribute replacer to have been called")
	}

	if !builderCalled {
		t.Fatalf("TestAlways: expected message builder to havve been called")
	}

	_, err = time.Parse(time.RFC3339Nano, entry.Time)

	if err != nil {
		t.Fatalf(
			"TestAlways: error parsing time '%s'; %s",
			entry.Time,
			err.Error())
	}

	if entry.Verbosity != "ALWAYS" {
		t.Fatalf(
			"TestAlways: expected verbosity to be 'ALWAYS', got '%s'",
			entry.Verbosity)
	}

	if entry.Msg != "always" {
		t.Fatalf(
			"TestAlways: expected msg to be 'always', got '%s'",
			entry.Msg)
	}

	if entry.Counters.Error1 != 1 {
		t.Fatalf(
			"TestAlways: expected Error1 to be 1, got %d",
			entry.Counters.Error1)
	}

	if entry.Counters.Error2 != 0 {
		t.Fatalf(
			"TestAlways: expected Error2 to be 0, got %d",
			entry.Counters.Error2)
	}

	functionName := stacktraces.FunctionName()
	name, _, err := firstFunction(entry.StackTrace)

	if err != nil {
		t.Fatalf("TestAlways: error parsing stack frames; %s", err.Error())
	}

	if name != functionName {
		t.Fatalf("TestAlways: expected first stack frame to be for '%s', got '%s'", functionName, name)
	}

	if entry.StackTrace == "" {
		t.Fatalf("TestAlways: expected stack trace not to be empty")
	}

	combinedTags := append(baseTags, additionalTags...)

	if len(entry.Tags) != len(combinedTags) {
		t.Fatalf("TestAlways: expected length of %#v to be 2, got %d", entry.Tags, len(entry.Tags))
	}

	for _, val := range combinedTags {

		if !slices.Contains[[]string](entry.Tags, val) {
			t.Fatalf("TestAlwas: expected %#v to contain '%s'", entry.Tags, val)
		}
	}

	if entry.Foo != "bar" {
		t.Fatalf("TestAlways: expected foo to be 'bar', got '%s'", entry.Foo)
	}
}

func TestNilBuilder(t *testing.T) {

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
	ctx := context.Background()
	logger.Always(ctx, nil, STACKTRACE, nil)

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
			"TestNilBuilder: error unmarshaling log entry; %s",
			err.Error())
	}

	_, err = time.Parse(time.RFC3339Nano, entry.Time)

	if err != nil {
		t.Fatalf(
			"TestNilBuilder: error parsing time '%s'; %s",
			entry.Time,
			err.Error())
	}

	if entry.Verbosity != "ALWAYS" {
		t.Fatalf(
			"TestNilBuilder: expected verbosity to be 'ALWAYS', got '%s'",
			entry.Verbosity)
	}

	if entry.Msg != "" {
		t.Fatalf(
			"TestNilBuilder: expected msg to be empty, got '%s'",
			entry.Msg)
	}

	if entry.Counters.Error1 != 1 {
		t.Fatalf(
			"TestNilBuilder: expected Error1 to be 1, got %d",
			entry.Counters.Error1)
	}

	if entry.Counters.Error2 != 0 {
		t.Fatalf(
			"TestNilBuilder: expected Error2 to be 0, got %d",
			entry.Counters.Error2)
	}

	functionName := stacktraces.FunctionName()
	name, _, err := firstFunction(entry.StackTrace)

	if err != nil {
		t.Fatalf("TestNilBuilder: error parsing stack trace: %s", err.Error())
	}

	if name != functionName {
		t.Fatalf("TestNilBuilder: expected stack trace not start with '%s', got '%s'", functionName, name)
	}

	if len(entry.Tags) != 1 {
		t.Fatalf("TestNilBuilder: expected length of %#v to be 1", entry.Tags)
	}

	if !slices.Contains[[]string](entry.Tags, "test") {
		t.Fatalf("TestNilBuilder: expected %#v to contain 'test'", entry.Tags)
	}
}

func TestLazyEvaluation(t *testing.T) {

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
	ctx := context.Background()

	if logger.Enabled(ctx, TRACE) {
		t.Fatalf("TestLazyEvaluation: expected TRACE to be disabled by default")
	}

	logger.Trace(
		ctx,
		func() string {
			t.Fatalf("TestLazyEvaluation: msg builder should not be called")
			return "error"
		})

	writer.Flush()
	b := buffer.Bytes()

	if len(b) > 0 {
		t.Fatalf("TestLazyEvaluation: no output should be written")
	}

	if replacerCalled {
		t.Fatalf("TestLazyEvaluation: replacer should not be called")
	}
}

func TestTrace(t *testing.T) {

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)

	logger := New(writer, nil)
	logger.SetVerbosity(TRACE)
	ctx := context.Background()
	logger.Trace(
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
		t.Fatalf("TestTrace: error unmarshaling log entry; %s", err.Error())
	}

	if entry.Verbosity != "TRACE" {
		t.Fatalf(
			"TestTrace: expected verbosity 'TRACE', got '%s'",
			entry.Verbosity)
	}

	if entry.Msg != "trace" {
		t.Fatalf(
			"TestTrace: expected msg 'trace', got '%s'",
			entry.Msg)
	}

	if len(entry.Tags) != 0 {
		t.Fatalf("TestTrace: expected no tags, got %#v", entry.Tags)
	}
}

func TestFine(t *testing.T) {

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)

	logger := New(writer, nil)
	ctx := context.Background()
	logger.Fine(
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
		t.Fatalf("TestFine: error unmarshaling log entry; %s", err.Error())
	}

	if entry.Verbosity != "FINE" {
		t.Fatalf(
			"TestFine: expected verbosity 'FINE', got '%s'",
			entry.Verbosity)
	}

	if entry.Msg != "fine" {
		t.Fatalf(
			"TestFine: expected msg 'fine', got '%s'",
			entry.Msg)
	}
}

func TestOptional(t *testing.T) {

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)

	logger := New(writer, nil)
	ctx := context.Background()
	logger.Optional(
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
		t.Fatalf("TestOptional: error unmarshaling log entry; %s", err.Error())
	}

	if entry.Verbosity != "OPTIONAL" {
		t.Fatalf(
			"TestOptional: expected verbosity 'OPTIONAL', got '%s'",
			entry.Verbosity)
	}

	if entry.Msg != "optional" {
		t.Fatalf(
			"TestOptional: expected msg 'optional', got '%s'",
			entry.Msg)
	}
}

func TestIntTag(t *testing.T) {

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)

	options := LoggerOptions{
		BaseTags: []string{"test"},
	}

	logger := New(writer, &options)
	ctx := context.Background()
	logger.Always(ctx, nil, TAGS, 1)
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
		t.Fatalf("TestIntTag: error unmarshaling log entry: %s", err.Error())
	}

	if len(entry.Tags) != 2 {
		t.Fatalf("TestIntTag: expected 2 tags, got %d", len(entry.Tags))
	}

	if !slices.Contains[[]string](entry.Tags, "test") {
		t.Fatalf("TestIntTag: expected %#v to contain 'test'", entry.Tags)
	}

	if !slices.Contains[[]string](entry.Tags, "1") {
		t.Fatalf("TestIntTag: expected %#v to contain 'test'", entry.Tags)
	}
}

func TestStringTag(t *testing.T) {

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)

	options := LoggerOptions{
		BaseTags: []string{"test"},
	}

	logger := New(writer, &options)
	ctx := context.Background()
	logger.Always(ctx, nil, TAGS, "foo")
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
		t.Fatalf("TestIntTag: error unmarshaling log entry: %s", err.Error())
	}

	if len(entry.Tags) != 2 {
		t.Fatalf("TestIntTag: expected 2 tags, got %d", len(entry.Tags))
	}

	if !slices.Contains[[]string](entry.Tags, "test") {
		t.Fatalf("TestIntTag: expected %#v to contain 'test'", entry.Tags)
	}

	if !slices.Contains[[]string](entry.Tags, "foo") {
		t.Fatalf("TestIntTag: expected %#v to contain 'test'", entry.Tags)
	}
}

func TestUnrecognizedLevel(t *testing.T) {

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	ctx := context.Background()

	lgr := New(writer, nil)
	wrapped := lgr.wrapped

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
		t.Fatalf("TestUnrecognizedLevel: error unmarshaling entry: %s", err.Error())
	}

	if !strings.HasPrefix(entry.Verbosity, "ERROR+") {
		t.Fatalf("TestUnrecognizedLevel: expected verbosity to start with 'ERROR+', got '%s'", entry.Verbosity)
	}
}

func TestBaseAttributes(t *testing.T) {

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

func TestBaseTags(t *testing.T) {

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

func TestVerbosity(t *testing.T) {

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
