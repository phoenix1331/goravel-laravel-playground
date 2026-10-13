---
title: "08. Middleware"
---

# 08. Middleware

| Laravel | Goravel |
|---|---|
| A class with `handle(Request $request, Closure $next)` | A struct with `Handle(ctx http.Context)`, calling `ctx.Request().Next()` |

Same idea, different signature. Laravel middleware receives the request and a `$next` closure, and returns a response. Goravel middleware receives the context and returns nothing: it either calls `Next()` to continue or writes a response and calls `Abort()` to stop.

The example rejects API writes that are not JSON, with `415 Unsupported Media Type`.

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch08/EnsureJsonRequest.php [EnsureJsonRequest.php]

<<< @/snippets/laravel/ch08/routes.php [Registering]

:::

## The Goravel way

`make artisan cmd="make:middleware EnsureJsonRequest"` scaffolds the struct:

<<< @/../app/http/middleware/ensure_json_request.go#middleware{go}

It is attached to the whole `/api` group in `routes/api.go`:

```go
facades.Route().Prefix("api").Middleware(middleware.EnsureJson()).Group(func(router route.Router) {
	// ...
})
```

`Middleware()` also works on a single route: `router.Middleware(middleware.EnsureJson()).Post("clubs", ...)`.

### Global middleware

Laravel's `$middleware->append()` in `bootstrap/app.php` has a direct equivalent in `bootstrap/app.go`:

```go
foundation.Setup().
	WithMiddleware(func(handler configuration.Middleware) {
		handler.Append(middleware.EnsureJson())
	}).
	// ...
```

This repo does not use it for `EnsureJson`, because `/health` and the fallback route are not part of the API.

### Before and after

Code before `Next()` runs before the handler; code after it runs once the handler has finished, the same as code after `$next($request)` in Laravel:

```go
func (m *Timing) Handle(ctx http.Context) {
	start := time.Now()

	ctx.Request().Next()

	facades.Log().Infof("%s took %s", ctx.Request().Path(), time.Since(start))
}
```

## Testing it

::: code-group

<<< @/snippets/laravel/ch08/EnsureJsonRequestTest.php [Laravel]

<<< @/../tests/feature/middleware_test.go#test{go} [Goravel]

:::

The first test checks the database as well as the status, to prove the controller never ran.

## What is different and why

- **No return value.** Laravel middleware must return the response, which is how `$next` threads the request down and the response back up. Goravel builds on Gin, where middleware and handlers share one context and the chain is managed for you. You signal "continue" with `Next()` and "stop" with `Abort()`.
- **Return after `Abort()`.** `Abort()` stops later handlers from running, but it does not stop the rest of your function. Without the `return`, the code would fall through to `request.Next()`. Forgetting this is the most common Goravel middleware bug.
- **Middleware is a value, not a class name.** Laravel accepts class names and aliases (`'auth'`, `EnsureJsonRequest::class`) and builds them through the container. Goravel takes an instance, usually from a small constructor such as `EnsureJson()`. To pass options, add parameters to the constructor instead of Laravel's `middleware:param` string syntax.
- **`Signature()` names it.** The name identifies the middleware in places such as `WithoutMiddleware`, much like a Laravel alias.
- **Laravel's defaults are not here.** Laravel's `api` group comes with middleware already attached (throttling, bindings). Goravel's groups start empty, so add what you need, such as `httpmiddleware.Throttle()`.
