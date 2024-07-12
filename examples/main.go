package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/felttrip/ctxslog"
)

type ComplexData struct {
	IntField   int
	StrField   string
	BoolField  bool
	SliceField []string
}

func main() {
	slog.SetDefault(slog.New(ctxslog.NewHandler(slog.NewJSONHandler(os.Stdout, nil))))

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
}
