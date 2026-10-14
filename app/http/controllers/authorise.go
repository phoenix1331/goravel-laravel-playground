package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

// authorise loads the logged-in user and asks the gate whether they may
// perform ability. It returns the user, or a ready-made 401/403 response.
//
// Laravel:
//
//	Gate::authorize('update', $club);  // throws AuthorizationException: 403
//
// Goravel: the gate is not told who is asking, so the user is added to the
// arguments here. Inspect returns a response instead of throwing.
func authorise(ctx http.Context, ability string, arguments map[string]any) (*models.User, http.Response) {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return nil, ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "Unauthenticated."})
	}

	arguments["user"] = &user
	if response := facades.Gate().WithContext(ctx).Inspect(ability, arguments); !response.Allowed() {
		return nil, ctx.Response().Json(http.StatusForbidden, http.Json{"message": response.Message()})
	}

	return &user, nil
}
