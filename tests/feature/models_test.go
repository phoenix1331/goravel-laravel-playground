package feature

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/tests"
)

type ModelsTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestModelsTestSuite(t *testing.T) {
	suite.Run(t, new(ModelsTestSuite))
}

func (s *ModelsTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
}

func (s *ModelsTestSuite) TestCreateFillsIdAndTimestamps() {
	club := models.Club{Name: "Ashford Athletic", Ground: "The Rec"}
	s.Require().NoError(facades.Orm().Query().Create(&club))

	// Create writes the generated ID and timestamps back into the struct,
	// because it was given a pointer (&club), not a copy
	s.NotZero(club.ID)
	s.NotNil(club.CreatedAt)
}

func (s *ModelsTestSuite) TestTypedFieldsRoundTrip() {
	club := models.Club{Name: "Ashford Athletic", Ground: "The Rec"}
	s.Require().NoError(facades.Orm().Query().Create(&club))
	player := models.Player{ClubID: club.ID, Name: "Sam Keeper", Position: models.Goalkeeper, SquadNumber: 1}
	s.Require().NoError(facades.Orm().Query().Create(&player))

	var found models.Player
	s.Require().NoError(facades.Orm().Query().Where("squad_number", 1).FirstOrFail(&found))

	s.Equal(models.Goalkeeper, found.Position)
	s.Equal(1, found.SquadNumber)
}

func (s *ModelsTestSuite) TestPasswordIsNeverSerialised() {
	user := models.User{Name: "Alex Admin", Email: "alex@example.com", Password: "hashed"}

	encoded, err := json.Marshal(user)
	s.Require().NoError(err)

	s.NotContains(string(encoded), "password")
	s.NotContains(string(encoded), "hashed")
}

func (s *ModelsTestSuite) TestFixtureResultOnlyWhenFinished() {
	home, away := 2, 1
	fixture := models.Fixture{HomeClubID: 1, AwayClubID: 2, Status: models.Scheduled, KickOffAt: time.Now()}

	_, ok := fixture.Result()
	s.False(ok)

	fixture.Status, fixture.HomeScore, fixture.AwayScore = models.Finished, &home, &away
	result, ok := fixture.Result()
	s.True(ok)
	s.Equal(2, result.HomeScore)
}
