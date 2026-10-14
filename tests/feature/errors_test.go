// Laravel equivalent: tests/Feature/RecordResultTwiceTest.php
package feature

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/results"
	"goravel/tests"
)

type ErrorsTestSuite struct {
	suite.Suite
	tests.TestCase
	fixture models.Fixture
	admin   models.User
}

func TestErrorsTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorsTestSuite))
}

func (s *ErrorsTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())

	var home, away models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&home))
	s.Require().NoError(facades.Orm().Factory().Create(&away))
	s.admin = s.CreateAdmin(s.T(), home.ID)
	s.Require().NoError(facades.Orm().Factory().Create(&s.fixture, map[string]any{"HomeClubID": home.ID, "AwayClubID": away.ID}))
}

func (s *ErrorsTestSuite) TestRecordingTwiceReturnsTheSentinel() {
	s.Require().NoError(results.Record(&s.fixture, 1, 0))

	err := results.Record(&s.fixture, 3, 3)

	s.True(errors.Is(err, results.ErrAlreadyRecorded))
	s.Equal(1, *s.fixture.HomeScore) // unchanged
}

func (s *ErrorsTestSuite) TestTheApiTurnsTheSentinelInto409() {
	uri := fmt.Sprintf("/api/fixtures/%d/result", s.fixture.ID)
	request := s.LoginAs(s.T(), s.admin)

	response, err := request.Post(uri, strings.NewReader(`{"home_score":1,"away_score":0}`))
	s.Require().NoError(err)
	response.AssertOk()

	response, err = request.Post(uri, strings.NewReader(`{"home_score":3,"away_score":3}`))
	s.Require().NoError(err)
	response.AssertConflict()
	response.AssertExactJson(map[string]any{"message": "This fixture already has a result."})
}
