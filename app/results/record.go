// Laravel equivalent: app/Actions/RecordResult.php - see docs/chapters/14-console-commands.md
package results

import (
	"errors"
	"fmt"

	"goravel/app/events"
	"goravel/app/facades"
	"goravel/app/models"
)

// #region sentinel
// ErrAlreadyRecorded is a sentinel error: a package-level value that callers
// compare against with errors.Is, the way PHP code catches a specific
// exception class. Callers decide what it means to them (a 409, a skip).
var ErrAlreadyRecorded = errors.New("fixture already has a result")

// #endregion sentinel

// #region record
// Record marks a fixture as finished with the given score. Every way a result
// arrives (the sync command, the API) goes through here, so later chapters
// have one place to hook in.
func Record(fixture *models.Fixture, homeScore, awayScore int) error {
	if homeScore < 0 || awayScore < 0 {
		return fmt.Errorf("record fixture %d: scores cannot be negative", fixture.ID)
	}
	// #region wrap
	if fixture.Status == models.Finished {
		// %w wraps the sentinel: the message gains context ("record fixture
		// 7: fixture already has a result") and errors.Is still matches it
		return fmt.Errorf("record fixture %d: %w", fixture.ID, ErrAlreadyRecorded)
	}
	// #endregion wrap

	fixture.Status = models.Finished
	fixture.HomeScore = &homeScore
	fixture.AwayScore = &awayScore

	if err := facades.Orm().Query().Save(fixture); err != nil {
		return fmt.Errorf("record fixture %d: %w", fixture.ID, err)
	}

	// #region log
	// Structured context instead of a formatted string: each key becomes a
	// JSON field in the results channel, like Log::channel('results')
	// ->info('Result recorded', [...])
	facades.Log().Channel("results").With(map[string]any{
		"fixture_id": fixture.ID,
		"home_score": homeScore,
		"away_score": awayScore,
	}).Info("result recorded")
	// #endregion log

	// #region dispatch
	// Tell the rest of the app. Like event(new FixtureResultRecorded($fixture)),
	// but the payload is a list of typed arguments.
	recorded := &events.FixtureResultRecorded{}
	if err := facades.Event().Job(recorded, recorded.Args(fixture.ID)).Dispatch(); err != nil {
		return fmt.Errorf("record fixture %d: %w", fixture.ID, err)
	}
	// #endregion dispatch

	return nil
}

// #endregion record
