package feature

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/goravel/framework/support/path"
	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/results"
	"goravel/tests"
)

type LoggingTestSuite struct {
	suite.Suite
	tests.TestCase
	logFile string
}

func TestLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingTestSuite))
}

func (s *LoggingTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())

	// .env.testing sends the results channel to its own file
	s.logFile = path.Storage("logs/results-testing.log")
}

// logSize returns the log file's current size. The logger keeps the file
// open for the life of the process, so deleting it between tests would leave
// the logger writing to a file nobody can see; read only what is new instead.
func (s *LoggingTestSuite) logSize() int64 {
	info, err := os.Stat(s.logFile)
	if err != nil {
		return 0
	}

	return info.Size()
}

func (s *LoggingTestSuite) TestRecordingAResultWritesAJsonLine() {
	var home, away models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&home))
	s.Require().NoError(facades.Orm().Factory().Create(&away))
	var fixture models.Fixture
	s.Require().NoError(facades.Orm().Factory().Create(&fixture, map[string]any{"HomeClubID": home.ID, "AwayClubID": away.ID}))

	before := s.logSize()
	s.Require().NoError(results.Record(&fixture, 4, 2))

	contents, err := os.ReadFile(s.logFile)
	s.Require().NoError(err)
	lines := strings.Split(strings.TrimSpace(string(contents[before:])), "\n")
	s.Require().Len(lines, 1)

	// {"environment":"testing","extra":{...},"level":"info","message":"result recorded","time":"..."}
	var entry struct {
		Level   string
		Message string
		Extra   struct {
			FixtureID uint `json:"fixture_id"`
			HomeScore int  `json:"home_score"`
			AwayScore int  `json:"away_score"`
		}
	}
	s.Require().NoError(json.Unmarshal([]byte(lines[0]), &entry))
	s.Equal("info", entry.Level)
	s.Equal("result recorded", entry.Message)
	s.Equal(fixture.ID, entry.Extra.FixtureID)
	s.Equal(4, entry.Extra.HomeScore)
	s.Equal(2, entry.Extra.AwayScore)
}
