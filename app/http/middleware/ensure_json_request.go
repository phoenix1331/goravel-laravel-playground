// Package middleware holds HTTP middleware.
package middleware

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
)

// EnsureJsonRequest rejects API writes whose body is not JSON, with 415.
//
// Laravel:
//
//	public function handle(Request $request, Closure $next): Response
//	{
//	    if (in_array($request->method(), ['POST', 'PUT', 'PATCH']) && ! $request->isJson()) {
//	        return response()->json([...], 415);  // returning stops the chain
//	    }
//	    return $next($request);
//	}
//
// Goravel: Handle returns nothing. Call Next() to continue, or write a
// response, call Abort() and return. Forgetting that return falls through
// to Next(), the most common Goravel middleware bug.
type EnsureJsonRequest struct{}

// EnsureJson is the constructor used in route definitions.
func EnsureJson() http.Middleware {
	return &EnsureJsonRequest{}
}

// Signature names the middleware, for logs and exclusion lists.
func (m *EnsureJsonRequest) Signature() string {
	return "ensure_json"
}

// Handle runs before the route handler. It has no return value: it either
// calls Next() to continue, or writes a response and calls Abort() to stop.
func (m *EnsureJsonRequest) Handle(ctx http.Context) {
	request := ctx.Request()

	switch request.Method() {
	case "POST", "PUT", "PATCH":
		contentType := request.Header("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			ctx.Response().Json(http.StatusUnsupportedMediaType, http.Json{
				"message": "Send the request body as JSON (Content-Type: application/json).",
			}).Abort()

			return
		}
	}

	request.Next()
}
