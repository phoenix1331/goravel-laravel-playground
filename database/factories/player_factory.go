// Laravel equivalent: database/factories/PlayerFactory.php - see docs/chapters/12-factories-seeders.md
package factories

import (
	"github.com/brianvoe/gofakeit/v7"
)

// PlayerFactory has no ClubID default: the caller always says which club,
// e.g. Factory().Create(&player, map[string]any{"ClubID": club.ID}).
// Position is the string value; the ORM decoder converts it to models.Position.
type PlayerFactory struct{}

func (f *PlayerFactory) Definition() map[string]any {
	return map[string]any{
		"Name":        gofakeit.FirstName() + " " + gofakeit.LastName(),
		"Position":    gofakeit.RandomString([]string{"goalkeeper", "defender", "midfielder", "forward"}),
		"SquadNumber": gofakeit.IntRange(1, 99),
	}
}
