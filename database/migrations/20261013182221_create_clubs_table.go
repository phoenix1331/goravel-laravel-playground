// Laravel equivalent: database/migrations/xxxx_create_clubs_table.php - see docs/chapters/11-migrations.md
package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20261013182221CreateClubsTable struct{}

// Signature The unique signature for the migration.
func (r *M20261013182221CreateClubsTable) Signature() string {
	return "20261013182221_create_clubs_table"
}

// Up Run the migrations.
func (r *M20261013182221CreateClubsTable) Up() error {
	if !facades.Schema().HasTable("clubs") {
		return facades.Schema().Create("clubs", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("ground")
			table.UnsignedSmallInteger("founded_year").Nullable()
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261013182221CreateClubsTable) Down() error {
	return facades.Schema().DropIfExists("clubs")
}
