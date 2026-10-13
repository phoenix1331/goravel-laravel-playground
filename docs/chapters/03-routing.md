---
title: "03. Routing"
---

# 03. Routing

| Laravel | Goravel |
|---|---|
| `routes/web.php`, `routes/api.php`, `Route::get()` | `routes/*.go`, `facades.Route().Get()` |

Routes read almost the same. A route maps a method and path to a handler, groups share a prefix, and a fallback catches everything else. The Go-specific part is what a handler is: any function with the signature `func(http.Context) http.Response`.

## The Laravel way

<<< @/snippets/laravel/ch03/api.php

Laravel loads `routes/api.php` because `bootstrap/app.php` names it in `withRouting()`, and adds the `/api` prefix for you.

## The Goravel way

<<< @/../routes/api.go#api{go}

The routes themselves are split into small functions that each receive the group's `router`. The league routes show the two kinds of handler:

<<< @/../routes/api.go#league{go}

The named handler is an ordinary package-level function:

<<< @/../routes/api.go#handler{go}

Route files are registered in `bootstrap/app.go`. Add a new file and forget this step, and its routes silently do not exist:

```go
WithRouting(func() {
	routes.Web()
	routes.Api()   // added for this chapter
	routes.Grpc()
}).
```

There is no automatic `/api` prefix either. `Prefix("api").Group(...)` adds it explicitly, and you can see it in the code.

Check what is registered with `make artisan cmd="route:list"`, the same as `php artisan route:list`. At this point in the project it showed:

```
GET|HEAD  /                  goravel/routes.Web.func1
GET|HEAD  api/league         goravel/routes.Api.func1.1
GET|HEAD  api/league/points  goravel/routes.leaguePoints
GET|HEAD  health             goravel/app/http/controllers.(*HealthController).Show
```

## The fallback route

::: code-group

<<< @/snippets/laravel/ch03/fallback.php [Laravel]

<<< @/../routes/web.go#fallback{go} [Goravel]

:::

This one caught the project out. With the Gin driver and no fallback, an unknown URL returns **200 with an empty body**, not 404. Global middleware (timeout, CORS) runs on unmatched requests, and nothing sets an error status. Laravel returns 404 without being asked, so in Goravel register a fallback in every app.

## Testing it

::: code-group

<<< @/snippets/laravel/ch03/RoutingTest.php [Laravel]

<<< @/../tests/feature/routing_test.go#test{go} [Goravel]

:::

Note the `3.0` in the Go test. When Go decodes JSON into `map[string]any`, every number becomes a `float64`, because JSON has a single number type. `3` (an `int`) and `3.0` (a `float64`) are not equal in Go, so the assertion has to use floats. PHP's `json_decode` picks `int` or `float` per value, so the Laravel test never has to think about it.

## What is different and why

- **A handler is a function type, not a callable string or array.** Laravel accepts closures, `[Controller::class, 'method']` arrays and invokable class names, and resolves them at runtime. Goravel accepts one thing: a value of type `func(http.Context) http.Response`. Closures, named functions and controller methods (chapter 04) all satisfy it, and the compiler rejects anything else.
- **No route model binding.** `Route::get('/clubs/{club}', ...)` with a type-hinted `Club $club` has no equivalent. Goravel gives you the raw parameter with `ctx.Request().Route("id")` and you load the model yourself. That is chapter 04.
- **Groups take a function.** `Route::prefix('api')->group(function () { ... })` becomes `Prefix("api").Group(func(router route.Router) { ... })`. Inside the group you register on the `router` you are handed, not on the facade.
- **Registration is explicit.** No file is loaded by convention. If `routes.Api()` is not called in `bootstrap/app.go`, nothing in `routes/api.go` exists.
