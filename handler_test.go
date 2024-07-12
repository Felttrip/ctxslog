package ctxslog_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/felttrip/ctxslog"
)

func TestNewHandler(t *testing.T) {
	mockHandler := &MockHandler{}
	handler := ctxslog.NewHandler(mockHandler)

	if handler == nil {
		t.Error("NewHandler returned nil")
	}
}

func TestHandlerEnabled(t *testing.T) {
	mockHandler := &MockHandler{}
	handler := ctxslog.NewHandler(mockHandler)

	ctx := context.Background()

	// Test with different log levels should just default down to the mock handler which is true for all
	levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}
	for _, level := range levels {
		enabled := handler.Enabled(ctx, level)
		if !enabled {
			t.Errorf("Handler.Enabled returned false for level %s", level)
		}
	}
}

func TestHandlerHandle(t *testing.T) {
	t.Run("Will Handle when nothing has been added to context", func(t *testing.T) {
		mockHandler := &MockHandler{}
		handler := ctxslog.NewHandler(mockHandler)

		ctx := context.Background()
		record := slog.Record{
			Level:   slog.LevelInfo,
			Message: "Test message",
		}

		err := handler.Handle(ctx, record)
		if err != nil {
			t.Errorf("Handler.Handle returned an error: %v", err)
		}

		if !mockHandler.HandleCalled {
			t.Error("MockHandler.Handle was not called")
		}

		if mockHandler.HandleRecord.Message != record.Message || mockHandler.HandleRecord.Level != record.Level {
			t.Error("MockHandler.Handle did not receive the correct record")
		}
	})
	t.Run("Will Handle when values have been added to context and add them to attrs", func(t *testing.T) {
		mockHandler := &MockHandler{}
		handler := ctxslog.NewHandler(mockHandler)

		ctx := context.Background()
		ctx = ctxslog.WithValue(ctx, "key1", "value1")
		ctx = ctxslog.WithValue(ctx, "key2", "value2")

		expectedAttrs := []slog.Attr{
			{Key: "key1", Value: slog.StringValue("value1")},
			{Key: "key2", Value: slog.StringValue("value2")},
		}
		record := slog.Record{
			Level:   slog.LevelInfo,
			Message: "Test message",
		}

		err := handler.Handle(ctx, record)
		if err != nil {
			t.Errorf("Handler.Handle returned an error: %v", err)
		}

		if !mockHandler.HandleCalled {
			t.Error("MockHandler.Handle was not called")
		}

		if mockHandler.HandleRecord.Message != record.Message || mockHandler.HandleRecord.Level != record.Level {
			t.Error("MockHandler.Handle did not receive the correct record")
		}

		if mockHandler.HandleRecord.NumAttrs() != len(expectedAttrs) {
			t.Error("MockHandler.WithAttrs did not receive the correct number of attributes")
		}
	})
}

type ComplexData struct {
	IntField   int
	StrField   string
	BoolField  bool
	SliceField []string
}

func ExampleHandler() {
	slog.SetDefault(slog.New(ctxslog.NewHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))))

	ctx := ctxslog.WithValue(context.Background(), "AccountID", 123456789)
	ctx = ctxslog.WithValue(ctx, "email", "noone@felttrip.com")
	ctx = ctxslog.WithValue(ctx, "sender", "greg@BailysInAShoe.lake")

	slog.InfoContext(ctx, "Info With Context")
	ctx = ctxslog.WithValues(context.Background(), map[string]interface{}{
		"AccountID": 987654321,
		"email":     "bob@TheBuilder.fake",
		"complexData": ComplexData{
			IntField:   123,
			StrField:   "DEADBEEF",
			BoolField:  true,
			SliceField: []string{"one", "two", "three"},
		},
	})

	slog.ErrorContext(ctx, "Error With Context")
	// Output:
	//{"level":"INFO","msg":"Info With Context","AccountID":123456789,"email":"noone@felttrip.com","sender":"greg@BailysInAShoe.lake"}
	//{"level":"ERROR","msg":"Error With Context","AccountID":987654321,"email":"bob@TheBuilder.fake","complexData":{"IntField":123,"StrField":"DEADBEEF","BoolField":true,"SliceField":["one","two","three"]}}

}

func TestHandlerWithAttrs(t *testing.T) {
	mockHandler := &MockHandler{}
	handler := ctxslog.NewHandler(mockHandler)

	attrs := []slog.Attr{
		{Key: "attr1", Value: slog.StringValue("value1")},
		{Key: "attr2", Value: slog.StringValue("value2")},
	}

	newHandler := handler.WithAttrs(attrs)

	if newHandler == nil {
		t.Error("Handler.WithAttrs returned nil")
	}

	if len(mockHandler.WithAttrsAttrs) != len(attrs) {
		t.Error("MockHandler.WithAttrs did not receive the correct number of attributes")
	}

	for i, attr := range attrs {
		if !mockHandler.WithAttrsAttrs[i].Equal(attr) {
			t.Errorf("MockHandler.WithAttrs did not receive the correct attribute at index %d", i)
		}
	}
}

func TestHandlerWithGroup(t *testing.T) {
	mockHandler := &MockHandler{WithGroupCalled: false}
	handler := ctxslog.NewHandler(mockHandler)
	groupName := "test-group"

	newHandler := handler.WithGroup(groupName)

	if newHandler == nil {
		t.Error("Handler.WithGroup returned nil")
	}

	if mockHandler.WithGroupCalled != true {
		t.Error("MockHandler.WithGroup was not called")
	}

	if mockHandler.WithGroupName != groupName {
		t.Error("MockHandler.WithGroup did not receive the correct group name")
	}
}

// MockHandler is a mock implementation of slog.Handler for testing purposes.
type MockHandler struct {
	HandleCalled    bool
	HandleRecord    slog.Record
	WithAttrsAttrs  []slog.Attr
	WithGroupCalled bool
	WithGroupName   string
}

func (m *MockHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return true
}

func (m *MockHandler) Handle(ctx context.Context, r slog.Record) error {
	m.HandleCalled = true
	m.HandleRecord = r
	return nil
}

func (m *MockHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	m.WithAttrsAttrs = attrs
	return m
}

func (m *MockHandler) WithGroup(name string) slog.Handler {
	m.WithGroupCalled = true
	m.WithGroupName = name
	return m
}
