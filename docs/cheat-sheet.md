# Laravel to Goravel cheat sheet

One row per concept. Each chapter link has the full side-by-side comparison.

| Chapter | Laravel | Goravel |
|---|---|---|
| [01. Installation and structure](/chapters/01-installation) | `laravel new`, `artisan`, directory structure | `goravel new`, `./artisan`, layout |
| [02. Configuration](/chapters/02-configuration) | `.env` and `config/*.php` | `.env`, `config/*.go`, `facades.Config()` |
| [03. Routing](/chapters/03-routing) | `routes/web.php`, `routes/api.php` | `routes/*.go`, route groups |
| [04. Controllers](/chapters/04-controllers) | Controller classes | Controller structs |
| [05. Service container and providers](/chapters/05-service-container) | Service container, `register()` / `boot()` | Container bindings, `Register` / `Boot` |
| [06. Facades](/chapters/06-facades) | Static proxies via `__callStatic` | `app/facades` package functions |
| [07. Form requests and validation](/chapters/07-validation) | Form requests, `$request->validate()` | Form requests, validation rules |
| [08. Middleware](/chapters/08-middleware) | Middleware classes with `handle()` | Middleware functions |
| [09. ORM models](/chapters/09-models) | Eloquent models | ORM models |
| [10. Relationships](/chapters/10-relationships) | Eloquent relationships, `with()` | ORM relationships, `With()` |
| [11. Migrations](/chapters/11-migrations) | Migrations with the schema builder | Migrations in Go |
| [12. Factories and seeders](/chapters/12-factories-seeders) | Model factories and seeders | Factories and seeders |
| [13. Query builder and raw SQL](/chapters/13-query-builder) | `DB::table()`, `DB::select()` | `facades.Orm().Query()`, `facades.DB()` |
| [14. Artisan commands](/chapters/14-console-commands) | `make:command`, `handle()` | Console commands (`make:command`) |
| [15. Events and listeners](/chapters/15-events) | Events, listeners, `event()` | Events and listeners |
| [16. Jobs and queues](/chapters/16-queues) | Jobs, `dispatch()`, queue workers | Jobs and queues |
| [17. Task scheduling](/chapters/17-scheduling) | `routes/console.php` schedule | Scheduler |
| [18. Cache](/chapters/18-cache) | `Cache::remember()` | Cache facade |
| [19. Authentication](/chapters/19-authentication) | Session guard, Sanctum | Auth with the JWT guard |
| [20. Gates and policies](/chapters/20-authorisation) | Gates and policies | Authorisation |
| [21. HTTP client](/chapters/21-http-client) | `Http::get()`, `Http::fake()` | HTTP client |
| [22. Logging](/chapters/22-logging) | `Log::info()`, channels | Logging |
| [23. Testing](/chapters/23-testing) | PHPUnit / Pest, HTTP tests, fakes | testify suites, HTTP tests, mocks |
| [24. Error handling](/chapters/24-error-handling) | Exceptions, `try` / `catch`, handler | Go errors |
