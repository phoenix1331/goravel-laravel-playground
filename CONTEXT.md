# Context

Teaching repo mapping Laravel concepts to Goravel, one chapter per concept, backed by a working Goravel API for a grassroots football club (clubs, players, fixtures, club-admin users). Audience: Laravel developers learning Go. Clarity beats cleverness. API only: no front end; `GET /` lists endpoints from `facades.Route().GetRoutes()`: HTML with the Goravel logo for browsers (`resources/views/index.tmpl`), JSON otherwise.

## Stack

| Concern | Choice |
|---|---|
| Language | Go 1.26 (Goravel 1.18 needs 1.25+; 1.25 is end of life and its stdlib advisories are only fixed in 1.26.9+) |
| Framework | Goravel 1.18, full scaffold (`goravel/goravel`), module name `goravel` |
| HTTP driver | `goravel/gin` |
| Database (dev) | MySQL 8.4 via `goravel/mysql` |
| Database (tests) | SQLite in memory via `goravel/sqlite` |
| Cache / queue / session | Redis via `goravel/redis`. Queue: `redis` in dev, `sync` in tests |
| Docs | VitePress in `docs/`, deployed to GitHub Pages |

## Services and ports

| Service | Image | Host port | Notes |
|---|---|---|---|
| `app` | `Dockerfile` (golang:1.26 + air) | 8000 | Runs as host UID/GID, project bind-mounted at `/app` |
| `mysql` | mysql:8.4 | 3307 | Host 3306 already in use on the dev machine |
| `redis` | redis:7-alpine | none | Internal only |

## Layout

| Path | Purpose |
|---|---|
| `main.go` | Entry point: `bootstrap.Boot().Start()` |
| `bootstrap/` | App builder (`app.go`), provider list (`providers.go`), migration list (`migrations.go`) |
| `config/` | One Go file per concern, each registers via `facades.Config().Add()` in `init()` |
| `app/facades/` | Thin facade functions, one per service |
| `app/http/` | Controllers, middleware, requests |
| `app/models/` | ORM models |
| `app/services/` | Plain Go services (e.g. `LeagueTable`), bound in the container by `app/providers/` |
| `app/results/` | `results.Record()`: the single path for recording a score (an action), used by `sync:fixtures` and the API |
| `app/console/commands/` | Artisan commands, listed in `bootstrap/commands.go` |
| `app/providers/` | Service providers, listed in `bootstrap/providers.go`. `AuthServiceProvider` defines the gates |
| `app/policies/` | `ClubPolicy`, `FixturePolicy`: club admins manage only their own club; the home admin records results |
| `routes/` | Route registration (`web.go`) |
| `database/` | Migrations, seeders, factories (Go code, not SQL) |
| `tests/` | testify suites |
| `database/factories/`, `database/seeders/` | gofakeit-based factories (keys are struct field names); `DatabaseSeeder` builds a 6-club demo league, admins `admin1..6@example.com` / `password` |
| `docs/` | VitePress site |

## Decisions

| Decision | Why |
|---|---|
| Gin over Fiber | Gin is built on `net/http`, so knowledge transfers to plain Go and the wider ecosystem. Fiber uses fasthttp, which is not `net/http` compatible |
| SQLite driver needs no CGO | `goravel/sqlite` uses `ncruces/go-sqlite3` (SQLite compiled to WebAssembly, run by wazero) |
| Cache, queue and session drivers read from `.env` | `package:install` hard-codes Redis as the default; env keys (`CACHE_STORE`, `QUEUE_CONNECTION`, `SESSION_DRIVER`) let tests switch to `memory`/`sync`/`file` |
| `air` live reload in the `app` container | Go compiles to a binary; air rebuilds and restarts on save. `database/` is watched because migrations are Go code |
| Container runs as host user | Files created in the bind mount (go.sum, tmp, storage) stay owned by the developer |
| Licence MIT | Matches Goravel |
| Queue worker and scheduler run inside the app process | Goravel starts runners (`goravel:http`, `goravel:queue`, `goravel:schedule`) from `app.Start()`; no `queue:work`. `APP_DISABLED_RUNNERS` (space separated) splits them into separate processes |
| Events and jobs carry IDs as typed args | Goravel payloads are `[]event.Arg` / `[]queue.Arg`, not objects; handlers reload from the database |
| League table cached as a JSON string | The Redis store cannot hold Go structs ("can't marshal"); the memory store can. Caching strings behaves the same on both. Never `Cache().Flush()` on shared Redis: Goravel runs `FLUSHALL` |
| JWT auth on writes | Reads and `POST /api/auth/login` are public; writes sit behind `middleware.Authenticate()`. `JWT_SECRET` via `jwt:secret`. Tests log in through the real endpoint with `s.LoginAs(t, user)`; `.env.testing` sets `BCRYPT_ROUNDS=4` and the user factory hashes `password` once. Goravel JWTs have no `jti`, so same-user logins in the same second share a token, and logging one out disables both |
| Authorisation in controllers via `authorise()` | Goravel's gate does not receive the user, so the helper loads it and adds it to the arguments; returns 401/403 JSON. Kept out of form request `Authorize`, which cannot produce a 403 |
| Fixture feed is an interface (`FileFeed`, `HttpFeed`) | `FIXTURE_FEED_DRIVER` picks file or http; the `fixtures` client lives in `config/http.go`; `/api/demo/results` is a local stand-in API. `facades.Http()` returns `client.Factory` (scaffold returned `client.Request`, hiding `Fake()`) |
| Expected failures are sentinel errors | `results.ErrAlreadyRecorded`, wrapped with `%w`, matched with `errors.Is` in the controller (409). Unexpected errors go to `serverError` (logged with the request, 500) |
| Mocks in `tests/mocked` | `mock.Factory()` replaces `foundation.App` process-wide, so mock tests live in a package that never boots the real app |
| Fixture `saving` event converts `kick_off_at` to UTC | SQLite keeps the original offset as text, MySQL converts to UTC; the event makes both return UTC (found by the UTC test failing on SQLite only) |
| Container keys are string constants beside the service (`services.LeagueTableBinding`) | Goravel binds by string, not interface name; a constant prevents typos |
| No `env_file` in `docker-compose.yml` | Goravel's config uses viper `AutomaticEnv`, so OS env vars override any `.env` file. Injecting `.env` as container env would silently override `.env.testing` in tests. The app reads `.env` from the bind mount instead |
| Pool size from `.env` (`DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`) | Lets environments tune the pool without editing `config/database.go` |
| `gofmt` + `go vet` instead of Pint | Built into the Go toolchain, no config file. `make lint` checks, `make fmt` fixes. Go has no `strict_types` equivalent because it is statically typed |

## Chapter conventions

Each chapter delivers working Go code, tests, a docs page and a Laravel comparison.

| Piece | Convention |
|---|---|
| Laravel snippet | Non-runnable PHP file at `docs/snippets/laravel/chNN/<Name>.php` |
| Go source | Real code with `// #region <name>` / `// #endregion <name>` markers around the parts shown in docs |
| Docs page | `docs/chapters/NN-<slug>.md`, VitePress code group: Laravel tab imports the `.php` file, Go tab imports the source region |
| Go file header | Key files start with `// Laravel equivalent: <path> - see docs/chapters/NN-<slug>.md` |
| Inline comments | Only where Go behaves differently from what a Laravel developer expects |
| Evolving code | A region imported by an earlier chapter must not change meaning later. When a later chapter rewrites code (e.g. ch07 validation replacing ch04 `Bind`), the earlier chapter shows the old version as a static code block |
| Snippet index | `docs/laravel-index.md` renders from `docs/laravel-index.data.js`, a VitePress data loader that pairs each snippet with the Go files in the same code group at build time |

## Testing

| Concern | Approach |
|---|---|
| Env | `.env.testing` (tracked), selected with `go test ./... --env=.env.testing`; Goravel walks up from the package dir to find it |
| Database | SQLite in memory, DSN `file:/goravel_test.db?vfs=memdb` |
| Why memdb, not `:memory:` | Each pooled connection to `:memory:` gets its own empty database (proven: second connection sees "no such table"). `cache=shared` is not supported by `ncruces/go-sqlite3`. A pool of 1 works but deadlocks if anything holds the connection. memdb shares one in-memory database across all connections in the process |
| Migrations | Once per package, in `TestMain` via `tests.Main(m)` |
| Reset per test | `s.TruncateTables()` in `SetupTest` deletes rows from every table except `migrations`. Chosen over re-migrating (slower) and per-test transactions (facades use a global connection, so a test transaction cannot wrap app code) |
| Isolation | `go test` runs each package as its own process, so each package gets its own database |
| Production-like run | `make test-mysql` / CI `test-mysql` job rerun the same suite on MySQL and the Redis cache by overriding env vars (OS env beats `.env.testing`). Uses the `goravel_test` database (`docker/mysql/init.sql`), Redis db 1, and `-p 1`, since packages share one database. Catches SQL dialect drift and memory-vs-Redis cache drift |
| Base suite | `tests.TestCase` embeds Goravel's `testing.TestCase` (`Http()`, `Seed()`) and adds `TruncateTables()`, `CreateAdmin()` and `LoginAs()` |

## Docs site

| Item | Detail |
|---|---|
| Tool | VitePress 1.6 at the repo root `package.json`, `vite` overridden to `^6.4.3` (Vite 5 dev-server advisories have no 5.x fix) |
| Config | `docs/.vitepress/config.mts`, base `/goravel-laravel-playground/` for GitHub Pages |
| Pages | `index.md`, `go-for-php-devs.md` (first in nav), `reading-the-comparisons.md` (worked example on `/health`), `where-to-find-it.md` (feature to code map), `chapters/NN-<slug>.md` |
| Go imports | `<<< @/../path/to/file.go#region{go}` pulls a `// #region` block from the real source at build time |
| Runs on | Host Node (`make docs-dev`, `make docs-build`), not the Go container |

## CI

| Workflow | Trigger | Jobs |
|---|---|---|
| `.github/workflows/ci.yml` | push to `main`, PRs | `lint` (gofmt, vet), `test-sqlite`, `test-mysql` (MySQL 8.4 service). Go version from `go.mod` |
| `.github/workflows/docs.yml` | push to `main`, manual | `npm ci`, `docs:build`. Pages deploy steps are commented out while the repo is private |

To publish: make the repo public, set Settings > Pages > Source: **GitHub Actions**, and uncomment the marked lines in `docs.yml`.

## Commands

All developer commands live in the `Makefile` and run inside the `app` container via `docker compose exec app`. Artisan is `go run . artisan <command>` (the `./artisan` script wraps this).

## Dependency audit

`osv-scanner` scans `go.mod` and `package-lock.json` in the pre-commit hook (Go call analysis off, as the host Go may not match `go.mod`). Accepted advisories live in `osv-scanner.toml`, each with a reason. Two are OpenTelemetry log packages whose fixed versions break `goravel/framework` v1.18.0, so they must be revisited on every framework upgrade.
