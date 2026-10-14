// Package migrations holds the schema migrations, registered in
// bootstrap/migrations.go.
package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/app/facades"
)

// M20261013182221CreateClubsTable creates the clubs table.
//
// Laravel:
//
//	return new class extends Migration
//	{
//	    public function up(): void
//	    {
//	        Schema::create('clubs', function (Blueprint $table) {
//	            $table->id();
//	            $table->string('name');
//	            // ...
//	        });
//	    }
//	};
//
// Goravel: a named struct, not an anonymous class, so the timestamp is in
// the type name to keep it unique. Laravel finds migrations by scanning the
// folder; Goravel lists them in bootstrap/migrations.go, which
// make:migration updates for you. Up and Down return errors rather than
// throwing.
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
