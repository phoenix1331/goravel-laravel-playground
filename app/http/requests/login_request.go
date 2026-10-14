// Laravel equivalent: app/Http/Requests/LoginRequest.php
package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type LoginRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (r *LoginRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *LoginRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{"email": "trim|lower"}
}

func (r *LoginRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"email":    "required|email",
		"password": "required|string",
	}
}
