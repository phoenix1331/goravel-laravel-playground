# Goravel for Laravel devs

A side-by-side guide for Laravel developers learning Go. It maps 24 Laravel concepts to their Goravel equivalents and is backed by a working grassroots football club API.

**Docs site:** https://phoenix1331.github.io/goravel-for-laravel-devs

## Who it is for

- Laravel developers curious about Go who want to start from something familiar
- Anyone evaluating Goravel against Laravel for a backend service

Each chapter shows the Laravel way, the Goravel way, and what is different and why: explicit errors, structs instead of classes, no magic methods, compile-time wiring. The Go snippets on the site are pulled from this repo's source at build time, so they are real, tested code.

## Stack

| Concern | Choice |
|---|---|
| Language | Go 1.26 |
| Framework | [Goravel](https://www.goravel.dev) 1.18 with the Gin HTTP driver |
| Database | MySQL 8.4 (dev), SQLite in memory (tests) |
| Cache / queue | Redis |
| Runtime | Docker Compose, with [air](https://github.com/air-verse/air) for live reload |
| Docs | VitePress on GitHub Pages |
| CI | GitHub Actions: gofmt, go vet, tests on SQLite and MySQL |

## Quick start

### Prerequisites

- Docker with Docker Compose
- Node.js (for the docs site and git hooks); run `npm install` once to activate the Husky hooks
- [osv-scanner](https://github.com/google/osv-scanner) for the pre-commit audit: `go install github.com/google/osv-scanner/cmd/osv-scanner@latest`
- Claude Code MCP servers (optional, configured in `.mcp.json`):
  - GitHub MCP authenticates via `gh auth token`; run `gh auth login` once per machine
  - MySQL MCP reads `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_DATABASE`. For this project export `MYSQL_HOST=127.0.0.1 MYSQL_PORT=3307 MYSQL_USER=goravel MYSQL_PASSWORD=secret MYSQL_DATABASE=goravel`
  - Playwright and Context7 need only `npx`

### Run it

```bash
cp .env.example .env
make up                          # app on http://localhost:8000, MySQL on 3307
make artisan cmd="key:generate"  # sets APP_KEY in .env
make run                         # restart the app with the new key
make migrate
make test                        # SQLite in memory
make test-mysql                  # same suite on MySQL, to catch dialect drift
make docs-dev                    # docs on http://localhost:5173
```

`make help` lists every target. Try `GET http://localhost:8000/health`, or open `request.http` in VS Code REST Client or a JetBrains IDE.

## Chapters

| # | Chapter | Laravel | Goravel |
|---|---|---|---|
| 01 | [Installation and structure](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/01-installation) | `laravel new`, `artisan`, directory structure | `goravel new`, `./artisan`, layout |
| 02 | [Configuration](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/02-configuration) | `.env` and `config/*.php` | `.env`, `config/*.go`, `facades.Config()` |
| 03 | [Routing](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/03-routing) | `routes/web.php`, `routes/api.php` | `routes/*.go`, route groups |
| 04 | [Controllers](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/04-controllers) | Controller classes | Controller structs |
| 05 | [Service container and providers](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/05-service-container) | Service container, `register()` / `boot()` | Container bindings, `Register` / `Boot` |
| 06 | [Facades](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/06-facades) | Static proxies via `__callStatic` | `app/facades` package functions |
| 07 | [Form requests and validation](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/07-validation) | Form requests, `$request->validate()` | Form requests, validation rules |
| 08 | [Middleware](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/08-middleware) | Middleware classes with `handle()` | Middleware functions |
| 09 | [ORM models](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/09-models) | Eloquent models | ORM models |
| 10 | [Relationships](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/10-relationships) | Eloquent relationships, `with()` | ORM relationships, `With()` |
| 11 | [Migrations](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/11-migrations) | Migrations with the schema builder | Migrations in Go |
| 12 | [Factories and seeders](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/12-factories-seeders) | Model factories and seeders | Factories and seeders |
| 13 | [Query builder and raw SQL](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/13-query-builder) | `DB::table()`, `DB::select()` | `facades.Orm().Query()`, `facades.DB()` |
| 14 | [Artisan commands](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/14-console-commands) | `make:command`, `handle()` | Console commands (`make:command`) |
| 15 | [Events and listeners](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/15-events) | Events, listeners, `event()` | Events and listeners |
| 16 | [Jobs and queues](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/16-queues) | Jobs, `dispatch()`, queue workers | Jobs and queues |
| 17 | [Task scheduling](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/17-scheduling) | `routes/console.php` schedule | Scheduler |
| 18 | [Cache](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/18-cache) | `Cache::remember()` | Cache facade |
| 19 | [Authentication](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/19-authentication) | Session guard, Sanctum | Auth with the JWT guard |
| 20 | [Gates and policies](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/20-authorisation) | Gates and policies | Authorisation |
| 21 | [HTTP client](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/21-http-client) | `Http::get()`, `Http::fake()` | HTTP client |
| 22 | [Logging](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/22-logging) | `Log::info()`, channels | Logging |
| 23 | [Testing](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/23-testing) | PHPUnit / Pest, HTTP tests, fakes | testify suites, HTTP tests, mocks |
| 24 | [Error handling](https://phoenix1331.github.io/goravel-for-laravel-devs/chapters/24-error-handling) | Exceptions, `try` / `catch`, handler | Go errors |

## Licence

MIT
