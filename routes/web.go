// Laravel equivalent: routes/web.php - see docs/chapters/03-routing.md
package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support"

	"goravel/app/facades"
	"goravel/app/http/controllers"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})

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
