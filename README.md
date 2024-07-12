# ctxslog
An [slog](https://pkg.go.dev/log/slog) handler to pull shared information from `context.Context` for use with structured logging.

# Installation 
```
go get github.com/felttrip/ctxslog
```

# Usage
* Add `ctxslog` into the handler chain when calling `SetDefault` with `slog.
* Use `ctxslog.WithValue` or `ctxslog.WithValues` to add Key Value pairs to the context to be logged
  * This operates the same way as `context.Context.WithValue` but with the addition of keeping track of which fields on the context
    should be logged.
* Use the `slog.InfoContext`, `slog.WarnContext`, or `slog.ErrorContext` functions providing the context that has the fields you want to log attached. 
```
slog.SetDefault(slog.New(ctxslog.NewHandler(slog.NewJSONHandler(os.Stdout, nil))))

ctx := ctxslog.WithValue(context.Background(), "AccountID", 123456789)
ctx = ctxslog.WithValue(ctx, "email", "noone@felttrip.com")
ctx = ctxslog.WithValue(ctx, "sender", "greg@BailysInAShoe.lake")

slog.InfoContext(ctx, "Info With Context")
fmt.Println()
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
```

Example Output
```
{"time":"2024-07-12T10:38:41.610863-06:00","level":"INFO","msg":"Info With Context","AccountID":123456789,"email":"noone@felttrip.com","sender":"greg@BailysInAShoe.lake"}

{"time":"2024-07-12T10:38:41.611199-06:00","level":"ERROR","msg":"Error With Context","complexData":{"IntField":123,"StrField":"DEADBEEF","BoolField":true,"SliceField":["one","two","three"]},"AccountID":987654321,"email":"bob@TheBuilder.fake"}
```