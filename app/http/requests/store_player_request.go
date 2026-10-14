package requests

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"goravel/app/models"
)

// StorePlayerRequest validates a new player.
//
// Laravel:
//
//	'position' => ['required', Rule::enum(Position::class)],
//
// Goravel: no enum rule, so the allowed values are joined into an in: rule
// from models.Positions.
type StorePlayerRequest struct {
	Name        string          `form:"name" json:"name"`
	Position    models.Position `form:"position" json:"position"`
	SquadNumber int             `form:"squad_number" json:"squad_number"`
}

func (r *StorePlayerRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *StorePlayerRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{
		"name":     "trim",
		"position": "trim|lower",
	}
}

func (r *StorePlayerRequest) Rules(ctx http.Context) map[string]any {
	positions := make([]string, len(models.Positions))
	for i, position := range models.Positions {
		positions[i] = string(position)
	}

	return map[string]any{
		"name":         "required|string|max:100",
		"position":     "required|in:" + strings.Join(positions, ","),
		"squad_number": "required|integer|min:1|max:99",
	}
}
