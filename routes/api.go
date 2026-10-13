// Laravel equivalent: routes/api.php - see docs/chapters/03-routing.md
package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"

	"goravel/app/facades"
	"goravel/app/http/controllers"
	"goravel/app/http/middleware"
	"goravel/app/services"
)

// #region api
func Api() {
	// Every route in the group runs through EnsureJson first, like
	// Route::middleware(EnsureJsonRequest::class)->prefix('api')->group(...)
	facades.Route().Prefix("api").Middleware(middleware.EnsureJson()).Group(func(router route.Router) {
		leagueRoutes(router)
		publicRoutes(router)
		authenticatedRoutes(router)
	})
}

// #endregion api

// #region league
func leagueRoutes(router route.Router) {
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
}

// #endregion league

// #region handler
func leaguePoints(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"win":  facades.Config().GetInt("club.points.win"),
		"draw": facades.Config().GetInt("club.points.draw"),
		"loss": facades.Config().GetInt("club.points.loss"),
	})
}

// #endregion handler

// #region public
// publicRoutes can be read by anyone. Logging in is public too.
func publicRoutes(router route.Router) {
	clubController := controllers.NewClubController()
	router.Get("clubs", clubController.Index)
	router.Get("clubs/{id}", clubController.Show)
	router.Get("clubs/{id}/players", controllers.NewPlayerController().Index)
	router.Get("fixtures", controllers.NewFixtureController().Index)

	router.Post("auth/login", controllers.NewAuthController().Login)

	// A stand-in for the league's results API, so FIXTURE_FEED_DRIVER=http
	// can be tried locally. It serves the same JSON as the file feed.
	router.Get("demo/results", func(ctx http.Context) http.Response {
		results, err := services.FileFeed{Path: facades.Config().GetString("club.feed.path")}.Results()
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{"message": "Feed unavailable"})
		}

		return ctx.Response().Success().Json(results)
	})
}

// #endregion public

// #region authenticated
// authenticatedRoutes need a bearer token, like a
// Route::middleware('auth:sanctum')->group(...) block.
func authenticatedRoutes(router route.Router) {
	router.Middleware(middleware.Authenticate()).Group(func(router route.Router) {
		authController := controllers.NewAuthController()
		router.Get("auth/me", authController.Me)
		router.Post("auth/logout", authController.Logout)

		clubController := controllers.NewClubController()
		router.Post("clubs", clubController.Store)
		router.Put("clubs/{id}", clubController.Update)
		router.Patch("clubs/{id}", clubController.Update)
		router.Delete("clubs/{id}", clubController.Destroy)

		router.Post("clubs/{id}/players", controllers.NewPlayerController().Store)
		router.Post("fixtures/{id}/result", controllers.NewFixtureController().RecordResult)
	})
}

// #endregion authenticated
