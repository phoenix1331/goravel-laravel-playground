package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/services"
	"goravel/tests"
)

type LeagueTestSuite struct {
	suite.Suite
	tests.TestCase
	clubs map[string]models.Club
}

func TestLeagueTestSuite(t *testing.T) {
	suite.Run(t, new(LeagueTestSuite))
}

// SetupTest builds a tiny league: Ashford beat Barnsley 3-0, Barnsley drew
// 1-1 with Crewe, and Dorking have not played. One fixture is still to come.
func (s *LeagueTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
	// The table is cached since chapter 18. Forget the key rather than Flush():
	// on Redis, Flush() runs FLUSHALL and empties every database on the server.
	facades.Cache().Forget(services.LeagueTableCacheKey)

	s.clubs = map[string]models.Club{}
	for _, name := range []string{"Ashford", "Barnsley", "Crewe", "Dorking"} {
		var club models.Club
		s.Require().NoError(facades.Orm().Factory().Create(&club, map[string]any{"Name": name}))
		s.clubs[name] = club
	}

	s.result("Ashford", "Barnsley", 3, 0)
	s.result("Barnsley", "Crewe", 1, 1)

	var upcoming models.Fixture
	s.Require().NoError(facades.Orm().Factory().Create(&upcoming, map[string]any{
		"HomeClubID": s.clubs["Crewe"].ID,
		"AwayClubID": s.clubs["Ashford"].ID,
	}))
}

func (s *LeagueTestSuite) result(home, away string, homeScore, awayScore int) {
	var fixture models.Fixture
	s.Require().NoError(facades.Orm().Factory().Create(&fixture, map[string]any{
		"HomeClubID": s.clubs[home].ID,
		"AwayClubID": s.clubs[away].ID,
		"Status":     string(models.Finished),
		"HomeScore":  &homeScore,
		"AwayScore":  &awayScore,
	}))
}

func (s *LeagueTestSuite) TestTableRanksClubsFromFinishedFixtures() {
	response, err := s.Http(s.T()).Get("/api/league/table")
	s.Require().NoError(err)
	response.AssertOk()

	var body struct {
		Data []struct {
			Position int
			Club     string
			Played   int
			Points   int
		}
	}
	s.Require().NoError(response.Bind(&body))

	s.Require().Len(body.Data, 4)
	s.Equal("Ashford", body.Data[0].Club)
	s.Equal(3, body.Data[0].Points)
	s.Equal(1, body.Data[0].Played) // the scheduled fixture is not counted
	s.Equal("Crewe", body.Data[1].Club)
	s.Equal("Barnsley", body.Data[2].Club)
	s.Equal("Dorking", body.Data[3].Club)
	s.Equal(0, body.Data[3].Played)
	s.Equal(4, body.Data[3].Position)
}

func (s *LeagueTestSuite) TestGoalsTotalsHomeAndAwayGoals() {
	response, err := s.Http(s.T()).Get("/api/league/goals")
	s.Require().NoError(err)
	response.AssertOk()

	var body struct {
		Data []struct {
			Club  string
			Goals int
		}
	}
	s.Require().NoError(response.Bind(&body))

	// Dorking have no finished fixtures, so the JOIN leaves them out
	s.Require().Len(body.Data, 3)
	s.Equal("Ashford", body.Data[0].Club)
	s.Equal(3, body.Data[0].Goals)
	s.Equal(1, body.Data[1].Goals) // Barnsley 0 + 1
	s.Equal(1, body.Data[2].Goals) // Crewe 1
}

func (s *LeagueTestSuite) TestBuilderGeneratesPortableSql() {
	sql := facades.DB().Table("fixtures").Where("status", models.Finished).ToSql().Get()

	s.Contains(sql, "status = ?")
}
