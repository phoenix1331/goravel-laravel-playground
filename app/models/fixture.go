package models

import (
	"time"

	"github.com/goravel/framework/contracts/database/factory"
	contractsorm "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/database/orm"

	"goravel/app/services"
	"goravel/database/factories"
)

type FixtureStatus string

const (
	Scheduled  FixtureStatus = "scheduled"
	InProgress FixtureStatus = "in_progress"
	Finished   FixtureStatus = "finished"
	Postponed  FixtureStatus = "postponed"
)

// Fixture is a match between two clubs.
//
// Laravel:
//
//	protected function casts(): array
//	{
//	    return ['kick_off_at' => 'datetime', 'status' => FixtureStatus::class];
//	}
//	public function homeClub(): BelongsTo { return $this->belongsTo(Club::class, 'home_club_id'); }
//	public function awayClub(): BelongsTo { return $this->belongsTo(Club::class, 'away_club_id'); }
//
// Goravel: time.Time is the datetime cast. Nullable scores are *int, because a
// plain int cannot be NULL and would read back as 0, a real score. Two
// relations to one table name their keys in a gorm tag.
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

	// Two relations to the same model. The names do not follow the Club + ID
	// convention, so the gorm tag names each foreign key explicitly, like
	// belongsTo(Club::class, 'home_club_id') in Laravel.
	HomeClub *Club `gorm:"foreignKey:HomeClubID" json:"home_club,omitempty"`
	AwayClub *Club `gorm:"foreignKey:AwayClubID" json:"away_club,omitempty"`
}

// IsFinished is an ordinary method.
//
// Laravel: an accessor, read as a property: $fixture->is_finished
// Goravel: just call it: fixture.IsFinished()
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

// DispatchesEvents hooks into the model's life cycle. Saving runs before
// every create and update.
//
// Laravel:
//
//	protected function kickOffAt(): Attribute
//	{
//	    return Attribute::set(fn ($value) => Carbon::parse($value)->utc());
//	}
//
// It converts kick-off to UTC. MySQL does this on its own, but SQLite stores
// the time as text with its original offset, so without this a 3pm BST
// kick-off reads back as 15:00+01:00 on SQLite and 14:00 UTC on MySQL.
func (f *Fixture) DispatchesEvents() map[contractsorm.EventType]func(contractsorm.Event) error {
	return map[contractsorm.EventType]func(contractsorm.Event) error{
		contractsorm.EventSaving: func(event contractsorm.Event) error {
			if kickOff, ok := event.GetAttribute("kick_off_at").(time.Time); ok {
				event.SetAttribute("kick_off_at", kickOff.UTC())
			}

			return nil
		},
	}
}

// Factory links the model to its factory, like the HasFactory trait.
func (f *Fixture) Factory() factory.Factory {
	return &factories.FixtureFactory{}
}
