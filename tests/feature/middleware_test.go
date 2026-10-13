// Laravel equivalent: tests/Feature/EnsureJsonRequestTest.php - see docs/chapters/08-middleware.md
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

// #region test
func (s *MiddlewareTestSuite) TestRejectsNonJsonWrites() {
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
	response, err := s.Http(s.T()).
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

// #endregion test
