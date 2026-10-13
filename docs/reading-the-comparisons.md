# Reading the comparisons

Every chapter follows the same pattern, shown here with the smallest feature in the repo: a `/health` endpoint that reports whether the app can reach its database.

- **Laravel tab:** a short PHP file from [`docs/snippets/laravel/`](https://github.com/phoenix1331/goravel-laravel-playground/tree/main/docs/snippets/laravel). These are illustrative and are not run.
- **Goravel tab:** pulled straight from the real source files at build time, so what you read is the code that is tested in CI.

## The route

::: code-group

<<< @/snippets/laravel/ch03/web.php [Laravel]

<<< @/../routes/web.go#health{go} [Goravel]

:::

## The handler

::: code-group

<<< @/snippets/laravel/ch03/HealthController.php [Laravel]

<<< @/../app/http/controllers/health_controller.go#show{go} [Goravel]

:::

Two things a Laravel developer will notice straight away:

1. **No `try` / `catch`.** `facades.Orm().DB()` and `PingContext` return an `error` value alongside their result. You check it with `if err != nil`. Chapter 24 covers why Go works this way.
2. **The response type is explicit.** The handler signature says it takes an `http.Context` and returns an `http.Response`. There is no magic conversion of arrays to JSON.

## The test

::: code-group

<<< @/snippets/laravel/ch23/HealthTest.php [Laravel]

<<< @/../tests/feature/health_test.go#test{go} [Goravel]

:::

The Go test runs against SQLite in memory. Chapter 23 explains the harness, including why the DSN is `file:/goravel_test.db?vfs=memdb` and not `:memory:`.

## Finding the Go source

Key Go files start with a comment naming their Laravel equivalent, for example:

```go
// Laravel equivalent: tests/TestCase.php - see docs/chapters/23-testing.md
```

So you can start from either side: the docs or the code.
