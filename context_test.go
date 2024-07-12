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
		val := ctx.Value(ctxKey("key1"))
		//we expect value1 to end up in the context
		if val != "value1" {
			t.Errorf("Expected value1, got %v", val)
		}

		loggedFields.mu.Lock()
		//we expect the key to be in the loggedFields
		if loggedFields.fields[0] != "key1" {
			t.Errorf("Expected key1, got %v", loggedFields.fields[0])
		}
		loggedFields.mu.Unlock()
		ctx.Done()
	})

	t.Run("Will Store Multiple Values On Context", func(t *testing.T) {
		ctx := context.Background()
		defer ctx.Done()

		ctx = WithValue(ctx, "key2", "value2")
		ctx = WithValue(ctx, "key3", "value3")

		// we expect value2 and value3 to end up in the context
		val := ctx.Value(ctxKey("key2"))
		if val != "value2" {
			t.Errorf("Expected value2, got %v", val)
		}
		val = ctx.Value(ctxKey("key3"))
		if val != "value3" {
			t.Errorf("Expected value3, got %v", val)
		}
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

		// we expect all the fields to end up in the context
		val := ctx.Value(ctxKey("AccountID"))
		if val != 987654321 {
			t.Errorf("Expected 987654321, got %v", val)
		}
		val = ctx.Value(ctxKey("email"))
		if val != "bob@TheBuilder.fake" {
			t.Errorf("Expected bob@TheBuilder.fak, got %v", val)
		}
		complexVal, ok := ctx.Value(ctxKey("complexData")).(ComplexData)
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
	})

	t.Run("Will panic with nil parent", func(t *testing.T) {
		defer func() { _ = recover() }()

		_ = WithValues(nil, map[string]interface{}{})
		t.Errorf("did not panic")
	})
}
