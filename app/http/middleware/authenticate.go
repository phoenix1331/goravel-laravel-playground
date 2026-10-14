// Laravel equivalent: the auth:sanctum / auth:api middleware
package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

// Authenticate requires a valid bearer token. It parses the token into the
// request context, so later code can call facades.Auth(ctx).User(&user).
type AuthenticateRequest struct{}

func Authenticate() http.Middleware {
	return &AuthenticateRequest{}
}

func (m *AuthenticateRequest) Signature() string {
	return "auth"
}

func (m *AuthenticateRequest) Handle(ctx http.Context) {
	token := ctx.Request().Header("Authorization")
	if token == "" {
		unauthenticated(ctx)
		return
	}

	// Parse checks the signature, expiry and logout list. Any error (expired,
	// invalid, logged out) means the request is not authenticated.
	if _, err := facades.Auth(ctx).Parse(token); err != nil {
		unauthenticated(ctx)
		return
	}

	ctx.Request().Next()
}

func unauthenticated(ctx http.Context) {
	ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "Unauthenticated."}).Abort()
}
