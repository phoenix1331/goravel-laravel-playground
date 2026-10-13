---
title: "13. Query builder and raw SQL"
---

# 13. Query builder and raw SQL

| Laravel | Goravel |
|---|---|
| `DB::table()`, `DB::select()` | `facades.DB().Table()`, `facades.DB().Select()` |

Goravel has two ways to query without models. `facades.Orm().Query()` is the ORM builder from chapters 09 and 10. `facades.DB()` is a lighter query builder that works on table names and scans rows into any struct, the equivalent of the `DB` facade. This chapter uses `facades.DB()` to build the league table that chapter 05's service has been waiting for.

## The Laravel way

<<< @/snippets/laravel/ch13/LeagueController.php

## The Goravel way

### The query builder

<<< @/../app/http/controllers/league_controller.go#build{go}

The rows scan straight into `[]services.Result`, the struct the league table service already takes. There are no tags on `Result`: `facades.DB()` converts each field name to snake case (`HomeClubID` to `home_club_id`) and matches it to a column. Laravel gives you `stdClass` rows that you map by hand. Go gives you typed structs.

### Raw SQL

Some queries are clearer in SQL. Goals per club have to count home goals and away goals, which live in different columns:

<<< @/../app/http/controllers/league_controller.go#goals{go}

### Placeholders across drivers

The query builder writes the placeholder style for the current driver:

```go
facades.DB().Table("fixtures").Where("status", "finished").ToSql().Get()
// SELECT * FROM fixtures WHERE status = ?

facades.DB().Table("fixtures").Where("status", "finished").ToRawSql().Get()
// SELECT * FROM fixtures WHERE status = 'finished'
```

On PostgreSQL the same builder call produces `status = $1`. **Raw SQL is not rewritten.** `facades.DB().Select(dest, sql, args...)` hands your string to the driver unchanged, so raw SQL needs `?` on MySQL and SQLite, and `$1, $2` on PostgreSQL. PHP's PDO accepts `?` on every driver, so Laravel developers never meet this. If an app might change databases, prefer the builder, and keep raw SQL in one place.

### Keeping SQL portable

The suite runs on SQLite and MySQL (chapter 23), so `goalsSQL` sticks to SQL both accept:

- `UNION ALL` in a derived table, aliased with `AS scored`. MySQL requires the alias.
- No `IFNULL` versus `COALESCE` choice to make: only finished fixtures are counted, so scores are never NULL.
- No date functions: these differ the most between engines (`DATE_FORMAT` against `strftime`).
- `SUM()` returns `DECIMAL` on MySQL and `INTEGER` on SQLite. Both scan into a Go `int`, so the difference never reaches the API.

## Testing it

::: code-group

<<< @/snippets/laravel/ch13/LeagueTest.php [Laravel]

<<< @/../tests/feature/league_test.go#test{go} [Goravel]

:::

The test's `SetupTest` builds a four-club league with factories: two results, one fixture still to play, and one club with no fixtures at all. That covers ranking, ignoring unplayed fixtures, and the club that has not played yet.

## What is different and why

- **Rows scan into structs.** `Get(&results)` fills a typed slice. A column the struct does not have is ignored, and a field with no column stays at its zero value. There is no `stdClass` and no array of arrays.
- **Raw SQL placeholders are driver-specific.** The builder handles them; raw strings do not. This is the most likely thing to break if the app ever moves to PostgreSQL.
- **`facades.DB()` against `facades.Orm().Query()`.** Use the ORM for models, relationships and events (the UTC `saving` event from chapter 11 only fires through the ORM). Use `facades.DB()` for reports, aggregates and anything that is not a model.
- **`ToSql()` and `ToRawSql()`** are `toSql()` and `toRawSql()`. They are useful in tests to pin down the SQL a builder call produces.
- **Errors are returned, every time.** `DB::select()` throws a `QueryException`. `facades.DB().Select()` returns an error, and the controller turns it into a logged 500 with `serverError`.
