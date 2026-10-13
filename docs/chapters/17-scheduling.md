---
title: "17. Task scheduling"
---

# 17. Task scheduling

| Laravel | Goravel |
|---|---|
| `routes/console.php`, `Schedule::command()`, a cron entry running `schedule:run` | `WithSchedule()` in `bootstrap/app.go`, `facades.Schedule().Command()`, a runner inside the app |

The fluent API is the same, almost method for method. What changes is how the schedule runs. Laravel needs cron to call `schedule:run` every minute, or a `schedule:work` process. Goravel's scheduler is another runner started by `app.Start()` ([chapter 16](/chapters/16-queues)), so the app container already runs it.

This chapter schedules chapter 14's `sync:fixtures` command for match days, plus a nightly clean-up task.

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch17/console.php [routes/console.php]

<<< @/snippets/laravel/ch17/crontab [Running it]

:::

## The Goravel way

This repo keeps the schedule in `bootstrap/schedule.go`, next to the other lists (`commands.go`, `jobs.go`, `migrations.go`):

<<< @/../bootstrap/schedule.go#schedule{go}

It is handed to the app builder with `WithSchedule(Schedule)` in `bootstrap/app.go`.

`schedule:list` shows what is registered and when it is next due:

```
$ make artisan cmd="schedule:list"

  */15 *  * * 6  artisan sync:fixtures ................ Next Due: 3 days after
  00   03 * * *  prune-old-notifications ............. Next Due: 7 hours after
```

| Laravel | Goravel |
|---|---|
| `->withoutOverlapping()` | `.SkipIfStillRunning()` (or `.DelayIfStillRunning()` to queue the next run) |
| `->onOneServer()` | `.OnOneServer()`, using a cache lock |
| `->name('...')` | `.Name("...")`, required for closures so they show up sensibly |
| `schedule:run` | `schedule:run` exists, for running due tasks once |
| `schedule:work` | Not needed: the runner is always on |

### Running it in the container

There is nothing to add to `docker-compose.yml`. The `app` container runs `air`, air runs the binary, and `app.Start()` starts the HTTP server, the queue worker and the scheduler together. Times are in the app timezone, UTC (`config/app.go`), so "Saturday" means Saturday in UTC.

## Testing it

::: code-group

<<< @/snippets/laravel/ch17/ScheduleTest.php [Laravel]

<<< @/../tests/feature/schedule_test.go#test{go} [Goravel]

:::

The tests read the registered events back with `facades.Schedule().Events()` and check their cron expressions. Tests call `bootstrap.Boot()` but never `Start()`, so no runner starts and nothing fires on its own. The prune test grabs the closure with `GetCallback()` and calls it directly.

`DailyAt("03:00")` produces the cron expression `00 03 * * *`, with leading zeros. The first version of the test expected `0 3 * * *` and failed. Both mean 3am, but the strings differ.

## What is different and why

- **No cron entry.** Go programs are long-running, so the scheduler is a goroutine in the app process that checks the schedule every second. Deploy the binary and scheduling works.
- **Every instance runs the scheduler.** This is the catch. Run three copies of the app behind a load balancer and `sync:fixtures` runs three times. Either add `.OnOneServer()` to each task (it takes a lock in the shared Redis cache) or disable the scheduler on all but one instance with `APP_DISABLED_RUNNERS="goravel:schedule"`. Laravel has the same problem with cron on several servers, but there you add the cron entry to one server on purpose; here it happens by default.
- **Second-level schedules.** `EverySecond()` to `EveryThirtySeconds()` exist because the runner ticks every second. Laravel added sub-minute tasks in version 10.
- **Closures are Go functions.** `Call(pruneOldNotifications)` takes a `func()`. It cannot return an error, so the task logs its own failures.
