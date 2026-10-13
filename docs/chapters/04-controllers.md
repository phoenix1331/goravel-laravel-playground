---
title: "04. Controllers"
---

# 04. Controllers

| Laravel | Goravel |
|---|---|
| Controller classes, `Route::apiResource()` | Controller structs, `router.Resource()` |

A Goravel controller is a struct whose methods are route handlers. `make:controller ClubController --resource` scaffolds the same five methods as Laravel's `--api` controller. What you lose is the container's magic: no constructor injection and no route model binding. What you gain is that every step is visible in the code.

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch04/ClubController.php [ClubController.php]

<<< @/snippets/laravel/ch04/routes.php [routes/api.php]

:::

## The Goravel way

The controller is a struct with a constructor function:

<<< @/../app/http/controllers/club_controller.go#struct{go}

It is registered once in `routes/api.go`. `Resource` maps the same five routes as `apiResource`, with the parameter named `{id}`:

```go
// One line registers index, store, show, update and destroy,
// like Route::apiResource('clubs', ClubController::class)
router.Resource("clubs", controllers.NewClubController())
```

[Chapter 19](/chapters/19-authentication) later splits these five routes so the reads stay public and the writes need a login, which `Resource` cannot do on its own.

A handler method looks like the closures from chapter 03, with a receiver `(r *ClubController)` in front:

<<< @/../app/http/controllers/club_controller.go#index{go}

### Standing in for route model binding

Laravel's `show(Club $club)` does three jobs: reads `{club}` from the URL, loads the row, and turns "not found" into a 404. In Goravel you write that once as a helper and call it from each method:

<<< @/../app/http/controllers/club_controller.go#find{go}

<<< @/../app/http/controllers/club_controller.go#show{go}

The `(*models.Club, http.Response)` return is a common Go shape: return the value or a reason to stop, and let the caller check which. Chapter 24 has more on returning errors instead of throwing them.

### Reading the request body

The first version of `Store` bound the JSON body straight into a struct:

```go
// clubInput is the JSON body accepted by Store and Update.
type clubInput struct {
	Name        string `json:"name"`
	Ground      string `json:"ground"`
	FoundedYear *int   `json:"founded_year"`
}

func (r *ClubController) Store(ctx http.Context) http.Response {
	var input clubInput
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid JSON body"})
	}

	club := models.Club{Name: input.Name, Ground: input.Ground, FoundedYear: input.FoundedYear}
	if err := facades.Orm().Query().Create(&club); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{"data": club})
}
```

`$request->only(...)` becomes a struct describing the expected JSON. `Bind` decodes the body into it. Fields not in the struct are ignored, which gives you the same protection as `$fillable` without a separate list. There was no validation, so an empty name was accepted. [Chapter 07](/chapters/07-validation) replaces `clubInput` with a form request, and the current `Store` lives there.

## Testing it

::: code-group

<<< @/snippets/laravel/ch04/ClubControllerTest.php [Laravel]

<<< @/../tests/feature/club_controller_test.go#test{go} [Goravel]

:::

`response.Bind(&body)` decodes the JSON response into a typed struct, so the test compares `string` and `int` fields rather than `float64` values from a `map[string]any` (the chapter 03 gotcha).

## What is different and why

- **No constructor injection.** Laravel's container reads the constructor's type hints and builds the dependencies for you. Go has no runtime type-hint reflection to drive that, so `NewClubController()` is an ordinary function you call yourself. To give the controller a dependency, add a struct field and pass it into the constructor. Chapter 05 shows the container side.
- **No route model binding.** The URL gives you a string; `RouteInt("id")` converts it, and you load the model. A small helper like `findClub` keeps that to one line per method.
- **Errors are returned, not thrown.** `FindOrFail` does not throw a `ModelNotFoundException` for the framework to catch. It returns an error, and `errors.Is(err, frameworkerrors.OrmRecordNotFound)` checks which kind it is. Every database call in the controller is followed by `if err != nil`.
- **Hide internal errors.** `serverError` logs the real error and sends a generic 500. Laravel's exception handler does this for you when `APP_DEBUG=false`. In Go, if you return `err.Error()` to the client, database details leak.
- **The struct shape is the API.** `Bind` and the JSON response both use the struct's `json:"..."` tags. Rename a field in the struct without its tag, and the API stays the same; change the tag, and the API changes.
