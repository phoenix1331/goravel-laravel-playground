// Laravel equivalent: app/Http/Requests/RecordResultRequest.php - see docs/chapters/15-events.md
package requests

import (
	"github.com/goravel/framework/contracts/http"
)

type RecordResultRequest struct {
	HomeScore int `form:"home_score" json:"home_score"`
	AwayScore int `form:"away_score" json:"away_score"`
}

func (r *RecordResultRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *RecordResultRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"home_score": "required|integer|min:0|max:99",
		"away_score": "required|integer|min:0|max:99",
	}
}
