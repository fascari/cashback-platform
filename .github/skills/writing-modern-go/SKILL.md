---
name: writing-modern-go
description: Use when writing or reviewing Go code to ensure modern idioms are used instead of legacy patterns. Covers Go 1.18 through 1.26 features including slices, maps, cmp, errors.AsType, new(val), and wg.Go.
---

# Writing Modern Go

Always use the most modern Go idiom available up to the project's Go version.
This project uses **Go 1.26.1** — all features below are available.

## Quick Reference

| Legacy pattern | Modern replacement | Since |
|---|---|---|
| `v := 42; &v` | `new(42)` | 1.26 |
| `errors.As(err, &target)` | `errors.AsType[T](err)` | 1.26 |
| `wg.Add(1); go func(){ defer wg.Done(); f() }()` | `wg.Go(f)` | 1.25 |
| `for i := 0; i < n; i++` | `for i := range n` | 1.22 |
| `for _, p := range strings.Split(s, ",")` | `for p := range strings.SplitSeq(s, ",")` | 1.24 |
| `ctx, cancel := context.WithCancel(context.Background())` in tests | `ctx := t.Context()` | 1.24 |
| `for i := 0; i < b.N; i++` in benchmarks | `for b.Loop()` | 1.24 |
| `json:"field,omitempty"` for time/structs/slices/maps | `json:"field,omitzero"` | 1.24 |
| `interface{}` | `any` | 1.18 |
| `if a > b { return a }; return b` | `max(a, b)` | 1.21 |
| `sort.Slice(s, less)` | `slices.SortFunc(s, cmp)` | 1.21 |
| manual loop to check membership | `slices.Contains(s, v)` | 1.21 |
| manual map copy loop | `maps.Clone(m)` | 1.21 |
| chain of `if x == "" { x = fallback }` | `cmp.Or(a, b, c)` | 1.22 |
| `reflect.TypeOf((*T)(nil)).Elem()` | `reflect.TypeFor[T]()` | 1.22 |
| manual loop to collect map keys | `slices.Sorted(maps.Keys(m))` | 1.23 |
| `err == target` | `errors.Is(err, target)` | 1.13 |
| `time.Now().Sub(start)` | `time.Since(start)` | 1.0 |
| `deadline.Sub(time.Now())` | `time.Until(deadline)` | 1.8 |
| `sync.Once` + wrapper function | `sync.OnceFunc(fn)` / `sync.OnceValue(fn)` | 1.21 |
| `errors.New` + manual join | `errors.Join(err1, err2)` | 1.20 |

## Go 1.26 — Critical (newest, most likely missed)

### `new(val)` — pointer to any value

```go
// Before
timeout := 30
cfg := Config{Timeout: &timeout}

// After
cfg := Config{Timeout: new(30)}
```

Type is inferred: `new(0)` → `*int`, `new("s")` → `*string`, `new(T{})` → `*T`.
Never use redundant casts like `new(int(0))`.

### `errors.AsType[T](err)` — type-safe error matching

```go
// Before
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    handle(pathErr)
}

// After
if pathErr, ok := errors.AsType[*os.PathError](err); ok {
    handle(pathErr)
}
```

## Go 1.25

### `wg.Go(fn)` — goroutine with WaitGroup

```go
// Before
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func() {
        defer wg.Done()
        process(item)
    }()
}
wg.Wait()

// After
var wg sync.WaitGroup
for _, item := range items {
    wg.Go(func() {
        process(item)
    })
}
wg.Wait()
```

## Go 1.24

### `t.Context()` — test context (ALWAYS use in tests)

```go
func TestFoo(t *testing.T) {
    ctx := t.Context()
    result := doSomething(ctx)
}
```

### `b.Loop()` — benchmark loop

```go
func BenchmarkFoo(b *testing.B) {
    for b.Loop() {
        doWork()
    }
}
```

### `omitzero` — JSON tag for time, structs, slices, maps

```go
type Config struct {
    Timeout time.Duration `json:"timeout,omitzero"`
    StartAt time.Time     `json:"start_at,omitzero"`
}
```

### `strings.SplitSeq` / `strings.FieldsSeq` — iterate without allocating

```go
for part := range strings.SplitSeq(s, ",") {
    process(part)
}
```

Also: `bytes.SplitSeq`, `bytes.FieldsSeq`.

## Go 1.22–1.23

### `for i := range n` — integer range loop

```go
for i := range len(items) {
    process(items[i])
}
```

### `cmp.Or` — first non-zero value

```go
name := cmp.Or(os.Getenv("NAME"), config.Name, "default")
```

### `slices.Sorted(maps.Keys(m))` — sorted map keys

```go
keys := slices.Sorted(maps.Keys(m))
```

### Enhanced `http.ServeMux`

```go
mux.HandleFunc("GET /api/{id}", handler)
id := r.PathValue("id")
```

## Go 1.21

### Built-in `min` / `max` / `clear`

```go
biggest := max(a, b, c)
clear(myMap)
```

### `slices` package

```go
slices.Contains(items, x)
slices.SortFunc(items, func(a, b T) int { return cmp.Compare(a.X, b.X) })
slices.Index(items, x)
slices.Reverse(items)
slices.Compact(items)
```

### `maps` package

```go
maps.Clone(m)
maps.Copy(dst, src)
maps.DeleteFunc(m, func(k K, v V) bool { return condition })
```

### `sync.OnceFunc` / `sync.OnceValue`

```go
initDB := sync.OnceFunc(func() { db = connect() })
getConfig := sync.OnceValue(func() Config { return loadConfig() })
```

## Go 1.18–1.20

### `any` instead of `interface{}`

### `strings.Cut` / `bytes.Cut`

```go
before, after, found := strings.Cut(s, ":")
```

### `strings.CutPrefix` / `strings.CutSuffix` (1.20)

```go
if rest, ok := strings.CutPrefix(s, "Bearer "); ok {
    token = rest
}
```

### `errors.Join` (1.20)

```go
return errors.Join(err1, err2)
```

### `context.WithCancelCause` / `context.Cause` (1.20)

```go
ctx, cancel := context.WithCancelCause(parent)
cancel(myErr)
// later: context.Cause(ctx) returns myErr
```

### Type-safe atomics (1.19)

```go
var flag atomic.Bool
flag.Store(true)

var ptr atomic.Pointer[Config]
ptr.Store(cfg)
```

## Common Mistakes

| Mistake | Fix |
|---|---|
| Using `new(int(0))` | Just `new(0)` — type is inferred |
| Using `errors.As` when `errors.AsType` is available | Always prefer `errors.AsType[T]` on Go 1.26+ |
| Using `wg.Add(1)` + `go func` + `defer wg.Done()` | Use `wg.Go(fn)` on Go 1.25+ |
| Using `context.Background()` in tests | Use `t.Context()` on Go 1.24+ |
| Using `omitempty` for `time.Time` or structs | Use `omitzero` on Go 1.24+ |
| Using `strings.Split` in `for range` | Use `strings.SplitSeq` on Go 1.24+ |
| Using `for i := 0; i < b.N; i++` | Use `for b.Loop()` on Go 1.24+ |
| Using `sort.Slice` | Use `slices.SortFunc` on Go 1.21+ |
| Using `interface{}` | Use `any` on Go 1.18+ |
