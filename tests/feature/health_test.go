package feature

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/tests"
)

func TestMain(m *testing.M) {
	os.Exit(tests.Main(m))
}

type HealthTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestHealthTestSuite(t *testing.T) {
	suite.Run(t, new(HealthTestSuite))
}

// SetupTest runs before each test, like setUp() in PHPUnit.
func (s *HealthTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
}

func (s *HealthTestSuite) TestHealthReportsDatabaseUp() {
	response, err := s.Http(s.T()).Get("/health")
	s.Require().NoError(err)

	response.AssertOk()
	response.AssertExactJson(map[string]any{
		"status":   "ok",
		"database": "up",
	})
}
