package feature

import (
	"testing"

	"github.com/goravel/framework/contracts/http/client"
	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/tests"
)

type HttpClientTestSuite struct {
	suite.Suite
	tests.TestCase
	fixture models.Fixture
}

func TestHttpClientTestSuite(t *testing.T) {
	suite.Run(t, new(HttpClientTestSuite))
}

func (s *HttpClientTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())

	var ashford, barnsley models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&ashford, map[string]any{"Name": "Ashford"}))
	s.Require().NoError(facades.Orm().Factory().Create(&barnsley, map[string]any{"Name": "Barnsley"}))
	s.Require().NoError(facades.Orm().Factory().Create(&s.fixture, map[string]any{
		"HomeClubID": ashford.ID, "AwayClubID": barnsley.ID,
	}))

	// Point sync:fixtures at the HTTP feed for these tests only
	facades.Config().Add("club.feed.driver", "http")
}

// TearDownTest puts things back: fakes and config are process-wide, and
// every test in the package shares the process.
func (s *HttpClientTestSuite) TearDownTest() {
	facades.Http().Reset()
	facades.Config().Add("club.feed.driver", "file")
}

func (s *HttpClientTestSuite) reload() models.Fixture {
	var fixture models.Fixture
	s.Require().NoError(facades.Orm().Query().FindOrFail(&fixture, s.fixture.ID))

	return fixture
}

func (s *HttpClientTestSuite) TestSyncReadsResultsFromTheApi() {
	// Like Http::fake(['feed.test/*' => Http::response([...], 200)])
	http := facades.Http()
	http.Fake(map[string]any{
		"https://feed.test/api/results": http.Response().Json(200, []map[string]any{
			{"home": "Ashford", "away": "Barnsley", "home_score": 2, "away_score": 0},
		}),
	}).PreventStrayRequests()

	s.Require().NoError(facades.Artisan().Call("sync:fixtures"))

	s.Equal(2, *s.reload().HomeScore)
	s.True(http.AssertSent(func(request client.Request) bool {
		return request.Method() == "GET" && request.Header("Accept") == "application/json"
	}))
}

func (s *HttpClientTestSuite) TestServerErrorFailsTheCommandAndChangesNothing() {
	facades.Http().Fake(map[string]any{"feed.test/*": 503})

	err := facades.Artisan().Call("sync:fixtures")

	s.ErrorContains(err, "status 503")
	s.Equal(models.Scheduled, s.reload().Status)
}

func (s *HttpClientTestSuite) TestStrayRequestsAreBlocked() {
	// Nothing matches the feed URL, and PreventStrayRequests stops the
	// request reaching the network
	facades.Http().Fake(map[string]any{"other.test/*": 200}).PreventStrayRequests()

	s.Error(facades.Artisan().Call("sync:fixtures"))
}
