---
title: "21. HTTP client"
---

# 21. HTTP client

| Laravel | Goravel |
|---|---|
| `Http::get()`, `Http::fake()`, `Http::assertSent()` | `facades.Http().Get()`, `facades.Http().Fake()`, `facades.Http().AssertSent()` |

Chapter 14's `sync:fixtures` read results from a JSON file. Now it can also call the league's results API. Goravel 1.18's HTTP client is modelled closely on Laravel's, fakes included, so this is one of the shortest chapters. The [HTTP client docs](https://www.goravel.dev/digging-deeper/http-client.html) have the full API.

## The Laravel way

<<< @/snippets/laravel/ch21/HttpFeed.php

## The Goravel way

A named client in `config/http.go` holds the base URL and timeout:

<<< @/../config/http.go#client{go}

`HttpFeed` sits next to chapter 14's `FileFeed` and satisfies the same `FixtureFeed` interface:

<<< @/../app/services/fixture_feed.go#http{go}

The command picks a feed from `--file` or `FIXTURE_FEED_DRIVER`:

<<< @/../app/console/commands/sync_fixtures.go#feed{go}

### Trying it locally

There is no real league API, so the app serves a stand-in at `GET /api/demo/results`, returning the same JSON as the file feed. The `fixtures` client points at it by default:

```bash
docker compose exec -e FIXTURE_FEED_DRIVER=http app go run . artisan sync:fixtures --dry-run
```

## Testing it

::: code-group

<<< @/snippets/laravel/ch21/HttpFeedTest.php [Laravel]

<<< @/../tests/feature/http_client_test.go#test{go} [Goravel]

:::

`.env.testing` points the client at `https://feed.test/api`, a domain that does not exist, so a test that forgets to fake it fails rather than calling anything real. `PreventStrayRequests()` makes that explicit.

The fakes and the config change are process-wide, so `TearDownTest` calls `facades.Http().Reset()` and puts the feed driver back. Laravel resets fakes between tests for you.

## What is different and why

- **4xx and 5xx are not errors.** `Get` only returns an error when the request did not complete (DNS, connection refused, timeout). A 503 is a successful request with a failed status, so check `response.Failed()`. Laravel's client behaves the same until you call `->throw()`; Go code checks the status every time.
- **The scaffold's facade was too narrow.** The generated `app/facades/http.go` returned `client.Request`, which has no `Fake()`. The container actually holds a `client.Factory` (a request builder plus the fake and assert methods), so this project changed the facade's return type. If you cannot find `Fake()` in your own app, check the facade.
- **Inject the client, do not reach for the facade.** `HttpFeed` receives a `client.Request` rather than calling `facades.Http()` itself. That keeps `services` free of framework lookups (it cannot import `facades` anyway: `facades` already imports `services`, and Go forbids import cycles).
- **Patterns are URL wildcards or client names.** `"feed.test/*"` matches any URL on that host. `"fixtures#/results"` matches the path on the named client only.
