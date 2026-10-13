---
title: "10. Relationships"
---

# 10. Relationships

| Laravel | Goravel |
|---|---|
| Relationship methods (`hasMany`, `belongsTo`), `with()`, `load()` | Struct fields for relations, `With()`, `Load()` |

In Eloquent a relationship is a method that returns a query object, and `$club->players` lazily runs that query the first time you touch it. In Goravel a relationship is a struct field. It stays empty until you load it, and nothing loads by accident.

## The Laravel way

<<< @/snippets/laravel/ch10/models.php

## The Goravel way

### Has many and belongs to

The relation is a field on the struct. `Club` gets a slice of players:

<<< @/../app/models/club.go#has-many{go}

`Player` gets a pointer to its club, next to the `ClubID` foreign key it already had:

```go
type Player struct {
	orm.Model
	ClubID      uint     `json:"club_id"`
	// ...
	Club *Club `json:"club,omitempty"`
}
```

The ORM pairs `Club` with `ClubID` by name, the same convention Laravel uses for `club_id`.

### Two foreign keys to the same table

A fixture has a home club and an away club, both in `clubs`. The field names do not match the convention, so a `gorm` tag names each key:

<<< @/../app/models/fixture.go#two-fks{go}

### Eager loading while querying: `With`

<<< @/../app/http/controllers/fixture_controller.go#index{go}

`With("HomeClub")` takes the **struct field name**, capitalised, not a method name. `With("homeClub")` fails at runtime. Nested relations use dots, as in Laravel: `With("HomeClub.Players")`.

### Loading onto a model you already have: `Load`

<<< @/../app/http/controllers/player_controller.go#index{go}

### Creating a related record

<<< @/../app/http/controllers/player_controller.go#store{go}

There is no `$club->players()->create(...)`. You set `ClubID` on the new struct and create it. For many-to-many and other cases, `facades.Orm().Query().Model(&club).Association("Players")` offers `Append`, `Replace`, `Delete`, `Clear` and `Count`.

The routes nest under the club, using the `{id}` that `findClub` from chapter 04 already reads:

```go
playerController := controllers.NewPlayerController()
router.Get("clubs/{id}/players", playerController.Index)
router.Post("clubs/{id}/players", playerController.Store)

router.Get("fixtures", controllers.NewFixtureController().Index)
```

[Chapter 19](/chapters/19-authentication) moves the `POST` behind authentication.

### Unique within a parent

::: code-group

<<< @/snippets/laravel/ch10/StorePlayerRequest.php [Laravel]

```go [Goravel]
taken, err := facades.Orm().Query().Model(&models.Player{}).
	Where("club_id", club.ID).Where("squad_number", request.SquadNumber).Exists()
```

:::

Goravel's `unique` rule has no `where` clause, so the "squad number is free at this club" check runs in the controller and returns a 422 in the same shape. The database also enforces it with a unique index on `(club_id, squad_number)`, covered in chapter 11.

## Testing it

::: code-group

<<< @/snippets/laravel/ch10/RelationshipsTest.php [Laravel]

<<< @/../tests/feature/relationships_test.go#test{go} [Goravel]

:::

## What is different and why

- **No lazy loading, so no accidental N+1.** `$fixture->homeClub` in a Blade loop quietly runs one query per row. In Goravel, `fixture.HomeClub` is just a field: `nil` if you did not load it. Reading a field on a `nil` pointer panics (`nil pointer dereference`), so a missing `With` shows up the first time the code runs instead of as a slow page later. Laravel's answer is `Model::preventLazyLoading()`; Goravel behaves that way all the time.
- **Relations are fields, so they need `omitempty`.** Without `json:"players,omitempty"`, every club in every response would carry `"players": null`. With it, the key only appears when the relation was loaded.
- **Relation names are Go identifiers.** `With("HomeClub")` matches the field, and the string is not checked at compile time. A typo is a runtime error, so keep a test that exercises each `With`.
- **Two queries, not a join.** Like Laravel's `with()`, `With` runs one query for fixtures and one `WHERE id IN (...)` query per relation, then stitches the results together. Use a join (chapter 13) when you need to filter or sort by a related column.
- **Saving a model does not save its relations.** Goravel omits associations on `Create` and `Save` by default, which matches Eloquent's `save()`. Attaching a new `Club` struct to a player and creating the player does not insert the club; with SQLite's foreign keys on, it fails because `club_id` is 0. To save related structs in one go, opt in with `facades.Orm().Query().Select(orm.Associations).Create(&player)`, using `orm` from `github.com/goravel/framework/database/orm`. Plain GORM does the opposite and saves associations by default, so older GORM tutorials will mislead you here.
