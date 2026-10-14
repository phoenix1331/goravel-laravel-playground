package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/tests"
)

type RoutingTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestRoutingTestSuite(t *testing.T) {
	suite.Run(t, new(RoutingTestSuite))
}

func (s *RoutingTestSuite) TestClosureRouteInsideApiGroup() {
	response, err := s.Http(s.T()).Get("/api/league")
	s.Require().NoError(err)

	response.AssertOk()
	response.AssertExactJson(map[string]any{"name": "Test League"})
}

func (s *RoutingTestSuite) TestNamedHandlerRoute() {
	response, err := s.Http(s.T()).Get("/api/league/points")
	s.Require().NoError(err)

	response.AssertOk()
	// JSON numbers decode into Go as float64, so compare against float64 values
	response.AssertExactJson(map[string]any{"win": 3.0, "draw": 1.0, "loss": 0.0})
}

func (s *RoutingTestSuite) TestRouteOutsideGroupIsNotPrefixed() {
	response, err := s.Http(s.T()).Get("/league")
	s.Require().NoError(err)

	response.AssertNotFound()
	response.AssertExactJson(map[string]any{"message": "Not Found"})
}
