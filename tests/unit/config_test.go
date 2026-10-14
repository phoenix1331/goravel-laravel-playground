package unit

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/tests"
)

func TestMain(m *testing.M) {
	os.Exit(tests.Main(m))
}

type ConfigTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestTypedLookups() {
	config := facades.Config()

	// .env.testing overrides the default
	s.Equal("Test League", config.GetString("club.league_name"))

	// Not set in .env.testing, so the default in config/club.go applies.
	// GetInt converts it: the caller gets an int, not a string or any.
	s.Equal(3, config.GetInt("club.points.win"))
	s.Equal(1, config.GetInt("club.points.draw"))
	s.Equal(0, config.GetInt("club.points.loss"))

	// Missing keys fall back to the default passed in
	s.Equal(42, config.GetInt("club.does_not_exist", 42))
}
