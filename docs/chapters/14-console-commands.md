---
title: "14. Artisan commands"
---

# 14. Artisan commands

| Laravel | Goravel |
|---|---|
| `make:command`, `$signature`, `handle(): int` | `make:command`, `Signature()` + `Extend()`, `Handle(ctx) error` |

A Goravel command is a struct with four methods: name, description, options and the code to run. The options are declared as Go values instead of being parsed out of a signature string.

The example is `sync:fixtures`, which reads results from the league's feed and records them. Here the feed is a JSON file; [chapter 21](/chapters/21-http-client) swaps in an HTTP API, and [chapter 17](/chapters/17-scheduling) runs the command on a schedule.

```bash
make artisan cmd="sync:fixtures --dry-run"
```

```
╭──────────────┬───────────────┬───────┬────────────────────────────────────╮
│ Home         │ Away          │ Score │ Outcome                            │
├──────────────┼───────────────┼───────┼────────────────────────────────────┤
│ Ilkley Town  │ Crewe United  │ 2-2   │ would record                       │
│ Ripon Rovers │ Ilkley Rovers │ 0-1   │ would record                       │
│ Nowhere FC   │ Ripon Rovers  │ -     │ skipped: unknown club "Nowhere FC" │
╰──────────────┴───────────────┴───────┴────────────────────────────────────╯
```

## The Laravel way

<<< @/snippets/laravel/ch14/SyncFixtures.php

## The Goravel way

`make artisan cmd="make:command SyncFixtures"` creates the struct and registers it in `bootstrap/commands.go`.

<<< @/../app/console/commands/sync_fixtures.go#signature{go}

<<< @/../app/console/commands/sync_fixtures.go#handle{go}

`newFeed` arrived in [chapter 21](/chapters/21-http-client) to choose between a file and an HTTP API. In this chapter the command built a `services.FileFeed{Path: path}` directly.

The feed is an interface, so the command does not care where results come from:

<<< @/../app/services/fixture_feed.go#feed{go}

### One place to record a result

The command does not update the fixture itself. It calls `results.Record`, a small function playing the part of a Laravel action class:

```go
// Record marks a fixture as finished with the given score. Every way a result
// arrives (the sync command, the API) goes through here, so later chapters
// have one place to hook in.
func Record(fixture *models.Fixture, homeScore, awayScore int) error {
	if homeScore < 0 || awayScore < 0 {
		return fmt.Errorf("record fixture %d: scores cannot be negative", fixture.ID)
	}

	fixture.Status = models.Finished
	fixture.HomeScore = &homeScore
	fixture.AwayScore = &awayScore

	if err := facades.Orm().Query().Save(fixture); err != nil {
		return fmt.Errorf("record fixture %d: %w", fixture.ID, err)
	}

	return nil
}
```

Chapter 15 adds an event to this function, and chapter 16 a queued job, without touching the command.

## Testing it

::: code-group

<<< @/snippets/laravel/ch14/SyncFixturesTest.php [Laravel]

<<< @/../tests/feature/sync_fixtures_test.go#test{go} [Goravel]

:::

`facades.Artisan().Call("sync:fixtures --file=...")` runs the command in the test process, like `$this->artisan(...)`, and returns the error from `Handle`.

The test feed lives in `tests/feature/testdata/`. Go runs each package's tests with the working directory set to that package's folder, so the relative path works, and the Go tools ignore any directory called `testdata` when building.

## What is different and why

- **Options are values, not a mini-language.** `{--file=} {--dry-run}` becomes a `StringFlag` and a `BoolFlag` in `Extend()`. More typing, but the compiler checks it, and `OptionBool("dry-run")` returns a real `bool`.
- **Return an error, not an exit code.** `Handle` returns `error`. `nil` exits 0; any error exits non-zero and is printed. There is no `self::SUCCESS` or `self::FAILURE`.
- **Registration is explicit.** Laravel discovers commands in `app/Console/Commands`. Goravel lists them in `bootstrap/commands.go`; `make:command` adds the entry.
- **No method injection.** Laravel resolves `handle(RecordResult $record)` from the container. Here `Handle` calls `results.Record` and builds a `FileFeed` directly. Chapter 21 shows how to pick the feed from config instead.
- **Every artisan call compiles the app.** `./artisan sync:fixtures` is `go run . artisan sync:fixtures`, so the first run after a change takes a few seconds to build. In production you run the compiled binary: `./app artisan sync:fixtures`.
