---
title: "12. Factories and seeders"
---

# 12. Factories and seeders

| Laravel | Goravel |
|---|---|
| `Club::factory()->count(3)->create()`, `fake()` | `facades.Orm().Factory().Count(3).Create(&clubs)`, gofakeit |

Factories work the same way in both: a definition returns default attributes, and callers override what they care about. Goravel's factories are smaller than Laravel's. There are no states, sequences or `for()`/`has()` relationship helpers, so relationships are set by passing foreign keys.

## The Laravel way

<<< @/snippets/laravel/ch12/ClubFactory.php

## The Goravel way

`make artisan cmd="make:factory ClubFactory"` creates the struct. `Definition` returns a map keyed by **struct field name**:

<<< @/../database/factories/club_factory.go#club{go}

[gofakeit](https://github.com/brianvoe/gofakeit) stands in for Faker. Its place names are American, so the factory picks from a list of English towns for a grassroots feel.

The model points at its factory with a `Factory()` method, the equivalent of the `HasFactory` trait:

<<< @/../app/models/club.go#factory{go}

Using it:

```go
var clubs []models.Club
facades.Orm().Factory().Count(3).Create(&clubs)            // three saved clubs

var club models.Club
facades.Orm().Factory().Make(&club)                         // built, not saved

var player models.Player
facades.Orm().Factory().Create(&player, map[string]any{     // overrides
	"ClubID":      club.ID,
	"SquadNumber": 7,
})
```

### Seeders

`make artisan cmd="make:seeder DatabaseSeeder"` creates the seeder and registers it in `bootstrap/seeders.go`, which `bootstrap/app.go` loads with `WithSeeders(Seeders)`.

::: code-group

<<< @/snippets/laravel/ch12/DatabaseSeeder.php [Laravel]

<<< @/../database/seeders/database_seeder.go#run{go} [Goravel]

:::

`make fresh` runs `migrate:fresh --seed`, giving a league of six clubs, 96 players, six admins and fifteen fixtures, seven of them already played. Every admin's password is `password`.

## Testing it

::: code-group

<<< @/snippets/laravel/ch12/FactoriesTest.php [Laravel]

<<< @/../tests/feature/factories_test.go#test{go} [Goravel]

:::

The earlier chapters' test helpers now use factories too. `createClub` in the chapter 04 tests went from building a struct by hand to one factory call with a name override.

## What is different and why

- **Keys are field names, not column names.** `"SquadNumber"`, not `"squad_number"`. The factory decodes the map into the struct, so a misspelt key is silently ignored rather than failing to compile. Check overrides in a test.
- **Pass values in the field's type, or something convertible.** `"Position": "forward"` works because the decoder converts a string to `models.Position`. A `*int` field such as `FoundedYear` needs a pointer: `"FoundedYear": &year`.
- **No relationship helpers.** Laravel's `for($club)` and `has(Player::factory()->count(16))` have no equivalent. Create the parent, then pass its ID as an override. It is more typing, but the order of inserts is obvious.
- **No states or sequences.** Laravel's `->state(['status' => 'finished'])` and `->sequence(...)` are replaced by plain Go: pass different overrides, or loop. The seeder's `squad` function is a loop where Laravel would use a sequence.
- **Hash once, not per user.** `UserFactory` stores a precomputed bcrypt hash of `password` instead of hashing on every call, the same trick Laravel's own `UserFactory` uses with its `static $password`.
- **Seeding gofakeit is global.** `gofakeit.Seed(2026)` makes the demo data repeatable, but it reseeds the shared generator for the whole process. Tests that run after the seeder in the same package get predictable values too, so do not rely on fake data being different between tests.
