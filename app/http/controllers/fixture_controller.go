// Laravel equivalent: app/Http/Controllers/FixtureController.php - see docs/chapters/10-relationships.md
package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

type FixtureController struct{}

func NewFixtureController() *FixtureController {
	return &FixtureController{}
}

// #region index
// Index lists fixtures with both clubs eager loaded, like
// Fixture::with(['homeClub', 'awayClub'])->orderBy('kick_off_at')->get().
func (r *FixtureController) Index(ctx http.Context) http.Response {
	var fixtures []models.Fixture
	err := facades.Orm().Query().
		With("HomeClub").
		With("AwayClub").
		Order("kick_off_at").
		Get(&fixtures)
	if err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": fixtures})
}

// #endregion index
