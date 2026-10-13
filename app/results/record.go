// Laravel equivalent: app/Actions/RecordResult.php - see docs/chapters/14-console-commands.md
package results

import (
	"fmt"

	"goravel/app/events"
	"goravel/app/facades"
	"goravel/app/models"
)

// #region record
// Record marks a fixture as finished with the given score. Every way a result
// arrives (the sync command, the API) goes through here, so later chapters
// have one place to hook in.
func Record(fixture *models.Fixture, homeScore, awayScore int) error {
	if homeScore < 0 || awayScore < 0 {
		return fmt.Errorf("record fixture %d: scores cannot be negative", fixture.ID)
	}

	fixture.Status = models.Finished
	fixture.HomeScore = &homeScore
	fixture.AwayScore = &awayScore

	if err := facades.Orm().Query().Save(fixture); err != nil {
		return fmt.Errorf("record fixture %d: %w", fixture.ID, err)
	}

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
