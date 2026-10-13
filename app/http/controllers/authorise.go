// Laravel equivalent: $this->authorize() / Gate::authorize() - see docs/chapters/20-authorisation.md
package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
)

// #region authorise
// authorise loads the logged-in user and asks the gate whether they may
// perform ability. It returns the user, or a ready-made 401/403 response.
// The user is added to arguments, because Goravel's gate does not know who
// is asking: Laravel passes the user to every gate for you.
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

// #endregion authorise
