// Laravel equivalent: bootstrap/app.php - see docs/chapters/01-installation.md
package bootstrap

import (
	"github.com/goravel/framework/contracts/event"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"

	"goravel/app/events"
	"goravel/app/listeners"
	"goravel/config"
	"goravel/routes"
)

// #region boot
func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithJobs(Jobs).
		WithCommands(Commands).
		WithSeeders(Seeders).
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Web()
			routes.Api()
			routes.Grpc()
		}).
		// #region events
		// Each event maps to its listeners, like $listen in EventServiceProvider
		WithEvents(func() map[event.Event][]event.Listener {
			return map[event.Event][]event.Listener{
				&events.FixtureResultRecorded{}: {
					&listeners.NotifyClubAdmins{},
					&listeners.ForgetLeagueTable{},
				},
			}
		}).
		// #endregion events
		WithSchedule(Schedule).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}

// #endregion boot
