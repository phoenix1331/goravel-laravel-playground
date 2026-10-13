// Laravel equivalent: tests/Feature/ClubValidationTest.php - see docs/chapters/07-validation.md
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

type ClubValidationTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestClubValidationTestSuite(t *testing.T) {
	suite.Run(t, new(ClubValidationTestSuite))
}

func (s *ClubValidationTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
}

// validationErrors posts a body and returns the "errors" object from the 422.
func (s *ClubValidationTestSuite) validationErrors(method, uri, body string) map[string]map[string]string {
	request := s.Http(s.T()).WithHeader("Content-Type", "application/json")

	send := request.Post
	if method == "PUT" {
		send = request.Put
	}
	response, err := send(uri, strings.NewReader(body))
	s.Require().NoError(err)
	response.AssertUnprocessableEntity()

	var decoded struct{ Errors map[string]map[string]string }
	s.Require().NoError(response.Bind(&decoded))

	return decoded.Errors
}

// #region test
func (s *ClubValidationTestSuite) TestRequiredFields() {
	errors := s.validationErrors("POST", "/api/clubs", `{}`)

	s.Contains(errors["name"], "required")
	s.Contains(errors["ground"], "required")
	s.NotContains(errors, "founded_year") // nullable
}

func (s *ClubValidationTestSuite) TestWhitespaceIsTrimmedBeforeRequiredCheck() {
	errors := s.validationErrors("POST", "/api/clubs", `{"name":"   ","ground":"Park"}`)

	s.Contains(errors["name"], "required")
}

func (s *ClubValidationTestSuite) TestFoundedYearRange() {
	errors := s.validationErrors("POST", "/api/clubs", `{"name":"Old FC","ground":"Park","founded_year":1800}`)

	s.Equal("The oldest football club was founded in 1857.", errors["founded_year"]["min"])
}

func (s *ClubValidationTestSuite) TestNameMustBeUnique() {
	s.Require().NoError(facades.Orm().Query().Create(&models.Club{Name: "Taken FC", Ground: "Park"}))

	errors := s.validationErrors("POST", "/api/clubs", `{"name":"Taken FC","ground":"Park"}`)

	s.Equal("A club with this name already exists.", errors["name"]["unique"])
}

// #endregion test

func (s *ClubValidationTestSuite) TestUpdateIgnoresOwnNameInUniqueCheck() {
	club := models.Club{Name: "Same Name FC", Ground: "Park"}
	s.Require().NoError(facades.Orm().Query().Create(&club))

	response, err := s.Http(s.T()).WithHeader("Content-Type", "application/json").
		Put(fmt.Sprintf("/api/clubs/%d", club.ID), strings.NewReader(`{"name":"Same Name FC","ground":"New Park"}`))
	s.Require().NoError(err)

	response.AssertOk()
}

func (s *ClubValidationTestSuite) TestUpdateRejectsAnotherClubsName() {
	s.Require().NoError(facades.Orm().Query().Create(&models.Club{Name: "First FC", Ground: "Park"}))
	second := models.Club{Name: "Second FC", Ground: "Park"}
	s.Require().NoError(facades.Orm().Query().Create(&second))

	errors := s.validationErrors("PUT", fmt.Sprintf("/api/clubs/%d", second.ID), `{"name":"First FC","ground":"Park"}`)

	s.Contains(errors["name"], "unique")
}
