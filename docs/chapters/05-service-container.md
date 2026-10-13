---
title: "05. Service container and providers"
---

# 05. Service container and providers

| Laravel | Goravel |
|---|---|
| `$this->app->singleton()`, `register()` / `boot()`, autowiring | `app.Singleton()`, `Register` / `Boot`, explicit `Make` |

Goravel has a container and service providers with the same two-phase life cycle as Laravel. The difference is how you get things out. Laravel's container reads constructor type hints and builds dependencies for you (autowiring). Goravel's container is a map from a key to a factory function, and you ask for things by key.

The example service is the league table: given a list of results, rank the clubs. Chapter 13 feeds it real fixtures from the database.

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch05/LeagueTable.php [LeagueTable.php]

<<< @/snippets/laravel/ch05/LeagueServiceProvider.php [LeagueServiceProvider.php]

<<< @/snippets/laravel/ch05/resolve.php [Resolving]

:::

## The Goravel way

### The interface and implementation

<<< @/../app/services/league_table.go#interface{go}

<<< @/../app/services/league_table.go#implementation{go}

`PointsLeagueTable` never says it implements `LeagueTable`. It has a `Standings` method with the right signature, so it satisfies the interface. The compiler checks this wherever a `PointsLeagueTable` is used as a `LeagueTable`.

### The provider

`make artisan cmd="make:provider LeagueServiceProvider"` creates the file and adds it to `bootstrap/providers.go` for you, the same as Laravel 11+ does with `bootstrap/providers.php`.

<<< @/../app/providers/league_service_provider.go#register{go}

The binding key is a string constant, `services.LeagueTableBinding` (`"club.league_table"`), not the interface's name. Go cannot use a type as a map key the way PHP uses `LeagueTable::class`.

`Relationship()` is Goravel-specific. It declares that this provider binds `club.league_table` and depends on config, so Goravel registers providers in dependency order. Laravel relies on the order of `bootstrap/providers.php` instead.

### Resolving

<<< @/../tests/unit/league_table_test.go#container{go}

`Make` returns `any`, Go's equivalent of `mixed`. The type assertion `first.(services.LeagueTable)` turns it back into something you can call methods on. Chapter 06 wraps this in a facade function so callers never see the `any`.

## Testing it

::: code-group

<<< @/snippets/laravel/ch05/LeagueTableTest.php [Laravel]

<<< @/../tests/unit/league_table_test.go#table{go} [Goravel]

:::

The table-driven test is Go's data provider. Each case is a struct literal in a slice, and `t.Run` turns each into a named subtest: `go test -run 'TestPointsLeagueTableRanking/goal_difference'` runs just one.

This test does not touch the container at all. `NewPointsLeagueTable(3, 1, 0)` is called directly, which is the usual Go style: construct the thing under test with plain function calls.

## What is different and why

- **No autowiring.** Laravel inspects `__construct(LeagueTable $table)` at runtime and resolves the argument. Go can inspect types at runtime (reflection), but not parameter type hints in a way that makes autowiring practical, and the Go community prefers explicit wiring anyway. Dependencies are passed into constructor functions by hand.
- **String keys, not class names.** Bindings are keyed by a string like `"club.league_table"`. Keep the key in a constant next to the service so callers cannot misspell it.
- **`any` and type assertions.** Whatever comes out of the container is `any`. Use the two-value form `v, ok := x.(T)` to check the type safely. The one-value form panics if the type is wrong.
- **Interfaces are satisfied implicitly.** No `implements`. Interfaces are usually defined where they are used and kept small (one or two methods is common), which makes them cheap to fake in tests.
- **You often do not need the container.** In Laravel the container is how everything is built. In Go, calling `NewPointsLeagueTable(...)` directly is normal and testable. Use the container for things that are configured once and shared, such as framework services, or when you want facades (chapter 06) to hand out a swappable implementation.
