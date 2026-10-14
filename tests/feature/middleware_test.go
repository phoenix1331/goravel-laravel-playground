package feature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/tests"
)

type MiddlewareTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
}

func (s *MiddlewareTestSuite) TestRejectsNonJsonWrites() {
	// No token needed: EnsureJson wraps the whole /api group, so it runs
	// before the Authenticate middleware on the write routes
	response, err := s.Http(s.T()).
		WithHeader("Content-Type", "application/x-www-form-urlencoded").
		Post("/api/clubs", strings.NewReader("name=Form+FC&ground=Park"))
	s.Require().NoError(err)

	response.AssertUnsupportedMediaType()

	// The controller never ran, so nothing was saved
	count, err := facades.Orm().Query().Model(&models.Club{}).Count()
	s.Require().NoError(err)
	s.Zero(count)
}

func (s *MiddlewareTestSuite) TestAllowsJsonWrites() {
	var user models.User
	s.Require().NoError(facades.Orm().Factory().Create(&user))

	response, err := s.LoginAs(s.T(), user).
		WithHeader("Content-Type", "application/json; charset=utf-8").
		Post("/api/clubs", strings.NewReader(`{"name":"Json FC","ground":"Park"}`))
	s.Require().NoError(err)

	response.AssertCreated()
}

func (s *MiddlewareTestSuite) TestReadsAreUnaffected() {
	response, err := s.Http(s.T()).Get("/api/clubs")
	s.Require().NoError(err)

	response.AssertOk()
}

func (s *MiddlewareTestSuite) TestRoutesOutsideTheGroupAreUnaffected() {
	response, err := s.Http(s.T()).Get("/health")
	s.Require().NoError(err)

	response.AssertOk()
}
