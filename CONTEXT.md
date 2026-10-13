# Context

Teaching repo mapping Laravel concepts to Goravel, one chapter per concept, backed by a working Goravel API for a grassroots football club (clubs, players, fixtures, club-admin users). Audience: Laravel developers learning Go. Clarity beats cleverness.

## Stack

| Concern | Choice |
|---|---|
| Language | Go 1.26 (Goravel 1.18 needs 1.25+; 1.25 is end of life and its stdlib advisories are only fixed in 1.26.9+) |
| Framework | Goravel 1.18, full scaffold (`goravel/goravel`), module name `goravel` |
| HTTP driver | `goravel/gin` |
| Database (dev) | MySQL 8.4 via `goravel/mysql` |
| Database (tests) | SQLite in memory via `goravel/sqlite` |
| Cache / queue / session | Redis via `goravel/redis`; queue defaults to `sync` |
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
| `routes/` | Route registration (`web.go`) |
| `database/` | Migrations, seeders, factories (Go code, not SQL) |
| `tests/` | testify suites |
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
| Snippet index | `docs/laravel-index.md` lists every PHP snippet against its Go counterpart |

## Commands

All developer commands live in the `Makefile` and run inside the `app` container via `docker compose exec app`. Artisan is `go run . artisan <command>` (the `./artisan` script wraps this).

## Dependency audit

`osv-scanner` scans `go.mod` and `package-lock.json` in the pre-commit hook (Go call analysis off, as the host Go may not match `go.mod`). Accepted advisories live in `osv-scanner.toml`, each with a reason. Two are OpenTelemetry log packages whose fixed versions break `goravel/framework` v1.18.0, so they must be revisited on every framework upgrade.
