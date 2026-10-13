// Laravel equivalent: app/Http/Controllers/ClubController.php - see docs/chapters/04-controllers.md
package controllers

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	frameworkerrors "github.com/goravel/framework/errors"

	"goravel/app/facades"
	"goravel/app/models"
)

// #region struct
// ClubController is a plain struct. Its methods become route handlers because
// they have the signature func(http.Context) http.Response.
type ClubController struct{}

// NewClubController is a constructor by convention only: Go has no
// __construct and nothing calls this for you. routes/api.go calls it once.
func NewClubController() *ClubController {
	return &ClubController{}
}

// #endregion struct

// #region index
// Index lists every club.
func (r *ClubController) Index(ctx http.Context) http.Response {
	var clubs []models.Club
	if err := facades.Orm().Query().Order("name").Get(&clubs); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": clubs})
}

// #endregion index

// #region show
// Show returns one club, or 404.
func (r *ClubController) Show(ctx http.Context) http.Response {
	club, response := findClub(ctx)
	if response != nil {
		return response
	}

	return ctx.Response().Success().Json(http.Json{"data": club})
}

// #endregion show

// #region store
// clubInput is the JSON body accepted by Store and Update.
type clubInput struct {
	Name        string `json:"name"`
	Ground      string `json:"ground"`
	FoundedYear *int   `json:"founded_year"`
}

// Store creates a club from the JSON body.
func (r *ClubController) Store(ctx http.Context) http.Response {
	var input clubInput
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid JSON body"})
	}

	club := models.Club{Name: input.Name, Ground: input.Ground, FoundedYear: input.FoundedYear}
	if err := facades.Orm().Query().Create(&club); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{"data": club})
}

// #endregion store

// Update replaces a club's fields from the JSON body.
func (r *ClubController) Update(ctx http.Context) http.Response {
	club, response := findClub(ctx)
	if response != nil {
		return response
	}

	var input clubInput
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "Invalid JSON body"})
	}

	club.Name, club.Ground, club.FoundedYear = input.Name, input.Ground, input.FoundedYear
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

	if _, err := facades.Orm().Query().Delete(club); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().NoContent()
}

// #region find
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

// #endregion find

// serverError logs the real error and returns a generic 500, so database
// details never reach the client.
func serverError(ctx http.Context, err error) http.Response {
	facades.Log().Error(err)

	return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": "Server Error"})
}
