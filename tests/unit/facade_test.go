// Laravel equivalent: tests/Unit/LeagueTableFacadeTest.php - see docs/chapters/06-facades.md
package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goravel/app/facades"
	"goravel/app/services"
)

// #region fake
// fakeLeagueTable satisfies services.LeagueTable just by having the method.
type fakeLeagueTable struct {
	calls int
}

func (f *fakeLeagueTable) Standings(results []services.Result) []services.Standing {
	f.calls++

	return []services.Standing{{ClubID: 99, Points: 100}}
}

// #endregion fake

// #region test
func TestLeagueTableFacadeResolvesBinding(t *testing.T) {
	standings := facades.LeagueTable().Standings([]services.Result{
		{HomeClubID: 1, AwayClubID: 2, HomeScore: 2, AwayScore: 0},
	})

	require.Len(t, standings, 2)
	assert.Equal(t, uint(1), standings[0].ClubID)
	assert.Equal(t, 3, standings[0].Points)
}

func TestLeagueTableFacadeCanBeSwapped(t *testing.T) {
	original, err := facades.App().Make(services.LeagueTableBinding)
	require.NoError(t, err)

	// Like Facade::swap(): put a different instance in the container. Fresh
	// clears the already-resolved singleton, otherwise Make keeps returning it.
	fake := &fakeLeagueTable{}
	facades.App().Instance(services.LeagueTableBinding, fake)
	facades.App().Fresh(services.LeagueTableBinding)

	// Put the real one back when this test finishes, so other tests in the
	// package (same process, same container) are not affected
	t.Cleanup(func() {
		facades.App().Instance(services.LeagueTableBinding, original)
		facades.App().Fresh(services.LeagueTableBinding)
	})

	standings := facades.LeagueTable().Standings(nil)

	assert.Equal(t, uint(99), standings[0].ClubID)
	assert.Equal(t, 1, fake.calls)
}

// #endregion test
