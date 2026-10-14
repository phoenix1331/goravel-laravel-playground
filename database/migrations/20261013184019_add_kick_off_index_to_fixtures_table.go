// Laravel equivalent: database/migrations/xxxx_add_kick_off_index_to_fixtures_table.php
package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20261013184019AddKickOffIndexToFixturesTable struct{}

// Signature The unique signature for the migration.
func (r *M20261013184019AddKickOffIndexToFixturesTable) Signature() string {
	return "20261013184019_add_kick_off_index_to_fixtures_table"
}

// Up Run the migrations.
func (r *M20261013184019AddKickOffIndexToFixturesTable) Up() error {
	// The fixture list is ordered by kick-off, so index it
	return facades.Schema().Table("fixtures", func(table schema.Blueprint) {
		table.Index("kick_off_at")
	})
}

// Down Reverse the migrations.
func (r *M20261013184019AddKickOffIndexToFixturesTable) Down() error {
	return facades.Schema().Table("fixtures", func(table schema.Blueprint) {
		table.DropIndex("kick_off_at")
	})
}
