---
title: "24. Error handling"
---

# 24. Error handling

| Laravel | Go |
|---|---|
| `throw`, `try` / `catch`, the exception handler | Errors are values: `return err`, `if err != nil`, `errors.Is` |

This is the biggest mental shift for a Laravel developer, and every chapter has already leaned on it: `if err != nil` follows almost every call in this codebase. This chapter adds one more example, recording a result twice, and pulls the patterns together. Go's own [error handling](https://go.dev/blog/error-handling-and-go) and [wrapping](https://go.dev/blog/go1.13-errors) posts are the best reference.

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch24/RecordResult.php [Throwing]

<<< @/snippets/laravel/ch24/handling.php [Handling]

:::

## The Go way

### Errors are returned

A Go function that can fail returns an `error` as its last value. The caller checks it immediately:

```go
var club models.Club
if err := facades.Orm().Query().Create(&club); err != nil {
	return serverError(ctx, err)
}
```

There is no `try`. Nothing unwinds the stack for you, so every call that can fail is followed by a decision about what to do. That is verbose, and it is the point: every failure path is visible in the code.

### Sentinel errors

A sentinel is a package-level error value that callers can recognise, the closest thing to a specific exception class:

<<< @/../app/results/record.go#sentinel{go}

### Wrapping with `%w`

<<< @/../app/results/record.go#wrap{go}

`fmt.Errorf` with `%w` builds a new error whose message adds context and which still contains the original. Each layer can add a little more, like a stack trace written by hand: `sync: record fixture 7: fixture already has a result`.

### Recognising an error: `errors.Is`

<<< @/../app/http/controllers/fixture_controller.go#is{go}

`errors.Is` walks the wrapped chain looking for the sentinel. A plain `==` would only compare the outermost error and miss it. For errors that carry data, `errors.As` does the same walk and extracts a typed value, like catching an exception and reading its properties. The framework uses this too: chapter 04 checks `errors.Is(err, frameworkerrors.OrmRecordNotFound)` to turn a missing row into a 404.

### There is no global handler

Laravel lets exceptions bubble up to one handler that renders them. In Go, the place that receives an error decides. This project's pattern:

| Situation | Where it is handled | Response |
|---|---|---|
| Expected failure (not found, already recorded, invalid input) | The controller, with `errors.Is` or the validator | 404, 409, 422 with a clear message |
| Unexpected failure (database down, a bug) | The controller passes it to `serverError` | Logged with the request (chapter 22), generic 500 |
| A `panic` (nil pointer, index out of range) | Goravel's recover middleware | Logged, bare 500 |

### Panics are not exceptions

`panic` exists, and Goravel's HTTP driver recovers from panics in handlers so one bad request does not crash the server: it logs the panic with the request and aborts with a 500. But panics are for programmer errors, such as reading a field on a `nil` pointer (chapter 10's missing `With`), not for control flow. Code in this project never panics on purpose, except the user factory, which cannot do anything useful if hashing fails.

## Testing it

<<< @/../tests/unit/errors_test.go#test{go}

<<< @/../tests/feature/errors_test.go#test{go}

## What is different and why

- **Errors are ordinary values.** An `error` is an interface with one method, `Error() string`. You return it, store it, compare it and wrap it like any other value.
- **Handle or return, never ignore.** The two honest options are to deal with an error or pass it up with context. Assigning it to `_` is a decision too, and should be rare and obvious. `go vet` does not catch an ignored error; linters such as `errcheck` do.
- **Add context as you return.** `fmt.Errorf("record fixture %d: %w", id, err)` is the habit to build. Without it, the log only says "record not found" with no clue which call failed.
- **Expected failures are not exceptional.** A duplicate result is part of normal operation, so it is a sentinel the caller checks for, not a crash. That is also how the framework reports a missing record.
- **The compiler does not force you to check.** Go makes errors visible, not mandatory: calling `facades.Orm().Query().Save(&club)` and ignoring what it returns compiles. The convention, code review and linters do the enforcing.
