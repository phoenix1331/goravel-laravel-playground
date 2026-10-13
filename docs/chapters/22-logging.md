---
title: "22. Logging"
---

# 22. Logging

| Laravel | Goravel |
|---|---|
| `Log::channel('results')->info('...', [context])` | `facades.Log().Channel("results").With(context).Info("...")` |

Two changes in this chapter: every recorded result now writes a structured line to its own channel, and unexpected server errors are logged with the request that caused them. Channels, drivers and levels work as in Laravel; see Goravel's [logging docs](https://www.goravel.dev/the-basics/logging.html) for the full list.

## The Laravel way

<<< @/snippets/laravel/ch22/logging.php

## The Goravel way

A `results` channel in `config/logging.go`, writing JSON:

<<< @/../config/logging.go#channel{go}

Context goes in with `With`, before the level method:

<<< @/../app/results/record.go#log{go}

Each call writes one line to `storage/logs/results.log`:

```json
{"environment":"local","extra":{"away_score":2,"fixture_id":1,"home_score":4},"level":"info","message":"result recorded","time":"2026-10-13 20:03:34.021"}
```

The 500 helper from chapter 04 now attaches the request too:

<<< @/../app/http/controllers/club_controller.go#server-error{go}

The default channel is still `stack`, which writes to a daily file (`storage/logs/goravel-YYYY-MM-DD.log`), so `make logs` does not show these entries. Set `"print": true` on a channel to also echo it to the console.

## Testing it

::: code-group

<<< @/snippets/laravel/ch22/ResultLoggingTest.php [Laravel]

<<< @/../tests/feature/logging_test.go#test{go} [Goravel]

:::

The Go test reads the log file back and decodes the JSON line. `.env.testing` points the channel at `results-testing.log`, so tests never write into the development log.

The first version of this test deleted the file in `SetupTest` and failed, but only when the whole suite ran. An earlier test had already written to the channel, and the logger keeps its file open for the life of the process. Deleting the file left the logger writing to a file that no longer had a name. The test now notes the file size before acting and reads only what was added. A long-running Go process behaves like this in production too, so rotate logs with `copytruncate` or let the `daily` driver do it.

## What is different and why

- **Context is a method, not a second argument.** `->info('message', $context)` becomes `.With(context).Info("message")`. `Info` takes `...any` and joins the values like `fmt.Print`, so passing a map as a second argument would just print it into the message.
- **Fields land under `extra`.** Goravel's JSON formatter nests your context under `extra`, beside `level`, `message`, `time` and `environment`. Monolog's `JsonFormatter` uses `context`.
- **More built-in fields.** `Request()`, `Response()`, `User()`, `Tags()`, `Code()`, `Hint()` and `WithTrace()` attach common context without building a map yourself.
- **One process, open files.** PHP opens and closes the log on every request. A Go app opens it once, which is faster but means deleting or moving the file underneath it loses writes.
- **`Infof` for formatted strings.** `Infof("fixture %d recorded", id)` is `Info(sprintf(...))`. Prefer `With` plus a fixed message when the log will be searched or parsed.
