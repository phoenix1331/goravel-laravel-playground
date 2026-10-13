---
title: "07. Form requests and validation"
---

# 07. Form requests and validation

| Laravel | Goravel |
|---|---|
| `FormRequest` classes, `$request->validated()` | Form request structs, `ctx.Request().ValidateRequest()` |

Goravel 1.18's validator speaks Laravel: the rule strings are almost identical (`required|string|max:100`, `nullable`, `unique`, `min`, `max`), and form requests have `Authorize`, `Rules`, `Messages`, `Attributes` and `PrepareForValidation`. The difference is in the plumbing. Laravel validates before your controller runs and throws on failure. In Goravel the controller asks for validation and gets the errors back as a value.

This chapter replaces the unvalidated `Bind` from [chapter 04](/chapters/04-controllers).

## The Laravel way

::: code-group

<<< @/snippets/laravel/ch07/StoreClubRequest.php [StoreClubRequest.php]

<<< @/snippets/laravel/ch07/store.php [ClubController::store]

:::

## The Goravel way

`make artisan cmd="make:request StoreClubRequest"` scaffolds the file. The request is a struct whose fields receive the validated data, plus methods for the rules:

<<< @/../app/http/requests/store_club_request.go#request{go}

The rules live in a shared function so the update request can reuse them:

<<< @/../app/http/requests/store_club_request.go#rules{go}

The controller runs the validation itself and decides what to do with the result:

```go
// Store creates a club from a validated request.
func (r *ClubController) Store(ctx http.Context) http.Response {
	var request requests.StoreClubRequest
	// Two results: validationErrors for bad input (422), err for anything
	// that stopped validation running at all, such as an unreadable body
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": err.Error()})
	}
	if validationErrors != nil {
		return validationFailed(ctx, validationErrors)
	}

	club := models.Club{Name: request.Name, Ground: request.Ground, FoundedYear: request.FoundedYear}
	if err := facades.Orm().Query().Create(&club); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{"data": club})
}
```

[Chapter 20](/chapters/20-authorisation) adds an authorisation check before the validation.

<<< @/../app/http/controllers/club_controller.go#failed{go}

A failed request looks like this:

```json
{
  "message": "The given data was invalid.",
  "errors": {
    "founded_year": { "min": "The oldest football club was founded in 1857." },
    "ground": { "required": "The ground field is required." },
    "name": { "required": "The name field is required." }
  }
}
```

### Ignoring the current record in `unique`

The update request ignores the club being edited, so saving without renaming does not fail. The argument order is **not** Laravel's:

::: code-group

<<< @/snippets/laravel/ch07/UpdateClubRequest.php [Laravel]

```go [Goravel]
// unique:table,column,idColumn,ignoredValue
"unique:clubs,name,id,7"
```

:::

Copy a Laravel rule across unchanged and it silently checks the wrong column. The tests cover both directions: saving a club under its own name passes, and taking another club's name fails.

## Testing it

::: code-group

<<< @/snippets/laravel/ch07/ClubValidationTest.php [Laravel]

<<< @/../tests/feature/club_validation_test.go#test{go} [Goravel]

:::

There is no `assertJsonValidationErrors`. The test helper decodes the 422 body into `map[string]map[string]string`, and the assertions check for a field and rule name, such as `errors["name"]["unique"]`.

## What is different and why

- **No exceptions, so validation is a call you make.** Laravel resolves the form request from the type hint, validates it, and throws `ValidationException` before your method runs. Go has no exceptions to unwind with. `ValidateRequest` returns `(validation.Errors, error)` and the controller handles both. It is more code per action, but the whole flow is visible.
- **Two kinds of failure.** `validationErrors` means the input was bad (422). `err` means validation could not run, for example `Authorize` refused or the body was unreadable. Keep them separate.
- **The request struct is the validated data.** After the rules pass, Goravel binds the input into the struct's fields, typed. `request.FoundedYear` is an `*int` (nil when absent), not a string you cast. That is the equivalent of `$request->validated()`.
- **Trimming is opt-in.** Laravel's global `TrimStrings` middleware trims every input. Goravel does not, so `"   "` would pass `required`. `Filters()` with `"trim"` fixes that per field.
- **Error shape.** Laravel returns a list of messages per field. Goravel's `Errors().All()` returns a map of rule name to message per field. Either is fine for an API, but clients written against Laravel's shape will need a small change.
- **Rules are code.** `time.Now().Year()` in the rules is an ordinary function call. Laravel can do the same in `rules()`, so this one is a similarity: both build rule strings at runtime.
