---
title: "19. Authentication"
---

# 19. Authentication

| Laravel | Goravel |
|---|---|
| Session guard by default; Sanctum tokens for APIs | JWT guard by default, `facades.Auth(ctx)` |

Reads stay public. Writes now need a bearer token, which `POST /api/auth/login` hands out. This is a short tour of the code; Goravel's [authentication docs](https://www.goravel.dev/security/authentication.html) cover guards, providers and refresh in full.

```bash
curl -X POST localhost:8000/api/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"admin1@example.com","password":"password"}'
# {"access_token":"eyJhbGciOi...","expires_in":3600,"token_type":"Bearer"}
```

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch19/AuthController.php [AuthController.php]

<<< @/snippets/laravel/ch19/api.php [routes/api.php]

:::

## The Goravel way

### Logging in

<<< @/../app/http/controllers/auth_controller.go#login{go}

`config/auth.go` sets the default guard to `jwt`, and `make artisan cmd="jwt:secret"` writes `JWT_SECRET` to `.env`. Tokens are signed with that secret and last `JWT_TTL` minutes (60 by default).

### Protecting routes

The middleware parses the token. Parsing also stores the user's ID in the request context, which is how `facades.Auth(ctx).User(&user)` knows who is asking:

<<< @/../app/http/middleware/authenticate.go#middleware{go}

The routes are now split into public reads and authenticated writes:

<<< @/../routes/api.go#authenticated{go}

<<< @/../app/http/controllers/auth_controller.go#me{go}

## Testing it

::: code-group

<<< @/snippets/laravel/ch19/AuthTest.php [Laravel]

<<< @/../tests/feature/auth_test.go#test{go} [Goravel]

:::

There is no `Sanctum::actingAs()`. Goravel's JWT guard needs a real HTTP request context to issue a token, so `tests.TestCase` gets a `LoginAs` helper that logs in through the real endpoint and returns a request builder carrying the token:

<<< @/../tests/test_case.go#login{go}

Every earlier test that writes data now calls `s.LoginAs(s.T(), user)`. Logging in for every test costs a bcrypt check, so `.env.testing` sets `BCRYPT_ROUNDS=4` and the user factory hashes `password` once per process, the same two tricks Laravel uses.

## What is different and why

- **JWT, not sessions.** Laravel's default is a session cookie, and Sanctum adds API tokens stored in a database table. Goravel defaults to signed JWTs: nothing is stored on login, and any server with the secret can check a token.
- **No `Auth::attempt()`.** You find the user and call `facades.Hash().Check()` yourself, then `Login(&user)` issues the token.
- **`Auth` takes the request context.** `facades.Auth(ctx)`, every time. The guard reads and writes the current user on that context, not on a global.
- **Logout disables the token, and the token is not unique.** JWTs cannot be recalled, so `Logout()` adds the token to a block list in the cache until it would have expired. Goravel's tokens contain only the user ID and timestamps, with no random ID (`jti`), so **the same user logging in twice in the same second gets the identical token**. This project found that out when the logout test disabled the token another test was using. Sanctum tokens are random, so it cannot happen there.
- **Same message for both failures.** An unknown email and a wrong password both return "These credentials do not match our records.", so the endpoint does not reveal which emails exist.
- **Authentication is not authorisation.** Any logged-in user can currently edit any club. [Chapter 20](/chapters/20-authorisation) fixes that.
