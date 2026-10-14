// Laravel equivalent: tests/Feature/FactoriesTest.php
package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/database/seeders"
	"goravel/tests"
)

type FactoriesTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestFactoriesTestSuite(t *testing.T) {
	suite.Run(t, new(FactoriesTestSuite))
}

func (s *FactoriesTestSuite) SetupTest() {
	s.Require().NoError(s.TruncateTables())
}

func (s *FactoriesTestSuite) TestCountCreatesManyRows() {
	// Club::factory()->count(3)->create()
	var clubs []models.Club
	s.Require().NoError(facades.Orm().Factory().Count(3).Create(&clubs))

	s.Len(clubs, 3)
	s.NotZero(clubs[0].ID)
}

func (s *FactoriesTestSuite) TestAttributesOverrideTheDefinition() {
	// Player::factory()->for($club)->create(['squad_number' => 7])
	var club models.Club
	s.Require().NoError(facades.Orm().Factory().Create(&club))

	var player models.Player
	s.Require().NoError(facades.Orm().Factory().Create(&player, map[string]any{
		"ClubID":      club.ID,
		"SquadNumber": 7,
	}))

	s.Equal(club.ID, player.ClubID)
	s.Equal(7, player.SquadNumber)
	s.Contains(models.Positions, player.Position)
}

func (s *FactoriesTestSuite) TestMakeDoesNotSave() {
	// Club::factory()->make()
	var club models.Club
	s.Require().NoError(facades.Orm().Factory().Make(&club))

	s.NotEmpty(club.Name)
	s.Zero(club.ID)
}

func (s *FactoriesTestSuite) TestDatabaseSeederBuildsALeague() {
	// $this->seed(DatabaseSeeder::class)
	s.Seed(&seeders.DatabaseSeeder{})

	s.Equal(int64(6), s.count(&models.Club{}))
	s.Equal(int64(96), s.count(&models.Player{}))
	s.Equal(int64(15), s.count(&models.Fixture{}))
	s.Equal(int64(6), s.count(&models.User{}))
}

func (s *FactoriesTestSuite) count(model any) int64 {
	count, err := facades.Orm().Query().Model(model).Count()
	s.Require().NoError(err)

	return count
}
