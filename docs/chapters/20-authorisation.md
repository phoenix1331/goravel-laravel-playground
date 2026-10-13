---
title: "20. Gates and policies"
---

# 20. Gates and policies

| Laravel | Goravel |
|---|---|
| Policies discovered per model, `Gate::authorize('update', $club)` | Policy methods registered as named gates, `facades.Gate().Inspect("update-club", args)` |

Chapter 19 made writes require a login. Now the rules about *which* writes:

| Action | Who may do it |
|---|---|
| Create a club | A user who does not run a club yet; they become its admin |
| Edit or delete a club, add players | That club's admin |
| Record a result | The home club's admin |

The pieces are the same as Laravel's: policies hold the rules, gates name them, controllers ask before acting. For the full API (`Before`, `After`, `Any`, `None`), see Goravel's [authorisation docs](https://www.goravel.dev/security/authorization.html).

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch20/ClubPolicy.php [ClubPolicy.php]

<<< @/snippets/laravel/ch20/ClubController.php [Using it]

:::

## The Goravel way

### The policy

`make artisan cmd="make:policy ClubPolicy"` scaffolds a struct. Each method takes a context and an arguments map, and returns an allow or deny response:

<<< @/../app/policies/club_policy.go#policy{go}

### Registering the gates

Each ability is named and pointed at a policy method in a provider's `Boot`:

<<< @/../app/providers/auth_service_provider.go#gates{go}

### Checking in the controller

A small helper loads the user, adds them to the arguments, and asks the gate:

<<< @/../app/http/controllers/authorise.go#authorise{go}

Each write action calls it before doing anything else:

<<< @/../app/http/controllers/club_controller.go#check{go}

Creating a club also attaches the creator as admin, inside a transaction so both happen or neither does:

<<< @/../app/http/controllers/club_controller.go#store{go}

## Testing it

::: code-group

<<< @/snippets/laravel/ch20/AuthorisationTest.php [Laravel]

<<< @/../tests/feature/authorisation_test.go#test{go} [Goravel]

:::

Policies are plain methods, so they can also be unit tested with no HTTP and no database:

::: code-group

<<< @/snippets/laravel/ch20/ClubPolicyTest.php [Laravel]

<<< @/../tests/unit/policy_test.go#test{go} [Goravel]

:::

## What is different and why

- **The gate does not know who is asking.** Laravel passes the logged-in user to every gate and policy method. Goravel's gate gets a context and whatever arguments you pass, so the `authorise` helper loads the user and puts them in the map as `"user"`.
- **Arguments are a `map[string]any`.** Policies pull values out with type assertions (`arguments["club"].(*models.Club)`) and must handle a missing or wrong-typed value. This project denies in that case: a programming mistake should fail closed, not open.
- **No policy discovery.** Laravel maps `Club` to `ClubPolicy` by naming convention. Goravel has string abilities (`"update-club"`) registered one by one. An ability nobody defined is denied with the message "ability doesn't exist: ...", so a typo in an ability name shows up as an unexpected 403.
- **No exceptions for 403.** `Gate::authorize()` throws and the exception handler renders a 403. `Inspect()` returns a response the helper turns into a 403 itself. `Allows()` and `Denies()` exist for plain yes or no checks.
- **Form requests have `Authorize` too.** Goravel's form requests (chapter 07) can return an error from `Authorize`, but `ValidateRequest` reports it as an ordinary error, not a 403. This project keeps authorisation in the controller so a denial is always a 403 with a clear message.
