// Laravel equivalent: app/Http/Controllers/ClubController.php
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

// ClubController is a plain struct. Its methods become route handlers because
// they have the signature func(http.Context) http.Response.
type ClubController struct{}

// NewClubController is a constructor by convention only: Go has no
// __construct and nothing calls this for you. routes/api.go calls it once.
func NewClubController() *ClubController {
	return &ClubController{}
}

// Index lists every club.
func (r *ClubController) Index(ctx http.Context) http.Response {
	var clubs []models.Club
	if err := facades.Orm().Query().Order("name").Get(&clubs); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": clubs})
}

// Show returns one club, or 404.
func (r *ClubController) Show(ctx http.Context) http.Response {
	club, response := findClub(ctx)
	if response != nil {
		return response
	}

	return ctx.Response().Success().Json(http.Json{"data": club})
}

// Store creates a club and makes the logged-in user its admin.
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
