package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

type M20261013185753CreateResultNotificationsTable struct{}

// Signature The unique signature for the migration.
func (r *M20261013185753CreateResultNotificationsTable) Signature() string {
	return "20261013185753_create_result_notifications_table"
}

// Up Run the migrations.
func (r *M20261013185753CreateResultNotificationsTable) Up() error {
	if !facades.Schema().HasTable("result_notifications") {
		return facades.Schema().Create("result_notifications", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id")
			table.UnsignedBigInteger("fixture_id")
			table.String("message")
			table.TimestampsTz()

			table.Foreign("user_id").References("id").On("users").CascadeOnDelete()
			table.Foreign("fixture_id").References("id").On("fixtures").CascadeOnDelete()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20261013185753CreateResultNotificationsTable) Down() error {
	return facades.Schema().DropIfExists("result_notifications")
}
