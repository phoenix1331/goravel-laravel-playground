// Laravel equivalent: routes/web.php - see docs/chapters/03-routing.md
package routes

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/http/controllers"
)

func Web() {
	// The app is API only, so / lists the endpoints instead of rendering a page
	facades.Route().Get("/", controllers.NewIndexController().Index)

	facades.Route().Static("public", "./public")

	// #region health
	healthController := controllers.NewHealthController()
	facades.Route().Get("/health", healthController.Show)
	// #endregion health

	// #region fallback
	// Without a fallback, the Gin driver answers unknown URLs with an empty
	// 200: global middleware runs and nothing sets a status. Laravel returns
	// 404 by default, so match that.
	facades.Route().Fallback(func(ctx http.Context) http.Response {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Not Found",
		})
	})
	// #endregion fallback
}
