---
title: "06. Facades"
---

# 06. Facades

| Laravel | Goravel |
|---|---|
| `Cache::get()`, static proxies via `__callStatic` | `facades.Cache().Get()`, plain functions in `app/facades` |

Goravel keeps the name and the purpose: short, global access to a service in the container. The mechanism is completely different. A Goravel facade is an ordinary function that resolves the service and returns it, typed. There is no static class and no magic method.

## The Laravel way

<<< @/snippets/laravel/ch06/LeagueTable.php

`LeagueTable::standings($results)` is not a real static method. PHP calls `__callStatic('standings', ...)` on the `Facade` base class, which resolves the accessor from the container and forwards the call. The `@method` docblock exists only so your IDE knows the method is there.

## The Goravel way

Every facade in the app, built-in or custom, lives in `app/facades` as one small file. The built-in cache facade is the whole file:

```go
package facades

import "github.com/goravel/framework/contracts/cache"

func Cache() cache.Cache {
	return App().MakeCache()
}
```

A custom facade for the chapter 05 league table follows the same pattern:

<<< @/../app/facades/league_table.go#facade{go}

Callers write `facades.LeagueTable().Standings(results)`. The extra `()` is the giveaway: you call a function, it returns a `services.LeagueTable`, and you call a method on that. Your editor knows the return type, so there is no docblock to maintain.

### Installable facades

Since Goravel 1.17 the facades are files in your app rather than functions inside the framework. That is why `app/facades` exists, and why you can add your own next to them. The `goravel-lite` scaffold starts with only `App`, `Artisan`, `Config` and `Process`, and you add the rest as you need them:

```bash
./artisan package:install Cache       # adds app/facades/cache.go, config/cache.go, the provider
./artisan package:uninstall Cache     # removes them again
```

This repo uses the full scaffold, which installs every facade up front.

## Testing it

::: code-group

<<< @/snippets/laravel/ch06/LeagueTableFacadeTest.php [Laravel]

<<< @/../tests/unit/facade_test.go#test{go} [Goravel]

:::

The fake is any type with a `Standings` method:

<<< @/../tests/unit/facade_test.go#fake{go}

There is no `shouldReceive()`. Instead, put a different instance in the container with `Instance`, the same idea as `Facade::swap()`. Two Goravel-specific details:

1. **Call `Fresh` after `Instance`.** The binding is a singleton, and once it has been resolved the container caches the instance. `Instance` alone changes the binding but not the cache, so the facade keeps returning the real service. This project's first attempt failed for exactly that reason.
2. **Restore it yourself.** Laravel resets facades between tests. In Go every test in a package shares one process and one container, so `t.Cleanup` puts the original back.

For framework facades such as `Cache` or `Orm`, Goravel ships generated mocks in `github.com/goravel/framework/testing/mock` (`mock.Factory().Cache()`), which are closer to `Cache::shouldReceive()`. Chapter 23 uses them.

## What is different and why

- **A facade is a function, not a class.** Go has no static methods and no `__callStatic`, so there is nothing to proxy through. A function that returns the service gives the same short syntax with full type checking.
- **Typed all the way.** `facades.LeagueTable()` returns `services.LeagueTable`. Call a method that does not exist and the build fails. Laravel finds out at runtime, or relies on IDE helper docblocks.
- **Resolution happens on every call.** Like Laravel, each call asks the container. For a singleton that is a cheap map lookup.
- **A failed lookup returns `nil`.** The built-in pattern logs the error and returns `nil`, and calling a method on `nil` panics. In practice this only happens if a provider is missing from `bootstrap/providers.go`, and you find out on the first call.
- **Dependency injection is still an option.** Facades are convenient in controllers. For code you want to unit test without the container, pass the `services.LeagueTable` in as a constructor argument instead (chapter 05).
