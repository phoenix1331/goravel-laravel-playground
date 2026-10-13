// Laravel equivalent: tests/Feature/AuthorisationTest.php - see docs/chapters/20-authorisation.md
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

type AuthorisationTestSuite struct {
	suite.Suite
	tests.TestCase
	home, away models.Club
	homeAdmin  models.User
	awayAdmin  models.User
	fixture    models.Fixture
}

func TestAuthorisationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorisationTestSuite))
}

func (s *AuthorisationTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())

	s.Require().NoError(facades.Orm().Factory().Create(&s.home, map[string]any{"Name": "Ashford"}))
	s.Require().NoError(facades.Orm().Factory().Create(&s.away, map[string]any{"Name": "Barnsley"}))
	s.homeAdmin = s.CreateAdmin(s.T(), s.home.ID)
	s.awayAdmin = s.CreateAdmin(s.T(), s.away.ID)
	s.Require().NoError(facades.Orm().Factory().Create(&s.fixture, map[string]any{
		"HomeClubID": s.home.ID, "AwayClubID": s.away.ID,
	}))
}

// #region test
func (s *AuthorisationTestSuite) TestAdminsCannotEditAnotherClub() {
	response, err := s.LoginAs(s.T(), s.awayAdmin).
		Put(fmt.Sprintf("/api/clubs/%d", s.home.ID), strings.NewReader(`{"name":"Hijacked","ground":"Park"}`))
	s.Require().NoError(err)

	response.AssertForbidden()
	response.AssertExactJson(map[string]any{"message": "You can only manage your own club."})
}

func (s *AuthorisationTestSuite) TestAdminsCannotAddPlayersToAnotherClub() {
	response, err := s.LoginAs(s.T(), s.awayAdmin).
		Post(fmt.Sprintf("/api/clubs/%d/players", s.home.ID), strings.NewReader(`{"name":"Mole","position":"forward","squad_number":9}`))
	s.Require().NoError(err)

	response.AssertForbidden()
}

func (s *AuthorisationTestSuite) TestOnlyTheHomeAdminRecordsTheResult() {
	uri := fmt.Sprintf("/api/fixtures/%d/result", s.fixture.ID)
	body := `{"home_score":0,"away_score":5}`

	response, err := s.LoginAs(s.T(), s.awayAdmin).Post(uri, strings.NewReader(body))
	s.Require().NoError(err)
	response.AssertForbidden()

	response, err = s.LoginAs(s.T(), s.homeAdmin).Post(uri, strings.NewReader(body))
	s.Require().NoError(err)
	response.AssertOk()
}

func (s *AuthorisationTestSuite) TestCreatingAClubMakesYouItsAdmin() {
	var newcomer models.User
	s.Require().NoError(facades.Orm().Factory().Create(&newcomer))

	response, err := s.LoginAs(s.T(), newcomer).Post("/api/clubs", strings.NewReader(`{"name":"Crewe","ground":"Gresty Road"}`))
	s.Require().NoError(err)
	response.AssertCreated()

	var club models.Club
	s.Require().NoError(facades.Orm().Query().Where("name", "Crewe").FirstOrFail(&club))
	s.Require().NoError(facades.Orm().Query().FindOrFail(&newcomer, newcomer.ID))
	s.Require().NotNil(newcomer.ClubID)
	s.Equal(club.ID, *newcomer.ClubID)
}

func (s *AuthorisationTestSuite) TestAdminsCannotCreateASecondClub() {
	response, err := s.LoginAs(s.T(), s.homeAdmin).Post("/api/clubs", strings.NewReader(`{"name":"Second","ground":"Park"}`))
	s.Require().NoError(err)

	response.AssertForbidden()
}

// #endregion test
