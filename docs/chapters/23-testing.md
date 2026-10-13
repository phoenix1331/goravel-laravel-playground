---
title: "23. Testing"
---

# 23. Testing

| Laravel | Goravel |
|---|---|
| PHPUnit or Pest, `RefreshDatabase`, `actingAs`, fakes | `go test` with testify suites, a hand-built harness, `LoginAs`, mocks |

Every chapter so far has come with tests, so this one is a tour of the harness behind them and the decisions it involved. Goravel's [testing docs](https://www.goravel.dev/testing/getting-started.html) cover the framework's own helpers, including Docker-based databases this project chose not to use.

```bash
make test         # SQLite in memory, memory cache, sync queue
make test-mysql   # the same suite on MySQL and Redis
```

## The shape of a test

::: code-group

<<< @/snippets/laravel/ch23/TestCase.php [Laravel]

<<< @/../tests/feature/club_controller_test.go#test{go} [Goravel]

:::

| PHPUnit / Laravel | testify / Goravel |
|---|---|
| A class extending `TestCase` | A struct embedding `suite.Suite` and `tests.TestCase` |
| `test_*` methods | `Test*` methods |
| `setUp()` | `SetupTest()` (and `SetupSuite()` for once per suite) |
| `$this->assertSame()` | `s.Equal()`; `s.Require().Equal()` stops the test on failure |
| Data providers | Table-driven tests: a slice of cases and `t.Run` (chapter 05) |
| `phpunit --filter name` | `go test -run 'Suite/TestName' ./tests/feature/` |

The one-line `func TestClubControllerTestSuite(t *testing.T) { suite.Run(...) }` in every file is what hands the suite to `go test`, which only knows about plain `Test*` functions.

## The harness

### Environment

::: code-group

<<< @/snippets/laravel/ch23/phpunit.xml [Laravel: phpunit.xml]

<<< @/../.env.testing{ini} [Goravel: .env.testing]

:::

`make test` passes `--env=.env.testing`. Real environment variables override the file, which is how `make test-mysql` and CI switch the same suite to MySQL and Redis without a second env file.

### In-memory SQLite

<<< @/../tests/test_case.go#main{go}

Three decisions, each tried and tested before settling:

| Question | Answer |
|---|---|
| `:memory:`, shared cache, or something else? | `file:/goravel_test.db?vfs=memdb`. With `:memory:` each pooled connection gets its own empty database, so a query on a second connection fails with "no such table". `cache=shared` is not supported by the driver. A pool of one connection works but deadlocks if anything holds it |
| Does the driver need CGO? | No. `goravel/sqlite` uses `ncruces/go-sqlite3`, which runs SQLite as WebAssembly |
| How is each test reset? | Migrate once per package in `TestMain`, then `TruncateTables()` in `SetupTest` |

<<< @/../tests/test_case.go#truncate{go}

Laravel's `RefreshDatabase` wraps each test in a transaction. That does not work here: the app reaches the database through global facades, so a transaction opened by the test cannot wrap the code under test. Deleting rows is fast on a small in-memory database.

`go test` runs each package as a separate process, so `tests/feature`, `tests/unit` and `tests/mocked` each get their own in-memory database and never interfere.

### Logging in

<<< @/../tests/test_case.go#login{go}

There is no `actingAs`, because the JWT guard needs a real request context. `LoginAs` logs in through the endpoint, made cheap by `BCRYPT_ROUNDS=4` (chapter 19).

### Mocks need their own package

Goravel's generated mocks replace the whole application object, process-wide:

::: code-group

<<< @/snippets/laravel/ch23/mock.php [Laravel]

<<< @/../tests/mocked/forget_league_table_test.go#test{go} [Goravel]

:::

Put that in `tests/feature` and every other test in the package would talk to a mock app. So mock-based tests live in `tests/mocked`, which never boots the real app. Most tests in this repo check outcomes in the database instead, which survives refactoring better.

## SQLite is not MySQL

Running fast tests on SQLite is a trade-off. The `test-mysql` job exists because SQLite and the memory cache accept things production would not. These are the differences this project ran into or designed around:

| Found in | SQLite / memory | MySQL / Redis | Fix |
|---|---|---|---|
| Chapter 11 | Keeps a time's original offset (`15:00+01:00`) | Converts to UTC | Model `saving` event converts kick-off to UTC |
| Chapter 13 | `SUM()` returns an integer | `SUM()` returns a decimal | Scan into `int`; both convert |
| Chapter 18 | The memory cache stores Go structs | Redis "can't marshal" a struct | Cache JSON strings |
| Chapter 18 | `Flush()` clears the memory cache | `Flush()` runs `FLUSHALL` on the whole Redis server | Tests forget specific keys |

Both runs use the same test code. If the suite only ran on SQLite and the memory cache, the timezone and struct-caching bugs would have shipped.

## What is different and why

- **No test runner to configure.** `go test` is built into the toolchain. There is no `phpunit.xml`, and the harness is ordinary Go code in `tests/test_case.go`.
- **Packages are processes.** Isolation comes from package boundaries, not from resetting a container between tests. Inside a package, state that outlives a test (config, fakes, cache, open log files) must be put back by the test that changed it.
- **Fakes vary by facade.** The HTTP client has Laravel-style fakes (chapter 21). Cache, queue and events do not; this project tests their effects or uses the generated mocks.
- **`go test` caches results.** It skips packages whose code has not changed. `make test` passes `-count=1` because these tests depend on things the cache cannot see, such as the database.
