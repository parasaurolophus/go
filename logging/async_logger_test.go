// Copright Kirk Rader 2024

package logging

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"slices"
	"testing"
)

func TestAsyncTrace(t *testing.T) {
	count := 0
	ch := make(chan int)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := NewAsync(writer, nil)
	defer logger.Stop()
	logger.SetVerbosity(TRACE)
	logger.Trace(
		func() string {
			count += 1
			ch <- count
			close(ch)
			return "trace"
		})
	for c := range ch {
		if c != 1 {
			t.Errorf("expected 1, got %d", c)
		}
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestAsyncTraceContext(t *testing.T) {
	count := 0
	ch := make(chan int)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := NewAsync(writer, nil)
	defer logger.Stop()
	logger.SetVerbosity(TRACE)
	logger.TraceContext(
		context.Background(),
		func() string {
			count += 1
			ch <- count
			close(ch)
			return "trace"
		})
	for c := range ch {
		if c != 1 {
			t.Errorf("expected 1, got %d", c)
		}
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestAsyncFine(t *testing.T) {
	count := 0
	ch := make(chan int)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := NewAsync(writer, nil)
	defer logger.Stop()
	logger.Fine(
		func() string {
			count += 1
			ch <- count
			close(ch)
			return "trace"
		})
	for c := range ch {
		if c != 1 {
			t.Errorf("expected 1, got %d", c)
		}
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestAsyncFineContext(t *testing.T) {
	count := 0
	ch := make(chan int)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := NewAsync(writer, nil)
	defer logger.Stop()
	logger.FineContext(
		context.Background(),
		func() string {
			count += 1
			ch <- count
			close(ch)
			return "trace"
		})
	for c := range ch {
		if c != 1 {
			t.Errorf("expected 1, got %d", c)
		}
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestAsyncOptional(t *testing.T) {
	count := 0
	ch := make(chan int)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := NewAsync(writer, nil)
	defer logger.Stop()
	logger.Optional(
		func() string {
			count += 1
			ch <- count
			close(ch)
			return "optional"
		})
	for c := range ch {
		if c != 1 {
			t.Errorf("expected 1, got %d", c)
		}
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestAsyncOptionalContext(t *testing.T) {
	count := 0
	ch := make(chan int)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := NewAsync(writer, nil)
	defer logger.Stop()
	logger.OptionalContext(
		context.Background(),
		func() string {
			count += 1
			ch <- count
			close(ch)
			return "optional"
		})
	for c := range ch {
		if c != 1 {
			t.Errorf("expected 1, got %d", c)
		}
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestAsyncAlways(t *testing.T) {
	count := 0
	ch := make(chan int)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := NewAsync(writer, nil)
	defer logger.Stop()
	logger.Always(
		func() string {
			count += 1
			ch <- count
			close(ch)
			return "optional"
		})
	for c := range ch {
		if c != 1 {
			t.Errorf("expected 1, got %d", c)
		}
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestAsyncAlwaysContext(t *testing.T) {
	count := 0
	ch := make(chan int)
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	logger := NewAsync(writer, nil)
	defer logger.Stop()
	logger.AlwaysContext(
		context.Background(),
		func() string {
			count += 1
			ch <- count
			close(ch)
			return "optional"
		})
	for c := range ch {
		if c != 1 {
			t.Errorf("expected 1, got %d", c)
		}
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestAsyncBaseAttributes(t *testing.T) {
	logger := NewAsync(os.Stdout, nil)
	defer logger.Stop()
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

func TestAsyncBaseTags(t *testing.T) {
	logger := NewAsync(os.Stdout, nil)
	defer logger.Stop()
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

func TestAsyncVerbosity(t *testing.T) {
	logger := NewAsync(os.Stdout, nil)
	defer logger.Stop()
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

func TestAsyncEnable(t *testing.T) {

	logger := NewAsync(os.Stdout, nil)
	defer logger.Stop()
	if logger.Enabled(TRACE) {
		t.Errorf("TRACE should be disabled by default")
	}
	if !logger.Enabled(FINE) {
		t.Errorf("FINE should be enabled by default")
	}
}

func TestAsyncIsEnableContext(t *testing.T) {

	logger := NewAsync(os.Stdout, nil)
	defer logger.Stop()
	if logger.EnabledContext(context.Background(), TRACE) {
		t.Errorf("TRACE should be disabled by default")
	}
	if !logger.EnabledContext(context.Background(), FINE) {
		t.Errorf("FINE should be enabled by default")
	}
}
