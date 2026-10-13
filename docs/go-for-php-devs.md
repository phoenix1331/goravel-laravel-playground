# Go for PHP devs

Just enough Go to read the code in this playground. For the language itself, work through [A Tour of Go](https://go.dev/tour/) (an afternoon) and keep [Effective Go](https://go.dev/doc/effective_go) nearby.

## Running code

PHP is interpreted per request. Go compiles the whole program into one binary, which then runs until stopped.

| PHP | Go |
|---|---|
| `php artisan serve` | `go run .` (compile and run) or `go build` (just compile) |
| Edit, refresh | Edit, rebuild: `air` does this on save ([chapter 01](/chapters/01-installation)) |
| Errors found when the line runs | Type errors, typos and unused imports stop the build |

## Packages and visibility

Every directory is a package, and every file starts with `package name`. Code imports other packages by path: `import "goravel/app/models"`.

There are no `public` or `private` keywords. **A capital first letter means exported** (usable from other packages); lower case means private to the package:

```go
func Record(...)        // results.Record: callable from anywhere
func findOpenFixture()  // only inside package commands
type Club struct {
	Name string           // exported field: the ORM and JSON encoder can see it
	secret string         // unexported: invisible to both
}
```

This catches out every PHP developer once: a lower-case struct field silently never reaches the database or the JSON.

## Variables and types

```go
var name string = "Ashford"   // explicit
name := "Ashford"             // declare and infer, inside functions
count := 0                    // int
```

Every type has a **zero value** instead of `null`: `0`, `""`, `false`, or `nil` for pointers, slices, maps and interfaces. To represent "no value" for a number or string, use a pointer (`*int`), which can be `nil` ([chapter 09](/chapters/09-models)).

| PHP | Go |
|---|---|
| `array` (list) | slice: `[]string{"a", "b"}`, `append(s, "c")` |
| `array` (keys) | map: `map[string]int{"win": 3}` |
| `foreach ($items as $i => $item)` | `for i, item := range items` |
| `mixed` | `any` |

## Structs and methods instead of classes

```go
type Standing struct {
	Points       int
	GoalsFor     int
	GoalsAgainst int
}

// A method: a function with a receiver
func (s Standing) GoalDifference() int {
	return s.GoalsFor - s.GoalsAgainst
}
```

No inheritance. Embedding (`orm.Model` inside `Club`) copies the embedded struct's fields and methods in, which covers most of what `extends` was used for. Constructors are ordinary functions, by convention `NewThing(...)`.

## Pointers

`&club` is the address of `club`; `*models.Club` is the type "pointer to a Club". Pass a pointer when the function must change the value or fill it in:

```go
var club models.Club
facades.Orm().Query().Create(&club) // fills club.ID; without & it would fill a copy
```

Methods that change their struct use a pointer receiver: `func (c *Club) Rename(...)`.

## Interfaces are implicit

```go
type FixtureFeed interface {
	Results() ([]FeedResult, error)
}
```

Any type with a `Results() ([]FeedResult, error)` method satisfies `FixtureFeed`. There is no `implements` keyword ([chapter 05](/chapters/05-service-container)).

## Multiple return values and errors

Functions can return several values, and the last is usually an `error`:

```go
results, err := feed.Results()
if err != nil {
	return err
}
```

No exceptions: errors are values you check and return. [Chapter 24](/chapters/24-error-handling) covers wrapping and `errors.Is`.

## `defer`

`defer` runs a call when the surrounding function returns, whatever path it takes. It replaces `finally`:

```go
file, err := os.Open(path)
if err != nil {
	return err
}
defer file.Close()
```

## Goroutines, briefly

`go doSomething()` runs a function concurrently. You rarely write this in a Goravel app, but it explains why one process can serve HTTP, run the queue worker and the scheduler at once ([chapter 16](/chapters/16-queues)), and why anything stored in a package-level variable is shared by every request.

## Tooling

| PHP | Go |
|---|---|
| Composer | `go mod`, `go get` (built in) |
| PHPUnit / Pest | `go test` (built in) |
| Pint / PHP-CS-Fixer | `gofmt` (built in, one style, no options) |
| PHPStan / Psalm | the compiler, plus `go vet` |
