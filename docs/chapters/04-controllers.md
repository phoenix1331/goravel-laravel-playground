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

<<< @/../routes/api.go#resource{go}

A handler method looks like the closures from chapter 03, with a receiver `(r *ClubController)` in front:

<<< @/../app/http/controllers/club_controller.go#index{go}

### Standing in for route model binding

Laravel's `show(Club $club)` does three jobs: reads `{club}` from the URL, loads the row, and turns "not found" into a 404. In Goravel you write that once as a helper and call it from each method:

<<< @/../app/http/controllers/club_controller.go#find{go}

<<< @/../app/http/controllers/club_controller.go#show{go}

The `(*models.Club, http.Response)` return is a common Go shape: return the value or a reason to stop, and let the caller check which. Chapter 24 has more on returning errors instead of throwing them.

### Reading the request body

<<< @/../app/http/controllers/club_controller.go#store{go}

`$request->only(...)` becomes a struct describing the expected JSON. `Bind` decodes the body into it. Fields not in the struct are ignored, which gives you the same protection as `$fillable` without a separate list. There is no validation yet: an empty name is accepted. [Chapter 07](/chapters/07-validation) fixes that with a form request.

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
