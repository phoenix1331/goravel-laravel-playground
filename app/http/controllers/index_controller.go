// Laravel equivalent: a closure on Route::get('/') returning JSON - see docs/chapters/03-routing.md
package controllers

import (
	"sort"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
)

type IndexController struct{}

func NewIndexController() *IndexController {
	return &IndexController{}
}

type endpoint struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// #region index
// Index describes the API and lists every endpoint. The list is read from the
// router at request time, so it can never fall out of date.
func (r *IndexController) Index(ctx http.Context) http.Response {
	var endpoints []endpoint
	for _, route := range facades.Route().GetRoutes() {
		switch route.Method {
		case "STATIC":
			continue
		case "RESOURCE":
			// A resource route is registered once but serves five endpoints
			endpoints = append(endpoints,
				endpoint{"GET", route.Path},
				endpoint{"POST", route.Path},
				endpoint{"GET", route.Path + "/{id}"},
				endpoint{"PUT|PATCH", route.Path + "/{id}"},
				endpoint{"DELETE", route.Path + "/{id}"},
			)
		default:
			// GET routes also answer HEAD; leave that detail out
			endpoints = append(endpoints, endpoint{strings.TrimSuffix(route.Method, "|HEAD"), route.Path})
		}
	}

	sort.SliceStable(endpoints, func(i, j int) bool {
		return endpoints[i].Path < endpoints[j].Path
	})

	return ctx.Response().Success().Json(http.Json{
		"name":        facades.Config().GetString("app.name"),
		"description": "API only: every endpoint returns JSON. See request.http for example requests.",
		"docs":        "https://github.com/phoenix1331/goravel-for-laravel-devs",
		"endpoints":   endpoints,
	})
}

// #endregion index
