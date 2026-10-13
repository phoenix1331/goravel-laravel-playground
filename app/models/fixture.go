// Laravel equivalent: app/Models/Fixture.php - see docs/chapters/09-models.md
package models

import (
	"time"

	"github.com/goravel/framework/database/orm"

	"goravel/app/services"
)

type FixtureStatus string

const (
	Scheduled  FixtureStatus = "scheduled"
	InProgress FixtureStatus = "in_progress"
	Finished   FixtureStatus = "finished"
	Postponed  FixtureStatus = "postponed"
)

// #region fixture
type Fixture struct {
	orm.Model
	HomeClubID uint          `json:"home_club_id"`
	AwayClubID uint          `json:"away_club_id"`
	KickOffAt  time.Time     `json:"kick_off_at"`
	Venue      string        `json:"venue"`
	Status     FixtureStatus `json:"status"`
	// Pointers, because "no score yet" (NULL) is different from 0
	HomeScore *int `json:"home_score"`
	AwayScore *int `json:"away_score"`
}

// #endregion fixture

// #region methods
// IsFinished is an ordinary method. Where Laravel might use an accessor
// ($fixture->is_finished), Go calls a method: fixture.IsFinished().
func (f *Fixture) IsFinished() bool {
	return f.Status == Finished && f.HomeScore != nil && f.AwayScore != nil
}

// Result converts a finished fixture for the league table. The second return
// value is false when there is no result yet: Go's "comma ok" idiom.
func (f *Fixture) Result() (services.Result, bool) {
	if !f.IsFinished() {
		return services.Result{}, false
	}

	return services.Result{
		HomeClubID: f.HomeClubID,
		AwayClubID: f.AwayClubID,
		HomeScore:  *f.HomeScore,
		AwayScore:  *f.AwayScore,
	}, true
}

// #endregion methods
