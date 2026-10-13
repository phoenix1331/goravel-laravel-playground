---
title: "01. Installation and structure"
---

# 01. Installation and structure

| Laravel | Goravel |
|---|---|
| `laravel new`, `php artisan`, `composer.json` | `goravel new`, `./artisan`, `go.mod` |

Goravel's layout will look familiar within seconds. The big difference is underneath: Laravel is interpreted per request, Goravel is compiled once into a single binary.

## The Laravel way

```bash
composer global require laravel/installer
laravel new club-app
cd club-app
php artisan serve        # or: php artisan dev
```

Every request runs `public/index.php`, which loads the autoloader, builds the application from `bootstrap/app.php`, handles one request and throws everything away.

::: code-group

<<< @/snippets/laravel/ch01/index.php [public/index.php]

<<< @/snippets/laravel/ch01/app.php [bootstrap/app.php]

:::

## The Goravel way

```bash
go install github.com/goravel/installer/goravel@latest
goravel new club-app
cd club-app
go run .                 # or: air, for live reload
```

This repo was created by cloning the full scaffold instead, which is equivalent:

```bash
git clone --depth=1 https://github.com/goravel/goravel.git club-app
cd club-app && go mod tidy && cp .env.example .env
./artisan key:generate
```

The entry points line up one to one:

::: code-group

<<< @/snippets/laravel/ch01/index.php [Laravel: public/index.php]

<<< @/../main.go#main{go} [Goravel: main.go]

:::

::: code-group

<<< @/snippets/laravel/ch01/app.php [Laravel: bootstrap/app.php]

<<< @/../bootstrap/app.go#boot{go} [Goravel: bootstrap/app.go]

:::

Both use a fluent builder: `Application::configure()->withRouting()->create()` against `foundation.Setup().WithRouting().Create()`. The difference is that Laravel's `withRouting` takes file paths it will `require`, while Goravel's takes a Go function that is compiled into the binary.

## Artisan

```bash
php artisan migrate          # Laravel
./artisan migrate            # Goravel
go run . artisan migrate     # what ./artisan actually runs
```

`./artisan` is a two-line shell script. It compiles the whole app and passes `artisan migrate` as arguments, and the binary decides whether to start the HTTP server or run a command. That is why the first artisan call after a code change takes a few seconds: it is a build.

In this repo, run everything through the Makefile so it happens inside the container: `make artisan cmd="route:list"`.

## Directory structure

| Laravel | Goravel | Notes |
|---|---|---|
| `public/index.php` | `main.go` | Entry point |
| `bootstrap/app.php` | `bootstrap/app.go` | Application builder |
| `bootstrap/providers.php` | `bootstrap/providers.go` | Provider list, as Go values not class names |
| `config/*.php` | `config/*.go` | Config is Go code ([chapter 02](/chapters/02-configuration)) |
| `app/Http/Controllers` | `app/http/controllers` | Lower-case folders: Go package names are lower case |
| `app/Models` | `app/models` | |
| `routes/web.php` | `routes/web.go` | |
| `database/migrations` | `database/migrations` | Migrations are Go structs, registered in `bootstrap/migrations.go` |
| `tests/Feature` | `tests/feature` | |
| `app/Facades` (rare) | `app/facades` | Every facade lives here ([chapter 06](/chapters/06-facades)) |
| `composer.json` | `go.mod` | Dependencies |
| `composer.lock` | `go.sum` | Checksums, not versions (see below) |
| `vendor/` | none | Modules live in a shared cache |

## go.mod and go.sum

`go.mod` is `composer.json`:

```go
module goravel           // the import path, like "App\\" in PSR-4 autoload

go 1.26.9                // minimum Go version, like "php": "^8.3"

require (
	github.com/goravel/framework v1.18.0
	github.com/goravel/gin v1.18.0
	github.com/goravel/mysql v1.18.0
	// ...
)
```

The module name is why imports look like `goravel/app/facades`: everything in the repo is imported by path under the module name.

| Composer | Go |
|---|---|
| `composer require vendor/pkg` | `go get github.com/vendor/pkg` |
| `composer install` | `go mod download` (or just build; Go fetches what it needs) |
| Clean up unused packages | `go mod tidy` |
| Packagist | None: packages are named by repository URL |
| Version ranges like `^1.18` resolved at install | Exact versions in `go.mod`, chosen by minimal version selection |

`go.sum` is not a lock file in Composer's sense. Versions are already exact in `go.mod`. `go.sum` holds a checksum for every module, so a changed or tampered download fails the build. Commit both, never edit `go.sum` by hand.

## Live reload with air

There is nothing to "refresh" in a compiled app. A code change means rebuild and restart. [air](https://github.com/air-verse/air) watches files and does both, configured by `.air.toml`:

```toml
[build]
  cmd = "go build -o ./tmp/main.exe ."
  bin = "./tmp/main.exe"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["storage", "tmp", "docs", "node_modules", ".claude"]
```

In this repo air is the `app` container's command, so `make up` gives you live reload. `database/` is watched on purpose: migrations are Go code here, so changing one needs a rebuild.

## What is different and why

- **One process, not one per request.** The binary boots once and keeps the app in memory. Anything you store in a package-level variable is shared across every request, which is fast but needs care. PHP's "share nothing" model gives you that isolation for free.
- **Compile errors happen before anything runs.** A typo in a controller stops `go build`, so air shows the error instead of the app restarting. In PHP you would find it when that route is hit.
- **Wiring is compiled in.** Laravel discovers routes, providers and commands from files at runtime. Goravel's lists are Go code (`bootstrap/providers.go`, `bootstrap/migrations.go`) that the compiler checks. Forgetting to register something is a common first mistake.
- **Deploying is copying one file.** `go build` produces a static binary with no PHP runtime, extensions or `vendor/` folder to ship.
