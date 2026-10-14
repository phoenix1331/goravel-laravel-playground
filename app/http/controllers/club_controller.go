// Package controllers holds the HTTP handlers.
//
// Laravel: app/Http/Controllers, one class per controller, resolved by the
// container with constructor injection and route model binding.
// Goravel: one struct per controller, built with a New... function, and
// every lookup, validation and authorisation step written out in the method.
package controllers

import (
	"errors"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	frameworkerrors "github.com/goravel/framework/errors"

	"goravel/app/facades"
	"goravel/app/http/requests"
	"goravel/app/models"
)

// ClubController manages clubs.
//
// Laravel:
//
//	class ClubController extends Controller
//	{
//	    public function __construct(private LeagueTable $table) {}  // injected
//	}
//
// Goravel: a plain struct. Its methods are route handlers because they have
// the signature func(http.Context) http.Response. To give it a dependency,
// add a field and pass it into NewClubController yourself.
type ClubController struct{}

// NewClubController is a constructor by convention only: Go has no
// __construct and nothing calls this for you. routes/api.go calls it once.
func NewClubController() *ClubController {
	return &ClubController{}
}

// Index lists every club.
//
// Laravel:
//
//	return response()->json(['data' => Club::orderBy('name')->get()]);
//
// Goravel: Get fills a slice you pass in, and returns an error you check.
func (r *ClubController) Index(ctx http.Context) http.Response {
	var clubs []models.Club
	if err := facades.Orm().Query().Order("name").Get(&clubs); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": clubs})
}

// Show returns one club, or 404.
//
// Laravel:
//
//	public function show(Club $club): JsonResponse  // route model binding
//	{
//	    return response()->json(['data' => $club]);
//	}
//
// Goravel: no route model binding, so findClub reads {id} and loads the row.
func (r *ClubController) Show(ctx http.Context) http.Response {
	club, response := findClub(ctx)
	if response != nil {
		return response
	}

	return ctx.Response().Success().Json(http.Json{"data": club})
}

// Store creates a club and makes the logged-in user its admin.
//
// Laravel:
//
//	public function store(StoreClubRequest $request): JsonResponse
//	{
//	    Gate::authorize('create', Club::class);
//	    $club = DB::transaction(function () use ($request) {
//	        $club = Club::create($request->validated());
//	        $request->user()->update(['club_id' => $club->id]);
//	        return $club;
//	    });
//	    return response()->json(['data' => $club], 201);
//	}
//
// Goravel: authorisation and validation are calls you make and check, not
// steps that run before the method and throw on failure.
func (r *ClubController) Store(ctx http.Context) http.Response {
	user, response := authorise(ctx, "create-club", map[string]any{})
	if response != nil {
		return response
	}

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

	// Create the club and attach the admin together, or not at all, like
	// DB::transaction(fn () => ...). Returning an error rolls back.
	club := models.Club{Name: request.Name, Ground: request.Ground, FoundedYear: request.FoundedYear}
	err = facades.Orm().Transaction(func(tx orm.Query) error {
		if err := tx.Create(&club); err != nil {
			return err
		}
		user.ClubID = &club.ID

		return tx.Save(user)
	})
	if err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{"data": club})
}

// Update replaces a club's fields from a validated request.
//
// Laravel:
//
//	public function update(UpdateClubRequest $request, Club $club): JsonResponse
//	{
//	    Gate::authorize('update', $club);
//	    $club->update($request->validated());
//	    return response()->json(['data' => $club]);
//	}
//
// Goravel: models have no update() or save(); pass the model to the query
// builder's Save.
func (r *ClubController) Update(ctx http.Context) http.Response {
	club, response := findClub(ctx)
	if response != nil {
		return response
	}
	if _, response := authorise(ctx, "update-club", map[string]any{"club": club}); response != nil {
		return response
	}

	var request requests.UpdateClubRequest
	// Two results: validationErrors for bad input (422), err for anything
	// that stopped validation running at all, such as an unreadable body
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": err.Error()})
	}
	if validationErrors != nil {
		return validationFailed(ctx, validationErrors)
	}

	club.Name, club.Ground, club.FoundedYear = request.Name, request.Ground, request.FoundedYear
	if err := facades.Orm().Query().Save(club); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": club})
}

// Destroy deletes a club.
//
// Laravel:
//
//	Gate::authorize('delete', $club);
//	$club->delete();
//	return response()->noContent();
func (r *ClubController) Destroy(ctx http.Context) http.Response {
	club, response := findClub(ctx)
	if response != nil {
		return response
	}
	if _, response := authorise(ctx, "delete-club", map[string]any{"club": club}); response != nil {
		return response
	}

	if _, err := facades.Orm().Query().Delete(club); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().NoContent()
}

// findClub loads the club named by the {id} route parameter. It stands in for
// Laravel's route model binding, which Goravel does not have. It returns either
// the club or a ready-made error response, never both.
//
// Laravel's Club::findOrFail() throws ModelNotFoundException and the handler
// renders a 404. Here FindOrFail returns an error, and errors.Is checks which
// kind it is.
func findClub(ctx http.Context) (*models.Club, http.Response) {
	var club models.Club
	err := facades.Orm().Query().FindOrFail(&club, ctx.Request().RouteInt("id"))
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return nil, ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Club not found"})
	}
	if err != nil {
		return nil, serverError(ctx, err)
	}

	return &club, nil
}

// validationFailed returns a 422 in the same shape as Laravel's, with one
// difference: each field maps rule names to messages, not a list of messages.
func validationFailed(ctx http.Context, validationErrors validation.Errors) http.Response {
	return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
		"message": "The given data was invalid.",
		"errors":  validationErrors.All(),
	})
}

// serverError logs the real error and returns a generic 500, so database
// details never reach the client. Request() attaches the method, path and
// headers to the log entry, the context Laravel's exception handler adds.
func serverError(ctx http.Context, err error) http.Response {
	facades.Log().Request(ctx.Request()).Error(err)

	return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": "Server Error"})
}
