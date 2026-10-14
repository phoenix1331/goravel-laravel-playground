package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/tests"
)

type SyncFixturesTestSuite struct {
	suite.Suite
	tests.TestCase
	fixture models.Fixture
}

func TestSyncFixturesTestSuite(t *testing.T) {
	suite.Run(t, new(SyncFixturesTestSuite))
}

func (s *SyncFixturesTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())

	var ashford, barnsley models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&ashford, map[string]any{"Name": "Ashford"}))
	s.Require().NoError(facades.Orm().Factory().Create(&barnsley, map[string]any{"Name": "Barnsley"}))
	s.Require().NoError(facades.Orm().Factory().Create(&s.fixture, map[string]any{
		"HomeClubID": ashford.ID,
		"AwayClubID": barnsley.ID,
	}))
}

// reload reads the fixture back from the database.
func (s *SyncFixturesTestSuite) reload() models.Fixture {
	var fixture models.Fixture
	s.Require().NoError(facades.Orm().Query().FindOrFail(&fixture, s.fixture.ID))

	return fixture
}

func (s *SyncFixturesTestSuite) TestRecordsResultsFromTheFeed() {
	// Go runs each package's tests from the package directory, so this path
	// is tests/feature/testdata/..., the conventional home for test fixtures
	err := facades.Artisan().Call("sync:fixtures --file=testdata/fixture_results.json")
	s.Require().NoError(err)

	fixture := s.reload()
	s.True(fixture.IsFinished())
	s.Equal(3, *fixture.HomeScore)
	s.Equal(1, *fixture.AwayScore)
}

func (s *SyncFixturesTestSuite) TestDryRunChangesNothing() {
	err := facades.Artisan().Call("sync:fixtures --file=testdata/fixture_results.json --dry-run")
	s.Require().NoError(err)

	s.Equal(models.Scheduled, s.reload().Status)
}

func (s *SyncFixturesTestSuite) TestMissingFileFails() {
	err := facades.Artisan().Call("sync:fixtures --file=testdata/missing.json")

	s.Error(err)
}
