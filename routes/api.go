// Laravel equivalent: routes/api.php - see docs/chapters/03-routing.md
package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"

	"goravel/app/facades"
	"goravel/app/http/controllers"
	"goravel/app/http/middleware"
)

// #region api
func Api() {
	// Every route in the group runs through EnsureJson first, like
	// Route::middleware(EnsureJsonRequest::class)->prefix('api')->group(...)
	facades.Route().Prefix("api").Middleware(middleware.EnsureJson()).Group(func(router route.Router) {
		// A closure, exactly like Route::get('/league', fn () => ...)
		router.Get("league", func(ctx http.Context) http.Response {
			return ctx.Response().Success().Json(http.Json{
				"name": facades.Config().GetString("club.league_name"),
			})
		})

		// A named function works anywhere a closure does: a handler is just
		// any func(http.Context) http.Response
		router.Get("league/points", leaguePoints)

		leagueController := controllers.NewLeagueController()
		router.Get("league/table", leagueController.Table)
		router.Get("league/goals", leagueController.Goals)

		// #region resource
		// One line registers index, store, show, update and destroy,
		// like Route::apiResource('clubs', ClubController::class)
		router.Resource("clubs", controllers.NewClubController())
		// #endregion resource

		// #region nested
		playerController := controllers.NewPlayerController()
		router.Get("clubs/{id}/players", playerController.Index)
		router.Post("clubs/{id}/players", playerController.Store)

		fixtureController := controllers.NewFixtureController()
		router.Get("fixtures", fixtureController.Index)
		router.Post("fixtures/{id}/result", fixtureController.RecordResult)
		// #endregion nested
	})
}

// #endregion api

// #region handler
func leaguePoints(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"win":  facades.Config().GetInt("club.points.win"),
		"draw": facades.Config().GetInt("club.points.draw"),
		"loss": facades.Config().GetInt("club.points.loss"),
	})
}

// #endregion handler
