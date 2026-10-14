package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// Show reports whether the app can reach its database.
//
// Laravel:
//
//	try {
//	    DB::connection()->getPdo();
//	} catch (Throwable) {
//	    return response()->json(['status' => 'error'], 503);
//	}
//
// Goravel: no try/catch. Each call that can fail returns an error to check.
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
