---
title: "18. Cache"
---

# 18. Cache

| Laravel | Goravel |
|---|---|
| `Cache::remember()`, `Cache::forget()` | `facades.Cache().Remember()`, `facades.Cache().Forget()` |

The league table from chapter 13 is rebuilt from every finished fixture on each request. It only changes when a result is recorded, so it is a good thing to cache: remember it until a result comes in, then forget it. The API is close to Laravel's. The trap is what you are allowed to put in the cache.

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch18/LeagueController.php [Caching]

<<< @/snippets/laravel/ch18/ForgetLeagueTable.php [Invalidating]

:::

## The Goravel way

<<< @/../app/http/controllers/league_controller.go#table{go}

The table itself is built by `buildLeagueTable`, the chapter 13 query code moved into its own function.

A second listener on chapter 15's event clears the cache:

<<< @/../app/listeners/forget_league_table.go#listener{go}

```go
&events.FixtureResultRecorded{}: {
	&listeners.NotifyClubAdmins{},
	&listeners.ForgetLeagueTable{},
},
```

### Cache only strings

The obvious version caches the rows directly:

```go
facades.Cache().Remember(key, ttl, func() (any, error) {
	return buildLeagueTable() // []leagueRow
})
```

With the memory store that the tests use, it works. With Redis it fails:

```
redis: can't marshal []services.Standing (implement encoding.BinaryMarshaler)
```

The memory store keeps whatever Go value you give it. The Redis store hands the value to the Redis client, which only knows how to store strings and numbers. Reading back is uneven too: Redis returns a `string` even if you stored a number. Laravel hides all of this by running every value through PHP's `serialize()`. Goravel does not.

So the controller encodes the rows to JSON itself, caches the string, and decodes it after `Remember`. That behaves the same on every store.

### Catching it in CI

The test suite runs on the memory store, so this bug would never show up there. The `test-mysql` CI job and `make test-mysql` now also switch the cache to Redis (database 1, so tests never touch development data). The production-like run is where store differences show up, just as it is for SQL dialects.

## Testing it

::: code-group

<<< @/snippets/laravel/ch18/LeagueTableCacheTest.php [Laravel]

<<< @/../tests/feature/cache_test.go#test{go} [Goravel]

:::

The middle test proves the cache is really used: it changes the database without going through `results.Record`, and the table does not change.

The tests reset with `Forget(key)`, not `Flush()`. On Redis, Goravel's `Flush()` runs `FLUSHALL`, which empties **every database on the Redis server**, including the development queue and sessions. Laravel's `flush()` runs `FLUSHDB`, which empties only the connection's database. Treat `Flush()` as dangerous anywhere Redis is shared.

## What is different and why

- **No automatic serialisation.** Laravel `serialize()`s anything into the cache. Goravel passes values straight to the store, so what works depends on the store. Encode structs yourself (JSON is simplest) and cache strings.
- **`Remember` returns `any`.** Convert it back with `cast.ToString` and decode. The type you get can differ by store: a number stored in Redis comes back as a string.
- **`Flush()` is broader than you think on Redis.** `FLUSHALL`, not `FLUSHDB`. Forget specific keys instead, or give tests their own Redis server.
- **Invalidation through events.** Forgetting the key in a listener means the code that records results does not need to know a cache exists. A Laravel cache tag or a model observer would do the same job.
- **Keys are prefixed.** `config/cache.go` prefixes keys with `APP_NAME`, so development (`ClubManager_cache:league_table`) and tests (`ClubManagerTest_cache:league_table`) never share a key.
