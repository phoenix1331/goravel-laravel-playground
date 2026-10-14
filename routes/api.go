// Package routes registers every HTTP route.
//
// Laravel: routes/api.php, loaded because bootstrap/app.php names it.
// Goravel: plain Go functions, called from WithRouting in bootstrap/app.go.
// Forget to call one and its routes do not exist.
package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"

	"goravel/app/facades"
	"goravel/app/http/controllers"
	"goravel/app/http/middleware"
	"goravel/app/services"
)

// Api registers the JSON API.
//
// Laravel:
//
//	Route::prefix('api')->middleware(EnsureJsonRequest::class)->group(function () {
//	    // ...
//	});
//
// Goravel: the same shape. Group takes a function that receives the group's
// router, and routes are registered on that router, not on the facade.
func Api() {
	facades.Route().Prefix("api").Middleware(middleware.EnsureJson()).Group(func(router route.Router) {
		leagueRoutes(router)
		publicRoutes(router)
		authenticatedRoutes(router)
	})
}

// leagueRoutes shows the two kinds of handler.
//
// Laravel:
//
//	Route::get('/league', fn () => ['name' => config('club.league_name')]);
//	Route::get('/league/points', LeaguePointsController::class);
//
// Goravel: a handler is any func(http.Context) http.Response, so a closure,
// a plain named function and a controller method are all accepted. Laravel
// wants a controller class for anything that is not a closure.
func leagueRoutes(router route.Router) {
	router.Get("league", func(ctx http.Context) http.Response {
		return ctx.Response().Success().Json(http.Json{
			"name": facades.Config().GetString("club.league_name"),
		})
	})

	router.Get("league/points", leaguePoints)

	leagueController := controllers.NewLeagueController()
	router.Get("league/table", leagueController.Table)
	router.Get("league/goals", leagueController.Goals)
}

func leaguePoints(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"win":  facades.Config().GetInt("club.points.win"),
		"draw": facades.Config().GetInt("club.points.draw"),
		"loss": facades.Config().GetInt("club.points.loss"),
	})
}

// publicRoutes can be read by anyone. Logging in is public too.
//
// Laravel:
//
//	Route::get('/clubs', [ClubController::class, 'index']);
//	Route::get('/clubs/{club}', [ClubController::class, 'show']);
//	Route::get('/clubs/{club}/players', [ClubPlayerController::class, 'index']);
//
// Goravel: controllers are values you construct, and you pass the method
// itself (clubController.Index), not a class name and a method string. The
// parameter is a plain {id}: there is no route model binding.
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

// authenticatedRoutes need a bearer token.
//
// Laravel:
//
//	Route::middleware('auth:sanctum')->group(function () {
//	    Route::apiResource('clubs', ClubController::class)->except(['index', 'show']);
//	    // ...
//	});
//
// Goravel: middleware is a value (middleware.Authenticate()), not an alias
// string. router.Resource("clubs", controller) exists, but it registers all
// five routes, so the writes are listed one by one to keep the reads public.
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
