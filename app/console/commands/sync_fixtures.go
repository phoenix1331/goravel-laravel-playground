// Laravel equivalent: app/Console/Commands/SyncFixtures.php
package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/results"
	"goravel/app/services"
)

type SyncFixtures struct{}

// Signature is the command name: ./artisan sync:fixtures
func (r *SyncFixtures) Signature() string {
	return "sync:fixtures"
}

func (r *SyncFixtures) Description() string {
	return "Record results from the league's fixture feed"
}

// Extend declares options, like {--file=} {--dry-run} in a Laravel signature.
func (r *SyncFixtures) Extend() command.Extend {
	return command.Extend{
		Category: "sync",
		Flags: []command.Flag{
			&command.StringFlag{Name: "file", Usage: "Read results from this JSON file instead of the configured feed"},
			&command.BoolFlag{Name: "dry-run", Usage: "Show what would change without saving"},
		},
	}
}

// Handle runs the command. Returning an error makes the command exit non-zero.
func (r *SyncFixtures) Handle(ctx console.Context) error {
	feedResults, err := newFeed(ctx.Option("file")).Results()
	if err != nil {
		ctx.Error(err.Error())
		return err
	}

	dryRun := ctx.OptionBool("dry-run")
	var rows [][]string
	for _, result := range feedResults {
		fixture, err := findOpenFixture(result.Home, result.Away)
		if err != nil {
			rows = append(rows, []string{result.Home, result.Away, "-", "skipped: " + err.Error()})
			continue
		}

		score := fmt.Sprintf("%d-%d", result.HomeScore, result.AwayScore)
		if dryRun {
			rows = append(rows, []string{result.Home, result.Away, score, "would record"})
			continue
		}
		if err := results.Record(fixture, result.HomeScore, result.AwayScore); err != nil {
			return err
		}
		rows = append(rows, []string{result.Home, result.Away, score, "recorded"})
	}

	ctx.Table([]string{"Home", "Away", "Score", "Outcome"}, rows)

	return nil
}

// newFeed picks where results come from: --file wins, then the configured
// driver. Both feeds satisfy services.FixtureFeed, so Handle does not care.
func newFeed(file string) services.FixtureFeed {
	if file != "" {
		return services.FileFeed{Path: file}
	}
	if facades.Config().GetString("club.feed.driver") == "http" {
		return services.HttpFeed{Request: facades.Http().Client("fixtures")}
	}

	return services.FileFeed{Path: facades.Config().GetString("club.feed.path")}
}

// findOpenFixture finds the unfinished fixture between two clubs, by name.
// Each pair of clubs meets once per season with the same home side, so home
// and away club identify the fixture.
func findOpenFixture(home, away string) (*models.Fixture, error) {
	var homeClub, awayClub models.Club
	if err := facades.Orm().Query().Where("name", home).First(&homeClub); err != nil || homeClub.ID == 0 {
		return nil, fmt.Errorf("unknown club %q", home)
	}
	if err := facades.Orm().Query().Where("name", away).First(&awayClub); err != nil || awayClub.ID == 0 {
		return nil, fmt.Errorf("unknown club %q", away)
	}

	var fixture models.Fixture
	err := facades.Orm().Query().
		Where("home_club_id", homeClub.ID).
		Where("away_club_id", awayClub.ID).
		Where("status <> ?", models.Finished).
		First(&fixture)
	if err != nil || fixture.ID == 0 {
		return nil, fmt.Errorf("no open fixture")
	}

	return &fixture, nil
}
