// Laravel equivalent: app/Http/Requests/StoreClubRequest.php - see docs/chapters/07-validation.md
package requests

import (
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/http"
)

// #region request
// StoreClubRequest is both the rules and the destination for the validated
// data: after validation passes, Goravel binds the input into these fields.
type StoreClubRequest struct {
	Name        string `form:"name" json:"name"`
	Ground      string `form:"ground" json:"ground"`
	FoundedYear *int   `form:"founded_year" json:"founded_year"`
}

// Authorize returns nil to allow the request, or an error to reject it.
// Chapter 20 adds the real check.
func (r *StoreClubRequest) Authorize(ctx http.Context) error {
	return nil
}

// Filters clean the input before the rules run.
func (r *StoreClubRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{
		"name":   "trim",
		"ground": "trim",
	}
}

func (r *StoreClubRequest) Rules(ctx http.Context) map[string]any {
	return clubRules(0)
}

func (r *StoreClubRequest) Messages(ctx http.Context) map[string]string {
	return clubMessages()
}

// #endregion request

// #region rules
// clubRules is shared by the store and update requests. Rules are built in
// code, so values such as the current year are just function calls.
// exceptID is the club being updated, or 0 when creating.
func clubRules(exceptID int) map[string]any {
	unique := "unique:clubs,name"
	if exceptID > 0 {
		// Goravel's argument order is table,column,idColumn,ignoredValue.
		// Laravel's is table,column,ignoredValue,idColumn.
		unique = fmt.Sprintf("unique:clubs,name,id,%d", exceptID)
	}

	return map[string]any{
		"name":         "required|string|max:100|" + unique,
		"ground":       "required|string|max:100",
		"founded_year": fmt.Sprintf("nullable|integer|min:1857|max:%d", time.Now().Year()),
	}
}

// #endregion rules

func clubMessages() map[string]string {
	return map[string]string{
		"name.unique":      "A club with this name already exists.",
		"founded_year.min": "The oldest football club was founded in 1857.",
	}
}
