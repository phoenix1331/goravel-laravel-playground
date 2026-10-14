package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20261013183449CreatePlayersTable struct{}

// Signature The unique signature for the migration.
func (r *M20261013183449CreatePlayersTable) Signature() string {
	return "20261013183449_create_players_table"
}

// Up Run the migrations.
func (r *M20261013183449CreatePlayersTable) Up() error {
	if !facades.Schema().HasTable("players") {
		return facades.Schema().Create("players", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("club_id")
			table.String("name")
			table.String("position", 20)
			table.UnsignedTinyInteger("squad_number")
			table.TimestampsTz()

			// Deleting a club deletes its players
			table.Foreign("club_id").References("id").On("clubs").CascadeOnDelete()
			// Two players in one club cannot share a shirt number
			table.Unique("club_id", "squad_number")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261013183449CreatePlayersTable) Down() error {
	return facades.Schema().DropIfExists("players")
}
