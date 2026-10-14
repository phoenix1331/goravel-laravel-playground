package factories

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// FixtureFactory makes a scheduled fixture on a Saturday at 3pm UTC within
// the next eight weeks. Callers set HomeClubID and AwayClubID.
type FixtureFactory struct{}

func (f *FixtureFactory) Definition() map[string]any {
	now := time.Now().UTC()
	daysToSaturday := (int(time.Saturday) - int(now.Weekday()) + 7) % 7
	saturday := now.AddDate(0, 0, daysToSaturday+7*gofakeit.IntRange(0, 7))

	return map[string]any{
		"KickOffAt": time.Date(saturday.Year(), saturday.Month(), saturday.Day(), 15, 0, 0, 0, time.UTC),
		"Venue":     gofakeit.RandomString(towns) + " Park",
		"Status":    "scheduled",
	}
}
