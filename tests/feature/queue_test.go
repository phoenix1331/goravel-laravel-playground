// Laravel equivalent: tests/Feature/SendResultNotificationsTest.php
package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/jobs"
	"goravel/app/models"
	"goravel/tests"
)

type QueueTestSuite struct {
	suite.Suite
	tests.TestCase
	fixture models.Fixture
	admin   models.User
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}

func (s *QueueTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())

	var home, away models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&home, map[string]any{"Name": "Ashford"}))
	s.Require().NoError(facades.Orm().Factory().Create(&away, map[string]any{"Name": "Barnsley"}))
	s.Require().NoError(facades.Orm().Factory().Create(&s.admin, map[string]any{"ClubID": &home.ID}))

	homeScore, awayScore := 2, 2
	s.Require().NoError(facades.Orm().Factory().Create(&s.fixture, map[string]any{
		"HomeClubID": home.ID, "AwayClubID": away.ID,
		"Status": string(models.Finished), "HomeScore": &homeScore, "AwayScore": &awayScore,
	}))
}

func (s *QueueTestSuite) TestJobHandleWritesNotifications() {
	// A job is a plain struct: call Handle directly to test the work itself
	job := &jobs.SendResultNotifications{}
	s.Require().NoError(job.Handle(s.fixture.ID))

	var notification models.ResultNotification
	s.Require().NoError(facades.Orm().Query().Where("user_id", s.admin.ID).FirstOrFail(&notification))
	s.Equal("Result: Ashford 2-2 Barnsley", notification.Message)
}

func (s *QueueTestSuite) TestJobWithoutArgumentsFails() {
	s.Error((&jobs.SendResultNotifications{}).Handle())
}

func (s *QueueTestSuite) TestTestsUseTheSyncQueue() {
	// .env.testing sets QUEUE_CONNECTION=sync, so Dispatch runs the job
	// immediately and the test can check its effect on the next line
	s.Equal("sync", facades.Config().GetString("queue.default"))

	job := &jobs.SendResultNotifications{}
	s.Require().NoError(facades.Queue().Job(job, job.Args(s.fixture.ID)).Dispatch())

	count, err := facades.Orm().Query().Model(&models.ResultNotification{}).Count()
	s.Require().NoError(err)
	s.Equal(int64(1), count)
}
