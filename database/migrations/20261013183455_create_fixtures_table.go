// Laravel equivalent: database/migrations/xxxx_create_fixtures_table.php - see docs/chapters/11-migrations.md
package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20261013183455CreateFixturesTable struct{}

// Signature The unique signature for the migration.
func (r *M20261013183455CreateFixturesTable) Signature() string {
	return "20261013183455_create_fixtures_table"
}

// #region create
// Up Run the migrations.
func (r *M20261013183455CreateFixturesTable) Up() error {
	if !facades.Schema().HasTable("fixtures") {
		return facades.Schema().Create("fixtures", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("home_club_id")
			table.UnsignedBigInteger("away_club_id")
			// Stored in UTC; see config/app.go "timezone"
			table.TimestampTz("kick_off_at")
			table.String("venue")
			table.String("status", 20).Default("scheduled")
			table.UnsignedTinyInteger("home_score").Nullable()
			table.UnsignedTinyInteger("away_score").Nullable()
			table.TimestampsTz()

			// Two foreign keys to the same table
			table.Foreign("home_club_id").References("id").On("clubs").CascadeOnDelete()
			table.Foreign("away_club_id").References("id").On("clubs").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261013183455CreateFixturesTable) Down() error {
	return facades.Schema().DropIfExists("fixtures")
}

// #endregion create
