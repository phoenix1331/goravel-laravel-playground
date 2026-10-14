// Package tests is the shared test harness, the equivalent of Laravel's
// tests/TestCase.php plus the testing settings in phpunit.xml.
//
// Laravel: class TestCase extends BaseTestCase, with RefreshDatabase and
// actingAs(). Goravel: suites embed this TestCase; TestMain migrates once,
// TruncateTables resets between tests, and LoginAs stands in for actingAs.
package tests

import (
	"fmt"
	"os"
	"strings"
	gotesting "testing"

	contractshttp "github.com/goravel/framework/contracts/testing/http"
	"github.com/goravel/framework/testing"

	// Registers the "memdb" SQLite VFS used by .env.testing. A plain :memory:
	// DSN gives every pooled connection its own empty database, so migrations
	// run on one connection and queries hit another. memdb shares one
	// in-memory database across all connections in the test process.
	_ "github.com/ncruces/go-sqlite3/vfs/memdb"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/bootstrap"
	"goravel/database/factories"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}

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

// CreateAdmin makes a user who administers the given club.
func (r *TestCase) CreateAdmin(t *gotesting.T, clubID uint) models.User {
	t.Helper()

	var admin models.User
	if err := facades.Orm().Factory().Create(&admin, map[string]any{"ClubID": &clubID}); err != nil {
		t.Fatalf("create admin: %v", err)
	}

	return admin
}

// LoginAs logs a user in through the real login endpoint and returns a
// request builder that sends their token and a JSON content type, the
// equivalent of $this->actingAs($user). The user must have been made by the
// factory, so their password is factories.Password.
func (r *TestCase) LoginAs(t *gotesting.T, user models.User) contractshttp.Request {
	t.Helper()

	body := fmt.Sprintf(`{"email":%q,"password":%q}`, user.Email, factories.Password)
	response, err := r.Http(t).WithHeader("Content-Type", "application/json").
		Post("/api/auth/login", strings.NewReader(body))
	if err != nil {
		t.Fatalf("login as %s: %v", user.Email, err)
	}

	var login struct {
		AccessToken string `json:"access_token"`
	}
	if err := response.Bind(&login); err != nil || login.AccessToken == "" {
		t.Fatalf("login as %s: no token in response", user.Email)
	}

	return r.Http(t).WithToken(login.AccessToken).WithHeader("Content-Type", "application/json")
}
