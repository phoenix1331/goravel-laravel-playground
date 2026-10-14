// Laravel equivalent: tests/Feature/LeagueTableCacheTest.php
package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/results"
	"goravel/app/services"
	"goravel/tests"
)

type CacheTestSuite struct {
	suite.Suite
	tests.TestCase
	fixture models.Fixture
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}

func (s *CacheTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
	facades.Cache().Forget(services.LeagueTableCacheKey)

	var home, away models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&home, map[string]any{"Name": "Ashford"}))
	s.Require().NoError(facades.Orm().Factory().Create(&away, map[string]any{"Name": "Barnsley"}))
	s.Require().NoError(facades.Orm().Factory().Create(&s.fixture, map[string]any{"HomeClubID": home.ID, "AwayClubID": away.ID}))
}

// leader returns the club at the top of the table and their points.
func (s *CacheTestSuite) leader() (string, int) {
	response, err := s.Http(s.T()).Get("/api/league/table")
	s.Require().NoError(err)
	response.AssertOk()

	var body struct {
		Data []struct {
			Club   string
			Points int
		}
	}
	s.Require().NoError(response.Bind(&body))

	return body.Data[0].Club, body.Data[0].Points
}

func (s *CacheTestSuite) TestTableIsCachedAfterFirstRequest() {
	s.leader()

	s.True(facades.Cache().Has(services.LeagueTableCacheKey))
}

func (s *CacheTestSuite) TestChangesBehindTheCachesBackAreNotSeen() {
	_, points := s.leader()
	s.Zero(points)

	// Update the database directly, skipping results.Record and its event
	_, err := facades.Orm().Query().Model(&s.fixture).Update(map[string]any{
		"status": models.Finished, "home_score": 1, "away_score": 0,
	})
	s.Require().NoError(err)

	// Still the cached table
	_, points = s.leader()
	s.Zero(points)
}

func (s *CacheTestSuite) TestRecordingAResultClearsTheCache() {
	s.leader()

	s.Require().NoError(results.Record(&s.fixture, 1, 0))

	s.False(facades.Cache().Has(services.LeagueTableCacheKey))
	club, points := s.leader()
	s.Equal("Ashford", club)
	s.Equal(3, points)
}
