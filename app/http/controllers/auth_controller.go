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

// Login exchanges an email and password for a JWT.
//
// Laravel (Sanctum):
//
//	$user = User::where('email', $request->email)->first();
//	if (! $user || ! Hash::check($request->password, $user->password)) {
//	    return response()->json(['message' => '...'], 401);
//	}
//	return ['access_token' => $user->createToken('api')->plainTextToken];
//
// Goravel: the default guard is a JWT. Login issues a signed token; nothing
// is stored. Auth always takes the request context: facades.Auth(ctx).
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

// Me returns the user the token belongs to. The Authenticate middleware has
// already parsed the token, so User can load them.
//
// Laravel: return response()->json(['data' => $request->user()]);
func (r *AuthController) Me(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "Unauthenticated."})
	}

	return ctx.Response().Success().Json(http.Json{"data": user})
}

// Logout disables the current token until it would have expired.
//
// Laravel: $request->user()->currentAccessToken()->delete();
//
// Goravel: a JWT cannot be deleted, so it goes on a block list in the cache.
// Tokens carry no unique ID: two logins by one user in the same second get
// the same token, and logging one out disables both.
func (r *AuthController) Logout(ctx http.Context) http.Response {
	if err := facades.Auth(ctx).Logout(); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().NoContent()
}
