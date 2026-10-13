// Laravel equivalent: app/Http/Controllers/AuthController.php - see docs/chapters/19-authentication.md
package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/http/requests"
	"goravel/app/models"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// #region login
// Login exchanges an email and password for a JWT.
func (r *AuthController) Login(ctx http.Context) http.Response {
	var request requests.LoginRequest
	validationErrors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": err.Error()})
	}
	if validationErrors != nil {
		return validationFailed(ctx, validationErrors)
	}

	// Look the user up and check the password ourselves: there is no
	// Auth::attempt(). The same message for both failures avoids telling an
	// attacker which emails exist.
	var user models.User
	if err := facades.Orm().Query().Where("email", request.Email).First(&user); err != nil {
		return serverError(ctx, err)
	}
	if user.ID == 0 || !facades.Hash().Check(request.Password, user.Password) {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "These credentials do not match our records."})
	}

	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{
		"token_type":   "Bearer",
		"access_token": token,
		"expires_in":   facades.Config().GetInt("jwt.ttl") * 60,
	})
}

// #endregion login

// #region me
// Me returns the user the token belongs to. The Authenticate middleware has
// already parsed the token, so User can load them.
func (r *AuthController) Me(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "Unauthenticated."})
	}

	return ctx.Response().Success().Json(http.Json{"data": user})
}

// #endregion me

// Logout disables the current token until it would have expired.
func (r *AuthController) Logout(ctx http.Context) http.Response {
	if err := facades.Auth(ctx).Logout(); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().NoContent()
}
