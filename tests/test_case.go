// Laravel equivalent: tests/TestCase.php - see docs/chapters/23-testing.md
package tests

import (
	"fmt"
	"os"
	gotesting "testing"

	"github.com/goravel/framework/testing"

	// Registers the "memdb" SQLite VFS used by .env.testing. A plain :memory:
	// DSN gives every pooled connection its own empty database, so migrations
	// run on one connection and queries hit another. memdb shares one
	// in-memory database across all connections in the test process.
	_ "github.com/ncruces/go-sqlite3/vfs/memdb"

	"goravel/app/facades"
	"goravel/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}

// #region main
// Main migrates the database once, then runs the package's tests. Call it
// from TestMain in every test package:
//
//	func TestMain(m *testing.M) { os.Exit(tests.Main(m)) }
//
// go test runs each package as its own process, so each package gets its own
// in-memory database with no cross-package interference.
func Main(m *gotesting.M) int {
	if err := facades.Artisan().Call("--no-ansi migrate"); err != nil {
		fmt.Fprintln(os.Stderr, "migrate failed:", err)
		return 1
	}

	return m.Run()
}

// #endregion main

// #region truncate
// TruncateTables deletes every row from every table except the migrations
// table, so each test starts empty without paying for a full re-migrate.
// Call it from SetupTest.
func (r *TestCase) TruncateTables() error {
	tables, err := facades.Schema().GetTables()
	if err != nil {
		return err
	}

	migrations := facades.Config().GetString("database.migrations.table", "migrations")
	for _, table := range tables {
		if table.Name == migrations {
			continue
		}
		// The query builder refuses an unconditional DELETE, so match every row explicitly
		if _, err := facades.DB().Table(table.Name).Where("1 = 1").Delete(); err != nil {
			return fmt.Errorf("truncate %s: %w", table.Name, err)
		}
	}

	return nil
}

// #endregion truncate
