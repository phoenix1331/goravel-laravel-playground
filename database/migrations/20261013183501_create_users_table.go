package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20261013183501CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20261013183501CreateUsersTable) Signature() string {
	return "20261013183501_create_users_table"
}

// Up Run the migrations.
func (r *M20261013183501CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		return facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID()
			// The club this user administers; null for users without a club
			table.UnsignedBigInteger("club_id").Nullable()
			table.String("name")
			table.String("email")
			table.String("password")
			table.TimestampsTz()

			table.Unique("email")
			table.Foreign("club_id").References("id").On("clubs").NullOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261013183501CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
