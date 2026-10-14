// Package requests holds form requests: validation rules plus the struct
// the validated input is bound into.
package requests

import (
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/http"
)

// StoreClubRequest validates a new club.
//
// Laravel:
//
//	class StoreClubRequest extends FormRequest
//	{
//	    public function rules(): array
//	    {
//	        return [
//	            'name' => ['required', 'string', 'max:100', Rule::unique('clubs', 'name')],
//	            'ground' => ['required', 'string', 'max:100'],
//	            'founded_year' => ['nullable', 'integer', 'min:1857', 'max:'.now()->year],
//	        ];
//	    }
//	}
//
// Goravel: the same rule strings, but the request is a struct, and it is
// also where the validated data lands: request.FoundedYear is a typed *int,
// the equivalent of $request->validated(). Validation does not run before
// the controller; the controller calls ctx.Request().ValidateRequest.
type StoreClubRequest struct {
	Name        string `form:"name" json:"name"`
	Ground      string `form:"ground" json:"ground"`
	FoundedYear *int   `form:"founded_year" json:"founded_year"`
}

// Authorize returns nil to allow the request, or an error to reject it.
// The real check is in the controller (see controllers/authorise.go):
// ValidateRequest reports an Authorize error as a plain error, not a 403.
func (r *StoreClubRequest) Authorize(ctx http.Context) error {
	return nil
}

// Filters clean the input before the rules run. Laravel's TrimStrings
// middleware trims every input globally; Goravel does not, so "   " would
// otherwise pass "required".
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

// clubRules is shared by the store and update requests. Rules are built in
// code, so values such as the current year are just function calls.
// exceptID is the club being updated, or 0 when creating.
//
// Watch the unique rule's argument order when copying from Laravel:
//
//	Laravel: unique:clubs,name,{ignoredValue},{idColumn}
//	Goravel: unique:clubs,name,{idColumn},{ignoredValue}
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

func clubMessages() map[string]string {
	return map[string]string{
		"name.unique":      "A club with this name already exists.",
		"founded_year.min": "The oldest football club was founded in 1857.",
	}
}
