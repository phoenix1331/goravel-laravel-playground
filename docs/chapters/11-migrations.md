---
title: "11. Migrations"
---

# 11. Migrations

| Laravel | Goravel |
|---|---|
| Anonymous migration classes, `Schema::create()`, `Blueprint` | Migration structs, `facades.Schema().Create()`, `schema.Blueprint` |

Goravel's schema builder is a near copy of Laravel's: `ID()`, `String()`, `Foreign().References().On()`, `Nullable()`, `Index()`. The differences are in how migrations are found and how they are named.

The schema now has four tables: `clubs`, `players`, `fixtures` and `users`.

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch11/create_fixtures_table.php [Create]

<<< @/snippets/laravel/ch11/add_short_name_to_clubs_table.php [Alter]

:::

## The Goravel way

`make artisan cmd="make:migration create_fixtures_table"` creates a struct with `Signature`, `Up` and `Down` methods, and adds it to `bootstrap/migrations.go`:

<<< @/../database/migrations/20261013183455_create_fixtures_table.go#create{go}

Altering a table uses `Table()`, like `Schema::table()`:

<<< @/../database/migrations/20261013184001_add_short_name_to_clubs_table.go#alter{go}

### Registration

Laravel finds migrations by scanning `database/migrations`. Goravel cannot scan for Go code at runtime, because it is compiled in, so every migration is listed in `bootstrap/migrations.go`:

```go
func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20261013182221CreateClubsTable{},
		&migrations.M20261013183449CreatePlayersTable{},
		// ...
	}
}
```

`make:migration` adds the line for you. A migration file that is not in this list never runs.

### Commands

| Laravel | Goravel (`make artisan cmd="..."`) |
|---|---|
| `php artisan migrate` | `migrate` (or `make migrate`) |
| `php artisan migrate:rollback` | `migrate:rollback` |
| `php artisan migrate:fresh --seed` | `migrate:fresh --seed` (or `make fresh`) |
| `php artisan migrate:status` | `migrate:status` |

## Storing kick-off in UTC

`config/app.go` sets the app timezone to UTC, and `kick_off_at` is a `TimestampTz` column. That is not enough on its own, and the test suite proved it:

- **MySQL** converts the time to UTC on the way in. A 3pm BST kick-off reads back as `14:00 UTC`.
- **SQLite** has no real timestamp type. It stores the value as text, offset included, and reads back `15:00 +01:00`.

Same instant, different representation, and the test that checks for 14:00 UTC failed on SQLite only. Laravel would fix this with a mutator. Goravel's equivalent is a model event:

::: code-group

<<< @/snippets/laravel/ch11/Fixture.php [Laravel]

<<< @/../app/models/fixture.go#events{go} [Goravel]

:::

This is the dialect drift the brief warned about, and the reason CI also runs the suite on MySQL.

## Testing it

::: code-group

<<< @/snippets/laravel/ch11/SchemaTest.php [Laravel]

<<< @/../tests/feature/migrations_test.go#test{go} [Goravel]

:::

## What is different and why

- **Migrations are structs, named by timestamp.** `M20261013183455CreateFixturesTable` is an ugly name, but Go type names must be unique within a package, and every migration lives in one package. The timestamp keeps them unique and in order. `Signature()` is the name recorded in the `migrations` table.
- **Explicit registration.** There is no directory scan, so a migration missing from `bootstrap/migrations.go` is silently skipped.
- **`Up` and `Down` return errors.** Laravel throws on failure. Goravel returns the error from `Create` or `Table`, and the migrator stops.
- **No `foreignId()->constrained()`.** This repo declares the column and the constraint separately: `UnsignedBigInteger("club_id")` plus `Foreign("club_id").References("id").On("clubs")`. Goravel has a `ForeignID` helper too, but the long form reads the same as Laravel's long form and works across every driver.
- **SQLite enforces foreign keys here.** Laravel turns SQLite foreign keys on by default; so does this project's driver. Creating a player with `club_id` 0 fails with `FOREIGN KEY constraint failed` on both databases, which keeps tests honest.
- **Down migrations matter more with SQLite.** Altering columns on SQLite is limited. Dropping a column works on SQLite 3.35+, which the bundled driver includes, but changing a column type usually means rebuilding the table. Keep test schemas simple.
