package models

import (
	"github.com/goravel/framework/contracts/database/factory"
	"github.com/goravel/framework/database/orm"

	"goravel/database/factories"
)

// Position is a named string type: Go's closest thing to a backed enum.
//
// Laravel:
//
//	enum Position: string
//	{
//	    case Goalkeeper = 'goalkeeper';
//	    // ...
//	}
//
// Weaker than a PHP enum: Position("striker") compiles, so validate input
// against Positions.
type Position string

const (
	Goalkeeper Position = "goalkeeper"
	Defender   Position = "defender"
	Midfielder Position = "midfielder"
	Forward    Position = "forward"
)

// Positions lists every valid position, for validation rules and seeders.
var Positions = []Position{Goalkeeper, Defender, Midfielder, Forward}

// Player is a member of a club's squad.
//
// Laravel:
//
//	public function club(): BelongsTo { return $this->belongsTo(Club::class); }
//
// Goravel: the relation is the Club field below, paired with ClubID by name.
type Player struct {
	orm.Model
	ClubID      uint     `json:"club_id"`
	Name        string   `json:"name"`
	Position    Position `json:"position"`
	SquadNumber int      `json:"squad_number"`

	// Belongs to: a pointer to the parent, nil until loaded. ClubID above is
	// the foreign key, matched by name (Club + ID).
	Club *Club `json:"club,omitempty"`
}

// Factory links the model to its factory, like the HasFactory trait.
func (p *Player) Factory() factory.Factory {
	return &factories.PlayerFactory{}
}
