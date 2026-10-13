// Laravel equivalent: tests/Feature/ScheduleTest.php - see docs/chapters/17-scheduling.md
package feature

import (
	"testing"
	"time"

	"github.com/goravel/framework/contracts/schedule"
	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/tests"
)

type ScheduleTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestScheduleTestSuite(t *testing.T) {
	suite.Run(t, new(ScheduleTestSuite))
}

func (s *ScheduleTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
}

// event finds a scheduled event by command or name.
func (s *ScheduleTestSuite) event(match func(schedule.Event) bool) schedule.Event {
	for _, event := range facades.Schedule().Events() {
		if match(event) {
			return event
		}
	}
	s.FailNow("scheduled event not found")

	return nil
}

// #region test
func (s *ScheduleTestSuite) TestSyncRunsEveryFifteenMinutesOnSaturdays() {
	sync := s.event(func(e schedule.Event) bool { return e.GetCommand() == "sync:fixtures" })

	s.Equal("*/15 * * * 6", sync.GetCron())
	s.True(sync.GetSkipIfStillRunning())
}

func (s *ScheduleTestSuite) TestPruneDeletesOnlyOldNotifications() {
	prune := s.event(func(e schedule.Event) bool { return e.GetName() == "prune-old-notifications" })
	s.Equal("00 03 * * *", prune.GetCron())

	var club models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&club))
	var admin models.User
	s.Require().NoError(facades.Orm().Factory().Create(&admin, map[string]any{"ClubID": &club.ID}))
	var fixture models.Fixture
	s.Require().NoError(facades.Orm().Factory().Create(&fixture, map[string]any{"HomeClubID": club.ID, "AwayClubID": club.ID}))

	old := models.ResultNotification{UserID: admin.ID, FixtureID: fixture.ID, Message: "old"}
	recent := models.ResultNotification{UserID: admin.ID, FixtureID: fixture.ID, Message: "recent"}
	s.Require().NoError(facades.Orm().Query().Create(&old))
	s.Require().NoError(facades.Orm().Query().Create(&recent))
	_, err := facades.Orm().Query().Model(&old).Update("created_at", time.Now().UTC().AddDate(0, 0, -31))
	s.Require().NoError(err)

	// Run the task's callback directly, without waiting for 3am
	prune.GetCallback()()

	var remaining []models.ResultNotification
	s.Require().NoError(facades.Orm().Query().Get(&remaining))
	s.Require().Len(remaining, 1)
	s.Equal("recent", remaining[0].Message)
}

// #endregion test
