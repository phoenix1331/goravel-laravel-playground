---
title: "16. Jobs and queues"
---

# 16. Jobs and queues

| Laravel | Goravel |
|---|---|
| `ShouldQueue` jobs, `::dispatch()`, `php artisan queue:work` | Job structs, `facades.Queue().Job(job, args).Dispatch()`, a worker inside the app process |

Chapter 15's listener wrote the notifications itself, inside the request that recorded the result. That work now moves into a queued job: the listener dispatches it, and a worker runs it later. In tests the queue is `sync`, so the job runs straight away; in development it is Redis.

The biggest surprise for a Laravel developer: there is no `queue:work`. The worker runs inside the same Go process as the web server.

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch16/SendResultNotifications.php [Job]

<<< @/snippets/laravel/ch16/NotifyClubAdmins.php [Listener]

<<< @/snippets/laravel/ch16/worker.sh [Running a worker]

:::

## The Goravel way

`make artisan cmd="make:job SendResultNotifications"` scaffolds the job and registers it in `bootstrap/jobs.go`.

<<< @/../app/jobs/send_result_notifications.go#job{go}

The listener from chapter 15 shrinks to a dispatch:

<<< @/../app/listeners/notify_club_admins.go#dispatch{go}

### Sync or Redis

The connection comes from `QUEUE_CONNECTION`, read in `config/queue.go`:

| Environment | `QUEUE_CONNECTION` | What `Dispatch()` does |
|---|---|---|
| Tests (`.env.testing`) | `sync` | Runs the job immediately, in the same call |
| Development (`.env`) | `redis` | Pushes the job to Redis; the worker picks it up |

`DispatchSync()` always runs immediately, whatever the connection, like `dispatch_sync()`.

### Where the worker runs

`app.Start()` in `main.go` starts every **runner** the installed packages provide: the HTTP server (`goravel:http`), the queue worker (`goravel:queue`) and the scheduler (`goravel:schedule`, chapter 17). The queue runner only starts when the connection is not `sync`. So with `QUEUE_CONNECTION=redis`, the `app` container is the web server and the worker at once:

```
app-1  | [HTTP] Listening on: http://0.0.0.0:8000
app-1  |   send_result_notifications .................. RUNNING
app-1  |   send_result_notifications ......... 20.443109ms DONE
```

To run workers separately, as you would with `queue:work`, start the same binary with the HTTP runner switched off. This repo reads the list from `APP_DISABLED_RUNNERS`:

<<< @/../config/app.go#runners{go}

```bash
APP_DISABLED_RUNNERS="goravel:http goravel:schedule" ./app   # a worker-only process
```

### Failed jobs

A job that returns an error is tried again up to the worker's `tries` limit (one by default), then recorded in `failed_jobs`. For finer control, a job can implement `ShouldRetry(err error, attempt int) (bool, time.Duration)`, the equivalent of Laravel's `$tries` and `backoff()`. `queue:failed` lists failed jobs and `queue:retry` re-queues them, as in Laravel.

## Testing it

::: code-group

<<< @/snippets/laravel/ch16/SendResultNotificationsTest.php [Laravel]

<<< @/../tests/feature/queue_test.go#test{go} [Goravel]

:::

The job is a plain struct with a `Handle` method, so the first test calls it directly. The chapter 15 tests still pass unchanged: on the sync connection, recording a result still produces notifications before the response is sent. There is no `Queue::fake()`; for "was it pushed?" assertions, Goravel's generated mocks can replace the queue facade.

## What is different and why

- **No `queue:work`.** A Go binary can run the web server and background workers as goroutines in one process, so Goravel does by default. One container gives you everything in development. In production you may still want separate processes, so web traffic and slow jobs scale independently; `APP_DISABLED_RUNNERS` gives you that from the same binary.
- **Jobs take arguments, not constructor properties.** A Laravel job serialises its public properties. A Goravel job is stateless: the payload is `[]queue.Arg`, each with a `Type` so the worker can rebuild it, and `Handle(args ...any)` receives the values. As with events, pass IDs and load fresh data in the job.
- **Jobs are registered by signature.** The worker receives a job name from Redis and looks it up in `bootstrap/jobs.go`. A job missing from that list cannot run on a real queue, even though it works on `sync`, because sync calls `Handle` directly. That makes it a mistake tests on `sync` will not catch; `make:job` registers it for you.
- **Graceful shutdown is built in.** When the process gets `SIGTERM`, `app.Start()` shuts runners down, letting the worker finish its current job. That is what Laravel's `queue:restart` and Supervisor stop signals are for.
