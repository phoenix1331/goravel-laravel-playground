---
title: "09. ORM models"
---

# 09. ORM models

| Laravel | Goravel |
|---|---|
| Eloquent models extending `Model`, attributes via `__get` | Structs embedding `orm.Model`, attributes as fields |

An Eloquent model is a class with no declared properties: columns arrive at runtime through `__get` and `__set`, and `$casts` says how to convert them. A Goravel model is a struct that declares every column as a typed field. The ORM (built on [GORM](https://gorm.io)) maps fields to columns by name, so `SquadNumber` reads and writes `squad_number`.

The domain now has four models: `Club` (from chapter 04), `Player`, `Fixture` and `User`. Their tables come from migrations covered in [chapter 11](/chapters/11-migrations).

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch09/Player.php [Player.php]

<<< @/snippets/laravel/ch09/Fixture.php [Fixture.php]

<<< @/snippets/laravel/ch09/User.php [User.php]

:::

## The Goravel way

`make artisan cmd="make:model Player"` creates a struct embedding `orm.Model`. You add the fields:

<<< @/../app/models/player.go#player{go}

`orm.Model` is embedded, not inherited. Embedding copies its fields (`ID`, `CreatedAt`, `UpdatedAt`) into `Player` as if they were declared there, so `player.ID` works. There is no parent class and no `parent::` call.

### Enums

<<< @/../app/models/player.go#position{go}

`Position` is a new type whose underlying type is `string`. The ORM stores it as the string, and the compiler stops you passing a plain `string` where a `Position` is expected without converting it. It is weaker than a PHP backed enum: `Position("striker")` compiles, so validate input against `Positions` (chapter 10 does).

### NULL vs zero values

<<< @/../app/models/fixture.go#fixture{go}

Every Go type has a zero value: `0` for `int`, `""` for `string`. A plain `int` cannot be NULL, so a NULL score would read back as `0`, which is a real score. A pointer (`*int`) can be `nil`, which maps to NULL. Use a pointer for every nullable column where the zero value means something.

### Accessors are methods

<<< @/../app/models/fixture.go#methods{go}

### Hiding fields

<<< @/../app/models/user.go#user{go}

`json:"-"` is the equivalent of `$hidden`. It applies to every JSON encoding of the struct, wherever it happens.

## Testing it

::: code-group

<<< @/snippets/laravel/ch09/ModelsTest.php [Laravel]

<<< @/../tests/feature/models_test.go#test{go} [Goravel]

:::

## What is different and why

- **No `__get`.** Every column is a declared field with a type. Misspell `player.SquadNumbr` and the build fails. In Eloquent, `$player->squad_numbr` is `null` at runtime.
- **No mass assignment, so no `$fillable`.** You build a struct literal, `models.Player{Name: ..., Position: ...}`. Request input is never passed straight in: it goes through a request struct first (chapter 07), and you copy the fields you want.
- **Types replace `$casts`.** A `time.Time` field is a datetime, a `*int` is a nullable integer, a `Position` is an enum. The conversion happens in the database driver, not in a cast list.
- **Pass pointers to the ORM.** `Create(&club)` takes the address of `club` so the ORM can write the new `ID` and timestamps back into it. Pass `club` without `&` and the ORM gets a copy, and your variable never sees the ID.
- **Table names are inferred, with an escape hatch.** `Player` maps to `players`, as in Laravel. Override it by adding a method `func (p *Player) TableName() string { return "squad_members" }`, the equivalent of `protected $table`.
- **Models are plain data.** A Goravel model has no `save()` or `delete()` methods. You hand it to the query builder: `facades.Orm().Query().Save(&club)`. That keeps models simple to construct in tests, at the cost of Eloquent's `$club->save()` style.
