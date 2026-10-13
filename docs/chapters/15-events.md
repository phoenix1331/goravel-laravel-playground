---
title: "15. Events and listeners"
---

# 15. Events and listeners

| Laravel | Goravel |
|---|---|
| Event classes holding models, `event()` / `::dispatch()`, discovered listeners | Event structs, `facades.Event().Job(event, args).Dispatch()`, listeners mapped in `bootstrap/app.go` |

When a fixture gets a result, two things should follow: the clubs' admins hear about it (this chapter), and the cached league table is cleared ([chapter 18](/chapters/18-cache)). Neither belongs in the controller or the sync command. An event lets the code that records the result say "this happened" without knowing who is listening.

The biggest difference from Laravel: a Goravel event does not carry objects. It carries a list of typed arguments, and here that is just the fixture ID.

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch15/FixtureResultRecorded.php [Event]

<<< @/snippets/laravel/ch15/NotifyClubAdmins.php [Listener]

<<< @/snippets/laravel/ch15/RecordResult.php [Dispatching]

:::

## The Goravel way

`make:event FixtureResultRecorded` and `make:listener NotifyClubAdmins` scaffold both files.

### The event

<<< @/../app/events/fixture_result_recorded.go#event{go}

`Args` is a helper this project adds. Goravel's own API takes `[]event.Arg` everywhere, and a helper on the event keeps the argument list in one place.

### The listener

<<< @/../app/listeners/notify_club_admins.go#listener{go}

### Wiring them together

Laravel discovers listeners by the type hint on `handle()`. Goravel needs an explicit map, passed to the app builder in `bootstrap/app.go`:

<<< @/../bootstrap/app.go#events{go}

### Dispatching

`results.Record` from chapter 14 now fires the event after saving, so the sync command and the new `POST /api/fixtures/{id}/result` endpoint both trigger it:

<<< @/../app/results/record.go#dispatch{go}

<<< @/../app/http/controllers/fixture_controller.go#record{go}

The notifications are rows in a new `result_notifications` table.

## Testing it

::: code-group

<<< @/snippets/laravel/ch15/FixtureResultRecordedTest.php [Laravel]

<<< @/../tests/feature/events_test.go#test{go} [Goravel]

:::

These tests check the outcome, notifications in the database, rather than "the event was dispatched". Goravel has no `Event::fake()` with `assertDispatched`. Its generated mocks (`mock.Factory().Event()`) can replace the event facade, but they assert on the exact `Job(...)` call, which ties the test to implementation details. Checking what the listener did is simpler and survives refactoring.

## What is different and why

- **Arguments, not objects.** A Laravel event holds `public Fixture $fixture`, and `SerializesModels` reduces it to an ID if it goes on a queue. A Goravel event's payload is `[]event.Arg{{Type: "uint", Value: id}}`. The `Type` string tells the queue how to rebuild the value, so the same event works with synchronous and queued listeners. Listeners load what they need from the ID.
- **Listeners receive `...any`.** `Handle(args ...any)` gets the values in dispatch order, untyped. `cast.ToUint(args[0])` converts safely, whatever numeric type arrived. Check `len(args)` before indexing: a wrong dispatch call is a runtime error, not a compile error.
- **Explicit registration.** No discovery. A listener missing from the `WithEvents` map never runs, and dispatching an event with no listeners at all returns an error (`EventListenerNotBind`).
- **`Queue()` decides sync or queued.** `Enable: false` runs the listener inside `Dispatch()`, like a Laravel listener without `ShouldQueue`. Chapter 16 changes this.
- **Errors stop the chain.** If a synchronous listener returns an error, `Dispatch()` returns it and later listeners do not run. `results.Record` passes that error back to the controller, which returns a 500. The result is already saved by then, so a listener failure does not lose the score.
