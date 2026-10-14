package routes

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/http/controllers"
)

// Web registers the routes outside /api.
//
// Laravel: routes/web.php, plus the health route and 404 handling you get
// for free from bootstrap/app.php and the exception handler.
func Web() {
	// The app is API only, so / lists the endpoints instead of rendering a page
	facades.Route().Get("/", controllers.NewIndexController().Index)

	facades.Route().Static("public", "./public")

	healthController := controllers.NewHealthController()
	facades.Route().Get("/health", healthController.Show)

	// Laravel:
	//
	//	Route::fallback(fn () => response()->json(['message' => 'Not Found'], 404));
	//
	// Optional in Laravel, needed here: without a fallback, the Gin driver
	// answers unknown URLs with an empty 200.
	facades.Route().Fallback(func(ctx http.Context) http.Response {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Not Found",
		})
	})
}
