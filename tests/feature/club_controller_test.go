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

type ClubControllerTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestClubControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ClubControllerTestSuite))
}

func (s *ClubControllerTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
}

// createClub makes a club with the factory, overriding only the name.
func (s *ClubControllerTestSuite) createClub(name string) models.Club {
	var club models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&club, map[string]any{"Name": name}))

	return club
}

func (s *ClubControllerTestSuite) TestIndexListsClubsByName() {
	s.createClub("Riverside Rovers")
	s.createClub("Ashford Athletic")

	response, err := s.Http(s.T()).Get("/api/clubs")
	s.Require().NoError(err)
	response.AssertOk()

	// Bind decodes the JSON into a typed struct: no float64 surprises
	var body struct{ Data []models.Club }
	s.Require().NoError(response.Bind(&body))
	s.Require().Len(body.Data, 2)
	s.Equal("Ashford Athletic", body.Data[0].Name)
	s.Equal("Riverside Rovers", body.Data[1].Name)
}

func (s *ClubControllerTestSuite) TestShowReturns404ForMissingClub() {
	response, err := s.Http(s.T()).Get("/api/clubs/999")
	s.Require().NoError(err)

	response.AssertNotFound()
	response.AssertExactJson(map[string]any{"message": "Club not found"})
}

func (s *ClubControllerTestSuite) TestStoreCreatesClub() {
	var user models.User
	s.Require().NoError(facades.Orm().Factory().Create(&user))

	body := strings.NewReader(`{"name":"Hillside United","ground":"Hillside Rec","founded_year":1987}`)
	response, err := s.LoginAs(s.T(), user).Post("/api/clubs", body)
	s.Require().NoError(err)
	response.AssertCreated()

	var club models.Club
	s.Require().NoError(facades.Orm().Query().Where("name", "Hillside United").FirstOrFail(&club))
	s.Equal("Hillside Rec", club.Ground)
	s.Require().NotNil(club.FoundedYear)
	s.Equal(1987, *club.FoundedYear)
}

func (s *ClubControllerTestSuite) TestUpdateChangesClub() {
	club := s.createClub("Old Name")

	admin := s.CreateAdmin(s.T(), club.ID)

	body := strings.NewReader(`{"name":"New Name","ground":"New Ground"}`)
	response, err := s.LoginAs(s.T(), admin).Put(fmt.Sprintf("/api/clubs/%d", club.ID), body)
	s.Require().NoError(err)
	response.AssertOk()

	var updated struct{ Data models.Club }
	s.Require().NoError(response.Bind(&updated))
	s.Equal("New Name", updated.Data.Name)
	s.Equal("New Ground", updated.Data.Ground)
}

func (s *ClubControllerTestSuite) TestDestroyDeletesClub() {
	club := s.createClub("Doomed FC")

	admin := s.CreateAdmin(s.T(), club.ID)

	response, err := s.LoginAs(s.T(), admin).Delete(fmt.Sprintf("/api/clubs/%d", club.ID), nil)
	s.Require().NoError(err)
	response.AssertNoContent()

	count, err := facades.Orm().Query().Model(&models.Club{}).Count()
	s.Require().NoError(err)
	s.Zero(count)
}
