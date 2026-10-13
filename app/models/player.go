// Laravel equivalent: app/Models/Player.php - see docs/chapters/09-models.md
package models

import (
	"github.com/goravel/framework/database/orm"
)

// #region position
// Position is a named string type: Go's closest thing to a backed enum.
type Position string

const (
	Goalkeeper Position = "goalkeeper"
	Defender   Position = "defender"
	Midfielder Position = "midfielder"
	Forward    Position = "forward"
)

// Positions lists every valid position, for validation rules and seeders.
var Positions = []Position{Goalkeeper, Defender, Midfielder, Forward}

// #endregion position

// #region player
type Player struct {
	orm.Model
	ClubID      uint     `json:"club_id"`
	Name        string   `json:"name"`
	Position    Position `json:"position"`
	SquadNumber int      `json:"squad_number"`
}

// #endregion player
