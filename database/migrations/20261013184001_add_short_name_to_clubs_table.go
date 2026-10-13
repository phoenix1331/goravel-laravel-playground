// Laravel equivalent: database/migrations/xxxx_add_short_name_to_clubs_table.php - see docs/chapters/11-migrations.md
package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20261013184001AddShortNameToClubsTable struct{}

// Signature The unique signature for the migration.
func (r *M20261013184001AddShortNameToClubsTable) Signature() string {
	return "20261013184001_add_short_name_to_clubs_table"
}

// #region alter
// Up Run the migrations.
func (r *M20261013184001AddShortNameToClubsTable) Up() error {
	// Table() alters an existing table, like Schema::table()
	return facades.Schema().Table("clubs", func(table schema.Blueprint) {
		// Short name for tight spaces such as the league table, e.g. "ASH"
		table.String("short_name", 10).Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20261013184001AddShortNameToClubsTable) Down() error {
	return facades.Schema().Table("clubs", func(table schema.Blueprint) {
		table.DropColumn("short_name")
	})
}

// #endregion alter
