package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type UpdateClubRequest struct {
	Name        string `form:"name" json:"name"`
	Ground      string `form:"ground" json:"ground"`
	FoundedYear *int   `form:"founded_year" json:"founded_year"`
}

func (r *UpdateClubRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateClubRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{
		"name":   "trim",
		"ground": "trim",
	}
}

// Rules ignores the club being updated in the unique check, so saving a club
// without renaming it does not fail.
//
// Laravel: Rule::unique('clubs', 'name')->ignore($this->route('club'))
func (r *UpdateClubRequest) Rules(ctx http.Context) map[string]any {
	return clubRules(ctx.Request().RouteInt("id"))
}

func (r *UpdateClubRequest) Messages(ctx http.Context) map[string]string {
	return clubMessages()
}
