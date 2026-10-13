// Laravel equivalent: tests/Feature/RelationshipsTest.php - see docs/chapters/10-relationships.md
package feature

import (
	"fmt"
	"strings"
	"testing"
	"time"

	testinghttp "github.com/goravel/framework/contracts/testing/http"
	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/tests"
)

type RelationshipsTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestRelationshipsTestSuite(t *testing.T) {
	suite.Run(t, new(RelationshipsTestSuite))
}

func (s *RelationshipsTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
}

func (s *RelationshipsTestSuite) club(name string) models.Club {
	var club models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&club, map[string]any{"Name": name}))

	return club
}

// postPlayer adds a player as the club's own admin.
func (s *RelationshipsTestSuite) postPlayer(clubID uint, body string) testinghttp.Response {
	response, err := s.LoginAs(s.T(), s.CreateAdmin(s.T(), clubID)).
		Post(fmt.Sprintf("/api/clubs/%d/players", clubID), strings.NewReader(body))
	s.Require().NoError(err)

	return response
}

// #region test
func (s *RelationshipsTestSuite) TestClubPlayersAreLoadedInSquadOrder() {
	club := s.club("Ashford Athletic")
	other := s.club("Riverside Rovers")
	for _, player := range []models.Player{
		{ClubID: club.ID, Name: "Nine", Position: models.Forward, SquadNumber: 9},
		{ClubID: club.ID, Name: "One", Position: models.Goalkeeper, SquadNumber: 1},
		{ClubID: other.ID, Name: "Elsewhere", Position: models.Defender, SquadNumber: 5},
	} {
		s.Require().NoError(facades.Orm().Query().Create(&player))
	}

	response, err := s.Http(s.T()).Get(fmt.Sprintf("/api/clubs/%d/players", club.ID))
	s.Require().NoError(err)
	response.AssertOk()

	var body struct{ Data models.Club }
	s.Require().NoError(response.Bind(&body))
	s.Require().Len(body.Data.Players, 2)
	s.Equal("One", body.Data.Players[0].Name)
	s.Equal("Nine", body.Data.Players[1].Name)
}

func (s *RelationshipsTestSuite) TestFixturesEagerLoadBothClubs() {
	home, away := s.club("Home Town"), s.club("Away Days")
	fixture := models.Fixture{
		HomeClubID: home.ID, AwayClubID: away.ID,
		KickOffAt: time.Date(2026, 11, 7, 15, 0, 0, 0, time.UTC),
		Venue:     "Home Town Park", Status: models.Scheduled,
	}
	s.Require().NoError(facades.Orm().Query().Create(&fixture))

	response, err := s.Http(s.T()).Get("/api/fixtures")
	s.Require().NoError(err)
	response.AssertOk()

	var body struct{ Data []models.Fixture }
	s.Require().NoError(response.Bind(&body))
	s.Require().Len(body.Data, 1)
	s.Equal("Home Town", body.Data[0].HomeClub.Name)
	s.Equal("Away Days", body.Data[0].AwayClub.Name)
}

// #endregion test

func (s *RelationshipsTestSuite) TestStorePlayerSetsClub() {
	club := s.club("Ashford Athletic")

	s.postPlayer(club.ID, `{"name":"New Signing","position":" Forward ","squad_number":10}`).AssertCreated()

	var player models.Player
	s.Require().NoError(facades.Orm().Query().Where("club_id", club.ID).FirstOrFail(&player))
	s.Equal(models.Forward, player.Position) // trimmed and lower-cased by filters
}

func (s *RelationshipsTestSuite) TestStorePlayerRejectsUnknownPosition() {
	club := s.club("Ashford Athletic")

	response := s.postPlayer(club.ID, `{"name":"Who","position":"striker","squad_number":10}`)

	response.AssertUnprocessableEntity()
	body, err := response.Json()
	s.Require().NoError(err)
	s.Contains(body["errors"], "position")
}

func (s *RelationshipsTestSuite) TestSquadNumberIsUniquePerClub() {
	club, other := s.club("Ashford Athletic"), s.club("Riverside Rovers")
	s.Require().NoError(facades.Orm().Query().Create(&models.Player{ClubID: club.ID, Name: "Ten", Position: models.Midfielder, SquadNumber: 10}))

	s.postPlayer(club.ID, `{"name":"Also Ten","position":"forward","squad_number":10}`).AssertUnprocessableEntity()

	// The same number at a different club is fine
	s.postPlayer(other.ID, `{"name":"Their Ten","position":"forward","squad_number":10}`).AssertCreated()
}
