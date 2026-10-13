---
title: "02. Configuration"
---

# 02. Configuration

| Laravel | Goravel |
|---|---|
| `.env` and `config/*.php`, read with `config()` | `.env` and `config/*.go`, read with `facades.Config()` |

The shape is identical: `.env` holds per-environment values, config files give them names and defaults, and the app reads the named values. The difference is that Goravel's config files are Go code that runs at start-up, and lookups are typed.

## The Laravel way

A config file returns an array. `env()` reads `.env`, and you cast if you care about types:

<<< @/snippets/laravel/ch02/club.php

Read it anywhere with `config('club.points.win')`, which returns `mixed`.

## The Goravel way

A config file is a Go file in the `config` package whose `init()` function registers a section:

<<< @/../config/club.go#club{go}

Read it with a typed getter:

```go
facades.Config().GetString("club.league_name")
facades.Config().GetInt("club.points.win")     // int, converted for you
facades.Config().GetBool("app.debug")
facades.Config().Get("club.points")           // any, when you need the raw value
```

The `.env` side is the same as Laravel:

```ini
LEAGUE_NAME="Sunday League Division One"
LEAGUE_POINTS_WIN=3
```

## Testing it

::: code-group

<<< @/snippets/laravel/ch02/ConfigTest.php [Laravel]

<<< @/../tests/unit/config_test.go#test{go} [Goravel]

:::

`.env.testing` sets `LEAGUE_NAME="Test League"`, and `make test` passes `--env=.env.testing`, the same job `phpunit.xml` or `.env.testing` does in Laravel.

## What is different and why

- **`init()` instead of `return [...]`.** Go has no "require a file and take its return value". Instead, every file in a package can have an `init()` function that Go runs automatically when the package is first imported. `bootstrap/app.go` calls `config.Boot()`, an empty function that exists only to make sure the `config` package is imported, which triggers every `init()` in it.
- **Typed getters do the casting.** `config()` returns `mixed`, so Laravel code often casts with `(int)`. `GetInt`, `GetString`, `GetBool` and `GetDuration` convert at the call site and always return that type. A value that cannot be converted returns the zero value (`0`, `""`, `false`) rather than failing.
- **Real environment variables win.** Goravel's config reader (viper) checks OS environment variables before the `.env` file. That is why this repo's `docker-compose.yml` does not use `env_file`: it would turn `.env` into OS variables that silently override `.env.testing` in tests. It is also how CI switches the test suite to MySQL without a second env file.
- **No `config:cache`.** Laravel caches config to skip parsing files on every request. A Go app reads config once at start-up and keeps it in memory, so there is nothing to cache. The flip side: changing `.env` needs a restart (`make run`).
- **Config is code, so it can hold code.** Look at `config/database.go`: connection entries include a `via` function that returns the driver. Laravel config arrays can hold closures too, but `config:cache` cannot serialise them, so it is rare. In Goravel it is normal.
