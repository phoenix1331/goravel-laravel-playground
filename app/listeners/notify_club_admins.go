// Laravel equivalent: app/Listeners/NotifyClubAdmins.php - see docs/chapters/15-events.md
package listeners

import (
	"fmt"

	"github.com/goravel/framework/contracts/event"
	"github.com/spf13/cast"

	"goravel/app/facades"
	"goravel/app/models"
)

// #region listener
// NotifyClubAdmins tells the admins of both clubs about a new result.
type NotifyClubAdmins struct{}

func (l *NotifyClubAdmins) Signature() string {
	return "notify_club_admins"
}

// Queue says whether to run on a queue. Enable: false runs it immediately, in
// the same request, like a Laravel listener without ShouldQueue.
func (l *NotifyClubAdmins) Queue(args ...any) event.Queue {
	return event.Queue{Enable: false}
}

// Handle receives the event's arguments as ...any, in the order they were
// dispatched. cast.ToUint converts whatever arrived into a uint.
func (l *NotifyClubAdmins) Handle(args ...any) error {
	if len(args) == 0 {
		return fmt.Errorf("notify club admins: missing fixture id")
	}
	fixtureID := cast.ToUint(args[0])

	var fixture models.Fixture
	err := facades.Orm().Query().With("HomeClub").With("AwayClub").FindOrFail(&fixture, fixtureID)
	if err != nil {
		return fmt.Errorf("notify club admins: %w", err)
	}

	var admins []models.User
	err = facades.Orm().Query().
		WhereIn("club_id", []any{fixture.HomeClubID, fixture.AwayClubID}).
		Get(&admins)
	if err != nil {
		return fmt.Errorf("notify club admins: %w", err)
	}

	message := fmt.Sprintf("Result: %s %d-%d %s",
		fixture.HomeClub.Name, *fixture.HomeScore, *fixture.AwayScore, fixture.AwayClub.Name)
	for _, admin := range admins {
		notification := models.ResultNotification{UserID: admin.ID, FixtureID: fixture.ID, Message: message}
		if err := facades.Orm().Query().Create(&notification); err != nil {
			return fmt.Errorf("notify club admins: %w", err)
		}
	}

	return nil
}

// #endregion listener
