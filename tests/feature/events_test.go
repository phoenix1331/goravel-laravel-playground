package feature

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/tests"
)

type EventsTestSuite struct {
	suite.Suite
	tests.TestCase
	fixture  models.Fixture
	admins   map[string]models.User
	outsider models.User
}

func TestEventsTestSuite(t *testing.T) {
	suite.Run(t, new(EventsTestSuite))
}

// SetupTest makes Ashford v Barnsley, an admin for each club, and an admin
// of a third club who should hear nothing.
func (s *EventsTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())

	s.admins = map[string]models.User{}
	clubs := map[string]models.Club{}
	for _, name := range []string{"Ashford", "Barnsley", "Crewe"} {
		var club models.Club
		s.Require().NoError(facades.Orm().Factory().Create(&club, map[string]any{"Name": name}))
		clubs[name] = club

		var admin models.User
		s.Require().NoError(facades.Orm().Factory().Create(&admin, map[string]any{"ClubID": &club.ID}))
		s.admins[name] = admin
	}

	s.Require().NoError(facades.Orm().Factory().Create(&s.fixture, map[string]any{
		"HomeClubID": clubs["Ashford"].ID,
		"AwayClubID": clubs["Barnsley"].ID,
	}))
}

func (s *EventsTestSuite) notificationsFor(user models.User) []models.ResultNotification {
	var notifications []models.ResultNotification
	s.Require().NoError(facades.Orm().Query().Where("user_id", user.ID).Get(&notifications))

	return notifications
}

func (s *EventsTestSuite) TestRecordingAResultNotifiesBothClubsAdmins() {
	response, err := s.LoginAs(s.T(), s.admins["Ashford"]).
		Post(fmt.Sprintf("/api/fixtures/%d/result", s.fixture.ID), strings.NewReader(`{"home_score":3,"away_score":1}`))
	s.Require().NoError(err)
	response.AssertOk()

	for _, club := range []string{"Ashford", "Barnsley"} {
		notifications := s.notificationsFor(s.admins[club])
		s.Require().Len(notifications, 1, club)
		s.Equal("Result: Ashford 3-1 Barnsley", notifications[0].Message)
	}
	s.Empty(s.notificationsFor(s.admins["Crewe"]))
}

func (s *EventsTestSuite) TestSyncCommandFiresTheSameEvent() {
	// sync:fixtures goes through results.Record too, so the listener runs
	s.Require().NoError(facades.Artisan().Call("sync:fixtures --file=testdata/fixture_results.json"))

	s.Len(s.notificationsFor(s.admins["Ashford"]), 1)
}

func (s *EventsTestSuite) TestInvalidScoreIsRejectedAndNothingFires() {
	response, err := s.LoginAs(s.T(), s.admins["Ashford"]).
		Post(fmt.Sprintf("/api/fixtures/%d/result", s.fixture.ID), strings.NewReader(`{"home_score":-1}`))
	s.Require().NoError(err)

	response.AssertUnprocessableEntity()
	s.Empty(s.notificationsFor(s.admins["Ashford"]))
}
