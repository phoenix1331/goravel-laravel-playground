package feature

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/tests"
)

type MigrationsTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestMigrationsTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationsTestSuite))
}

func (s *MigrationsTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
}

func (s *MigrationsTestSuite) TestSchemaHasExpectedColumnsAndIndexes() {
	schema := facades.Schema()

	s.True(schema.HasColumns("clubs", []string{"id", "name", "ground", "founded_year", "short_name"}))
	s.True(schema.HasColumns("fixtures", []string{"home_club_id", "away_club_id", "kick_off_at", "status"}))
	s.True(schema.HasIndex("fixtures", "fixtures_kick_off_at_index"))
}

func (s *MigrationsTestSuite) TestDatabaseEnforcesUniqueSquadNumberPerClub() {
	club := models.Club{Name: "Ashford Athletic", Ground: "The Rec"}
	s.Require().NoError(facades.Orm().Query().Create(&club))
	s.Require().NoError(facades.Orm().Query().Create(&models.Player{ClubID: club.ID, Name: "A", Position: models.Forward, SquadNumber: 9}))

	// The unique index catches what validation might miss, such as two requests at once
	err := facades.Orm().Query().Create(&models.Player{ClubID: club.ID, Name: "B", Position: models.Forward, SquadNumber: 9})

	s.Error(err)
}

func (s *MigrationsTestSuite) TestDeletingClubCascadesToPlayers() {
	club := models.Club{Name: "Doomed FC", Ground: "The Rec"}
	s.Require().NoError(facades.Orm().Query().Create(&club))
	s.Require().NoError(facades.Orm().Query().Create(&models.Player{ClubID: club.ID, Name: "A", Position: models.Forward, SquadNumber: 9}))

	_, err := facades.Orm().Query().Delete(&club)
	s.Require().NoError(err)

	count, err := facades.Orm().Query().Model(&models.Player{}).Count()
	s.Require().NoError(err)
	s.Zero(count)
}

func (s *MigrationsTestSuite) TestKickOffIsStoredInUtc() {
	home := models.Club{Name: "Home", Ground: "Home Park"}
	away := models.Club{Name: "Away", Ground: "Away Park"}
	s.Require().NoError(facades.Orm().Query().Create(&home))
	s.Require().NoError(facades.Orm().Query().Create(&away))

	// 3pm British Summer Time (UTC+1) is 2pm UTC
	bst := time.FixedZone("BST", 60*60)
	kickOff := time.Date(2026, 10, 10, 15, 0, 0, 0, bst)
	fixture := models.Fixture{HomeClubID: home.ID, AwayClubID: away.ID, KickOffAt: kickOff, Venue: "Home Park", Status: models.Scheduled}
	s.Require().NoError(facades.Orm().Query().Create(&fixture))

	var stored models.Fixture
	s.Require().NoError(facades.Orm().Query().FindOrFail(&stored, fixture.ID))

	s.True(stored.KickOffAt.Equal(kickOff), "same instant in time")
	s.Equal(time.UTC, stored.KickOffAt.Location())
	s.Equal(14, stored.KickOffAt.Hour())
}
