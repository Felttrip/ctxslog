package ctxslog

import (
	"context"
	"testing"
)

type ComplexData struct {
	IntField   int
	StrField   string
	BoolField  bool
	SliceField []string
}

func TestWithValue(t *testing.T) {
	t.Run("Will Store Single Value On Context", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()

		ctx = WithValue(ctx, "key1", "value1")
		m, ok := ctx.Value(fieldsKey).(*safeMap)
		if !ok {
			t.Errorf("Expected safeMap to be initialized in context")
		}
		//we expect value1 to end up in the contexts map
		m.mu.Lock()
		if m.fields["key1"] != "value1" {
			t.Errorf("Expected value1, got %v", m.fields["key1"])
		}
		m.mu.Unlock()
	})

	t.Run("Will Store Multiple Values On Context", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()

		ctx = WithValue(ctx, "key2", "value2")
		ctx = WithValue(ctx, "key3", "value3")

		m, ok := ctx.Value(fieldsKey).(*safeMap)
		if !ok {
			t.Errorf("Expected safeMap to be initialized in context")
		}
		//we expect value2 and value3 to end up in the contexts map
		m.mu.Lock()
		if m.fields["key2"] != "value2" {
			t.Errorf("Expected value2, got %v", m.fields["key2"])
		}
		if m.fields["key3"] != "value3" {
			t.Errorf("Expected value3, got %v", m.fields["key3"])
		}
		m.mu.Unlock()
	})

	t.Run("Will panic with nil parent", func(t *testing.T) {
		defer func() { _ = recover() }()

		_ = WithValue(nil, "", "")
		t.Errorf("did not panic")
	})
}

func TestWithValues(t *testing.T) {
	t.Run("Will Store Multiple Values On Context", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()
		cd := ComplexData{
			IntField:   123,
			StrField:   "DEADBEEF",
			BoolField:  true,
			SliceField: []string{"one", "two", "three"},
		}
		ctx = WithValues(ctx, map[string]interface{}{
			"AccountID":   987654321,
			"email":       "bob@TheBuilder.fake",
			"complexData": cd,
		})

		m, ok := ctx.Value(fieldsKey).(*safeMap)
		if !ok {
			t.Errorf("Expected safeMap to be initialized in context")
		}
		//we expect value1 to end up in the contexts map
		m.mu.Lock()
		if m.fields["AccountID"] != 987654321 {
			t.Errorf("Expected 987654321, got %v", m.fields["AccountID"])
		}
		if m.fields["email"] != "bob@TheBuilder.fake" {
			t.Errorf("Expected bob@TheBuilder.fake, got %v", m.fields["email"])
		}
		complexVal, ok := m.fields["complexData"].(ComplexData)
		if !ok {
			t.Errorf("mistatch type when retrieving ComplexData from context")
		}
		if complexVal.IntField != 123 {
			t.Errorf("Expected 123, got %v", complexVal.IntField)
		}
		if complexVal.StrField != "DEADBEEF" {
			t.Errorf("Expected 123, got %v", complexVal.StrField)
		}
		if complexVal.BoolField != true {
			t.Errorf("Expected 123, got %v", complexVal.BoolField)
		}
		if complexVal.SliceField[0] != "one" || complexVal.SliceField[1] != "two" || complexVal.SliceField[2] != "three" {
			t.Errorf("Expected []string{\"one\", \"two\", \"three\"}, got %v", complexVal.SliceField)
		}
		m.mu.Unlock()
	})

	t.Run("Will panic with nil parent", func(t *testing.T) {
		defer func() { _ = recover() }()

		_ = WithValues(nil, map[string]interface{}{})
		t.Errorf("did not panic")
	})
}
