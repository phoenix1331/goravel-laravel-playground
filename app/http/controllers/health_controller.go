// Laravel equivalent: the /up health route in bootstrap/app.php - see docs/chapters/03-routing.md
package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// #region show
// Show reports whether the app can reach its database.
func (r *HealthController) Show(ctx http.Context) http.Response {
	// No exceptions in Go: Orm().DB() hands back an error value to check
	db, err := facades.Orm().DB()
	if err == nil {
		err = db.PingContext(ctx.Context())
	}
	if err != nil {
		return ctx.Response().Json(http.StatusServiceUnavailable, http.Json{
			"status":   "error",
			"database": "down",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"status":   "ok",
		"database": "up",
	})
}

// #endregion show
