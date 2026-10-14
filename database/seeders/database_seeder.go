package seeders

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"goravel/app/facades"
	"goravel/app/models"
)

type DatabaseSeeder struct{}

// Signature The name and signature of the seeder.
func (s *DatabaseSeeder) Signature() string {
	return "DatabaseSeeder"
}

// Run builds a small league: six clubs, a squad and an admin for each, and a
// full round of fixtures, the first half already played.
//
// Laravel:
//
//	public function run(): void
//	{
//	    $clubs = Club::factory()->count(6)->create();
//	    // ...
//	}
//
// Goravel: the same idea, with every step returning an error to check.
// Seeders are listed in bootstrap/seeders.go; make fresh runs this one.
func (s *DatabaseSeeder) Run() error {
	// A fixed seed makes the demo data the same on every run
	if err := gofakeit.Seed(2026); err != nil {
		return err
	}

	var clubs []models.Club
	if err := facades.Orm().Factory().Count(6).Create(&clubs); err != nil {
		return err
	}

	for i, club := range clubs {
		if err := s.squad(club); err != nil {
			return err
		}

		admin := models.User{}
		if err := facades.Orm().Factory().Create(&admin, map[string]any{
			"ClubID": &club.ID,
			"Email":  fmt.Sprintf("admin%d@example.com", i+1),
		}); err != nil {
			return err
		}
	}

	return s.fixtures(clubs)
}

// squad gives a club sixteen players with shirt numbers 1 to 16.
func (s *DatabaseSeeder) squad(club models.Club) error {
	for number := 1; number <= 16; number++ {
		position := models.Midfielder
		switch {
		case number == 1 || number == 13:
			position = models.Goalkeeper
		case number <= 5 || number == 14:
			position = models.Defender
		case number >= 9 && number <= 11:
			position = models.Forward
		}

		var player models.Player
		if err := facades.Orm().Factory().Create(&player, map[string]any{
			"ClubID":      club.ID,
			"SquadNumber": number,
			"Position":    string(position),
		}); err != nil {
			return err
		}
	}

	return nil
}

// fixtures pairs every club with every other club once. The first half are
// finished on past Saturdays; the rest are scheduled on future ones.
func (s *DatabaseSeeder) fixtures(clubs []models.Club) error {
	var pairs [][2]models.Club
	for i := range clubs {
		for j := i + 1; j < len(clubs); j++ {
			pairs = append(pairs, [2]models.Club{clubs[i], clubs[j]})
		}
	}

	lastSaturday := previousSaturday(time.Now().UTC())
	for n, pair := range pairs {
		played := n < len(pairs)/2
		weeks := n - len(pairs)/2 // negative for played fixtures
		kickOff := lastSaturday.AddDate(0, 0, 7*weeks)
		if !played {
			kickOff = kickOff.AddDate(0, 0, 7)
		}

		attributes := map[string]any{
			"HomeClubID": pair[0].ID,
			"AwayClubID": pair[1].ID,
			"KickOffAt":  kickOff,
			"Venue":      pair[0].Ground,
		}
		if played {
			home, away := gofakeit.IntRange(0, 4), gofakeit.IntRange(0, 3)
			attributes["Status"] = string(models.Finished)
			attributes["HomeScore"] = &home
			attributes["AwayScore"] = &away
		}

		var fixture models.Fixture
		if err := facades.Orm().Factory().Create(&fixture, attributes); err != nil {
			return err
		}
	}

	return nil
}

// previousSaturday returns 3pm UTC on the most recent Saturday before t.
func previousSaturday(t time.Time) time.Time {
	daysBack := (int(t.Weekday()) - int(time.Saturday) + 7) % 7
	if daysBack == 0 {
		daysBack = 7
	}
	day := t.AddDate(0, 0, -daysBack)

	return time.Date(day.Year(), day.Month(), day.Day(), 15, 0, 0, 0, time.UTC)
}
