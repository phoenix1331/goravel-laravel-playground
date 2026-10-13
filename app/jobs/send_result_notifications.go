// Laravel equivalent: app/Jobs/SendResultNotifications.php - see docs/chapters/16-queues.md
package jobs

import (
	"fmt"

	"github.com/goravel/framework/contracts/queue"
	"github.com/spf13/cast"

	"goravel/app/facades"
	"goravel/app/models"
)

// #region job
// SendResultNotifications writes a notification for each admin of the two
// clubs in a fixture. It takes one argument: the fixture ID.
type SendResultNotifications struct{}

// Signature identifies the job on the queue. The worker looks the job up by
// this name, which is why every job is registered in bootstrap/jobs.go.
func (j *SendResultNotifications) Signature() string {
	return "send_result_notifications"
}

// Args builds the job's argument list for dispatching.
func (j *SendResultNotifications) Args(fixtureID uint) []queue.Arg {
	return []queue.Arg{{Type: "uint", Value: fixtureID}}
}

// Handle does the work. On the Redis queue it runs in the worker, later and
// outside any HTTP request; on the sync queue it runs straight away.
func (j *SendResultNotifications) Handle(args ...any) error {
	if len(args) == 0 {
		return fmt.Errorf("send result notifications: missing fixture id")
	}
	fixtureID := cast.ToUint(args[0])

	var fixture models.Fixture
	err := facades.Orm().Query().With("HomeClub").With("AwayClub").FindOrFail(&fixture, fixtureID)
	if err != nil {
		return fmt.Errorf("send result notifications: %w", err)
	}

	var admins []models.User
	err = facades.Orm().Query().
		WhereIn("club_id", []any{fixture.HomeClubID, fixture.AwayClubID}).
		Get(&admins)
	if err != nil {
		return fmt.Errorf("send result notifications: %w", err)
	}

	message := fmt.Sprintf("Result: %s %d-%d %s",
		fixture.HomeClub.Name, *fixture.HomeScore, *fixture.AwayScore, fixture.AwayClub.Name)
	for _, admin := range admins {
		notification := models.ResultNotification{UserID: admin.ID, FixtureID: fixture.ID, Message: message}
		if err := facades.Orm().Query().Create(&notification); err != nil {
			return fmt.Errorf("send result notifications: %w", err)
		}
	}

	return nil
}

// #endregion job
