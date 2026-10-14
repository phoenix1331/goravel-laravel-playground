package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goravel/app/facades"
	"goravel/app/services"
)

// A table-driven test: each case is a row of data, and one loop runs them all.
// It plays the same role as a PHPUnit data provider.
func TestPointsLeagueTableRanking(t *testing.T) {
	tests := []struct {
		name    string
		results []services.Result
		want    []uint // club IDs, top of the table first
	}{
		{
			name:    "a win beats a loss",
			results: []services.Result{{HomeClubID: 1, AwayClubID: 2, HomeScore: 0, AwayScore: 2}},
			want:    []uint{2, 1},
		},
		{
			name: "level on points, goal difference decides",
			results: []services.Result{
				{HomeClubID: 1, AwayClubID: 3, HomeScore: 1, AwayScore: 0},
				{HomeClubID: 2, AwayClubID: 3, HomeScore: 4, AwayScore: 0},
			},
			want: []uint{2, 1, 3},
		},
		{
			name:    "a draw is level on everything, so lower club ID first",
			results: []services.Result{{HomeClubID: 2, AwayClubID: 1, HomeScore: 1, AwayScore: 1}},
			want:    []uint{1, 2},
		},
	}

	table := services.NewPointsLeagueTable(3, 1, 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []uint
			for _, standing := range table.Standings(tt.results) {
				got = append(got, standing.ClubID)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPointsLeagueTableCountsResults(t *testing.T) {
	table := services.NewPointsLeagueTable(3, 1, 0)
	standings := table.Standings([]services.Result{
		{HomeClubID: 1, AwayClubID: 2, HomeScore: 3, AwayScore: 1},
		{HomeClubID: 2, AwayClubID: 1, HomeScore: 2, AwayScore: 2},
	})

	require.Len(t, standings, 2)
	assert.Equal(t, services.Standing{ClubID: 1, Played: 2, Won: 1, Drawn: 1, GoalsFor: 5, GoalsAgainst: 3, Points: 4}, standings[0])
	assert.Equal(t, 2, standings[0].GoalDifference())
}

func TestContainerResolvesLeagueTableSingleton(t *testing.T) {
	first, err := facades.App().Make(services.LeagueTableBinding)
	require.NoError(t, err)

	// Make returns `any`. A type assertion checks the value really is a
	// LeagueTable; ok is false (instead of a panic) if it is not.
	table, ok := first.(services.LeagueTable)
	require.True(t, ok)
	assert.NotEmpty(t, table.Standings([]services.Result{{HomeClubID: 1, AwayClubID: 2}}))

	// Singleton: the second Make returns the very same instance
	second, err := facades.App().Make(services.LeagueTableBinding)
	require.NoError(t, err)
	assert.Same(t, first, second)
}
