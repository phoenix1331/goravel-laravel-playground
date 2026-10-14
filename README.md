<div align="center">

<a href="https://www.goravel.dev"><img src="https://www.goravel.dev/logo.png" width="300" alt="Goravel"></a>

# Goravel playground for Laravel devs

[Goravel website](https://www.goravel.dev) · [Goravel documentation](https://www.goravel.dev/getting-started/installation.html) · [Goravel on GitHub](https://github.com/goravel/goravel)

</div>

<img width="3010" height="1650" alt="Screenshot 2026-10-13 202051" src="https://github.com/user-attachments/assets/36575dc9-448e-4983-8089-77ae85925e5b" />


> **A playground, not the docs.** [Goravel](https://www.goravel.dev) is built and maintained by the Goravel team and contributors. This is an independent, unofficial project: a small working Goravel app for Laravel developers to clone, run and break. It is not affiliated with or endorsed by the Goravel project, and the Goravel name and logo belong to their owners.
>
> For the full feature-by-feature comparison, read Goravel's own [Compare with Laravel](https://www.goravel.dev/prologue/compare-with-laravel.html) page and the [Goravel documentation](https://www.goravel.dev/getting-started/installation.html).

Something to play with: a grassroots football club API (clubs, players, fixtures, a league table) built with Goravel, where each feature sits next to a short note on how you would do it in Laravel and what Go does differently. Run it, poke the endpoints, change the code and watch the tests.

**Project page:** https://phoenix1331.github.io/goravel-laravel-playground (published once the repository is public)

> **API only.** There is no front end: every endpoint returns JSON. Open http://localhost:8000/ in a browser for a live list of endpoints (JSON clients get the same list as JSON), and use `request.http` to send requests.

## What this demonstrates

Goravel deliberately mirrors Laravel, so a Laravel developer can read most of it on sight. The interesting part is where the resemblance stops:

- **The familiar surface:** routing, controllers, form requests, middleware, an Eloquent-style ORM, migrations, factories, artisan commands, events, queues, scheduling, cache, JWT auth and policies, all in one small app.
- **Where Go changes the rules:** no exceptions (errors are returned values), no magic methods (every model field is declared and typed), no autowiring (dependencies are wired by hand), and no runtime discovery (routes, migrations, commands, jobs and listeners are registered in code the compiler checks).
- **Surprises worth knowing:** SQLite and MySQL disagreeing about timezones, Redis refusing to cache Go structs, Goravel's JWTs not being unique per login, and a scaffold facade that hid the HTTP fakes. Each is noted in the code where it bit.

## How to read it

There are no chapters. Open a file, and the Laravel version is written above the Go in the comments:

```go
// Show returns one club, or 404.
//
// Laravel:
//
//	public function show(Club $club): JsonResponse  // route model binding
//	{
//	    return response()->json(['data' => $club]);
//	}
//
// Goravel: no route model binding, so findClub reads {id} and loads the row.
func (r *ClubController) Show(ctx http.Context) http.Response {
	club, response := findClub(ctx)
	// ...
}
```

The comparisons cover the code you write every day. Framework plumbing (bootstrap, providers, the container) has ordinary comments; the [Goravel documentation](https://www.goravel.dev/getting-started/installation.html) covers it properly.

| Start here | What it compares |
|---|---|
| [`routes/api.go`](routes/api.go) | Route groups, middleware, closures and handlers, no route model binding |
| [`app/http/controllers/club_controller.go`](app/http/controllers/club_controller.go) | A full CRUD controller with no constructor injection |
| [`app/http/requests/store_club_request.go`](app/http/requests/store_club_request.go) | Form requests: familiar rules, typed validated data, the `unique` argument trap |
| [`app/http/middleware/ensure_json_request.go`](app/http/middleware/ensure_json_request.go) | `Next()` and `Abort()` instead of returning `$next($request)` |
| [`app/models/fixture.go`](app/models/fixture.go) | Structs instead of `__get`, relations as fields, types instead of `$casts` |
| [`database/migrations/`](database/migrations) | The schema builder in Go, registered in `bootstrap/migrations.go` |
| [`database/factories/club_factory.go`](database/factories/club_factory.go) | Factories keyed by struct field name |
| [`app/policies/club_policy.go`](app/policies/club_policy.go) | Policies whose gate is not told who the user is |
| [`app/events/`](app/events), [`app/listeners/`](app/listeners), [`app/jobs/`](app/jobs) | Typed arguments instead of serialised models, no `queue:work` |
| [`app/console/commands/sync_fixtures.go`](app/console/commands/sync_fixtures.go) | Flags as values, errors instead of exit codes |
| [`bootstrap/schedule.go`](bootstrap/schedule.go) | The scheduler without a cron entry |
| [`app/results/record.go`](app/results/record.go) | An action that returns a sentinel error instead of throwing |
| [`tests/test_case.go`](tests/test_case.go) | The test harness: in-memory SQLite, `LoginAs` instead of `actingAs` |

For the full feature-by-feature comparison, read Goravel's [Compare with Laravel](https://www.goravel.dev/prologue/compare-with-laravel.html) page.

## Tech stack

| Concern | Choice |
|---|---|
| Language | Go 1.26 |
| Framework | [Goravel](https://www.goravel.dev) 1.18 with the Gin HTTP driver |
| Database | MySQL 8.4 in development, SQLite in memory for tests |
| Cache, queue | Redis (the queue is `sync` in tests) |
| Auth | Goravel's JWT guard |
| Runtime | Docker Compose, with [air](https://github.com/air-verse/air) for live reload |
| Quality | `gofmt`, `go vet`, [osv-scanner](https://github.com/google/osv-scanner), Husky git hooks, GitHub Actions |

## Getting started

### Prerequisites

- Docker with Docker Compose. Port 8000 (app) and 3307 (MySQL) must be free
- Node.js, only for the Husky git hooks: run `npm install` once to activate them (optional if you are just exploring)
- [osv-scanner](https://github.com/google/osv-scanner), used by the pre-commit hook (also optional for exploring): `go install github.com/google/osv-scanner/cmd/osv-scanner@latest`
- Go itself is optional on the host: everything runs inside the `app` container
- Claude Code MCP servers (optional, configured in `.mcp.json`):
  - GitHub MCP authenticates via `gh auth token`; run `gh auth login` once per machine
  - MySQL MCP reads `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_USER`, `MYSQL_PASSWORD` and `MYSQL_DATABASE`. For this project: `export MYSQL_HOST=127.0.0.1 MYSQL_PORT=3307 MYSQL_USER=goravel MYSQL_PASSWORD=secret MYSQL_DATABASE=goravel`
  - Playwright and Context7 need nothing beyond `npx`

### Run it

1. `cp .env.example .env`
2. `make up`: builds and starts the app (http://localhost:8000), MySQL (port 3307) and Redis
3. `make artisan cmd="key:generate"`: sets `APP_KEY`
4. `make artisan cmd="jwt:secret"`: sets `JWT_SECRET` for logins
5. `make run`: restarts the app so it picks up both keys
6. `make fresh`: migrates and seeds a demo league
7. `make test`: runs the suite on SQLite in memory
8. `make test-mysql`: runs the same suite on MySQL and Redis

`make help` lists every target, including `logs`, `shell`, `lint`, `fmt` and `migrate`.

The demo league has six clubs, 96 players, 15 fixtures (seven already played) and one admin per club: `admin1@example.com` to `admin6@example.com`, password `password`. Admin N runs club N.

## Environment variables

Unused keys from the Goravel scaffold (mail, gRPC, AI) are left out.

| Key | Purpose | Default in `.env.example` |
|---|---|---|
| `APP_NAME` | App name; also prefixes cache keys | `ClubManager` |
| `APP_ENV` | Environment name | `local` |
| `APP_KEY` | 32-character encryption key, set by `key:generate` | empty |
| `APP_DEBUG` | Debug mode | `true` |
| `APP_URL`, `APP_HOST`, `APP_PORT` | Where the HTTP server listens | `http://localhost:8000`, `0.0.0.0`, `8000` |
| `APP_DISABLED_RUNNERS` | Space-separated runners to switch off, e.g. `goravel:http` for a worker-only process | empty |
| `JWT_SECRET` | Signs login tokens, set by `jwt:secret` | empty |
| `BCRYPT_ROUNDS` | Password hashing cost (`.env.testing` uses 4) | `12` |
| `LOG_CHANNEL`, `LOG_LEVEL` | Default log channel and level | `stack`, `debug` |
| `DB_CONNECTION` | `mysql` in development, `sqlite` in tests | `mysql` |
| `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD` | MySQL connection (inside Docker) | `mysql`, `3306`, `goravel`, `goravel`, `secret` |
| `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD` | Redis connection | `redis`, `6379`, empty |
| `CACHE_STORE` | `redis` or `memory` | `redis` |
| `QUEUE_CONNECTION` | `redis` or `sync` | `redis` |
| `SESSION_DRIVER`, `SESSION_LIFETIME` | Sessions (unused by the JWT API) | `file`, `120` |
| `LEAGUE_NAME` | Shown by `GET /api/league` | `Sunday League Division One` |
| `LEAGUE_POINTS_WIN`, `LEAGUE_POINTS_DRAW`, `LEAGUE_POINTS_LOSS` | League table points | `3`, `1`, `0` |
| `CLUB_MAX_SQUAD_SIZE` | Largest squad a club may register | `25` |
| `FIXTURE_FEED_DRIVER` | Where `sync:fixtures` reads results: `file` or `http` | `file` |
| `FIXTURE_FEED_PATH` | JSON file for the file feed | `database/data/fixture_results.json` |
| `FIXTURE_FEED_URL` | Base URL for the HTTP feed | `http://localhost:8000/api/demo` |

## Endpoints

http://localhost:8000/ lists every route live. `request.http` has a ready-made request for each one; it works in VS Code (REST Client extension) and JetBrains IDEs. Writes need a token: send the "Log in" request and paste its `access_token` into `@token` at the top of the file.

| Method | Path | Auth | What it does |
|---|---|---|---|
| `GET` | `/`, `/health` | No | Endpoint index; database health check |
| `GET` | `/api/league`, `/api/league/points` | No | League name and points from config |
| `GET` | `/api/league/table` | No | League table (query builder, cached in Redis) |
| `GET` | `/api/league/goals` | No | Goals per club (raw SQL) |
| `GET` | `/api/clubs`, `/api/clubs/{id}` | No | Clubs |
| `GET` | `/api/clubs/{id}/players` | No | A club with its squad |
| `GET` | `/api/fixtures` | No | Fixtures with both clubs eager loaded |
| `GET` | `/api/demo/results` | No | Stand-in results API for `FIXTURE_FEED_DRIVER=http` |
| `POST` | `/api/auth/login` | No | Email and password for a JWT |
| `GET`, `POST` | `/api/auth/me`, `/api/auth/logout` | Yes | Current user; disable the token |
| `POST` | `/api/clubs` | Yes | Create a club; you become its admin (only if you have no club) |
| `PUT`, `PATCH`, `DELETE` | `/api/clubs/{id}` | Yes, that club's admin | Edit or delete a club |
| `POST` | `/api/clubs/{id}/players` | Yes, that club's admin | Add a player |
| `POST` | `/api/fixtures/{id}/result` | Yes, the home club's admin | Record the final score: fires an event, queues notifications, clears the league table cache. A second result returns 409 |

Writes must send `Content-Type: application/json` (415 otherwise), validation failures return 422, and unknown URLs return a JSON 404.

From the command line: `make artisan cmd="sync:fixtures --dry-run"` previews importing results from the feed, and `make artisan cmd="schedule:list"` shows the scheduled tasks.

## Architecture notes

- **One process does everything.** `app.Start()` runs the HTTP server, the Redis queue worker and the scheduler as goroutines in one binary, so there is no `queue:work` and no cron entry. `APP_DISABLED_RUNNERS` splits them apart when you need separate processes. The catch: every copy of the app runs the scheduler unless you use `OnOneServer()`.
- **Tests run on in-memory SQLite, with a twist.** A plain `:memory:` database is per connection, so a second pooled connection sees an empty database. The suite uses the pure-Go driver's `memdb` file system (`file:/goravel_test.db?vfs=memdb`), migrates once per package in `TestMain`, and deletes rows between tests. No CGO is needed.
- **A second, production-like test run.** `make test-mysql` and the CI job rerun the same tests on MySQL and Redis by overriding environment variables. It caught a timezone difference (SQLite keeps a time's offset, MySQL converts to UTC) and a cache difference (Redis cannot store Go structs, the memory store can) that the SQLite run never would.
- **One path for recording results.** The API and the `sync:fixtures` command both call `results.Record`, which saves the score, logs it to a JSON channel and fires `FixtureResultRecorded`. Listeners queue the admin notifications and clear the cached league table. Expected failures are sentinel errors (`errors.Is(err, results.ErrAlreadyRecorded)` becomes a 409); unexpected ones are logged with the request and become a 500.
- **Authorisation is explicit.** Goravel's gate is not told who is asking, so a small `authorise` helper loads the user, passes them to the policy, and turns a denial into a 403.
- **Comparisons live with the code.** The Laravel equivalents are doc comments on the Go they describe, so they change with it and render anywhere the code does, including GitHub. The project page is a single static file in `docs/index.html`.
- **Dependency hygiene.** The pre-commit hook runs gofmt, `go vet` and osv-scanner. Go 1.26 was chosen over the brief's 1.25 because 1.25 is end of life, and three advisories that cannot be fixed without breaking Goravel 1.18 are documented in `osv-scanner.toml`.

## Licence

MIT. Goravel itself is MIT licensed by its authors.
