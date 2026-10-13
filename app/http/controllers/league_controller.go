// Laravel equivalent: app/Http/Controllers/LeagueController.php - see docs/chapters/13-query-builder.md
package controllers

import (
	"github.com/goravel/framework/contracts/http"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/services"
)

type LeagueController struct{}

func NewLeagueController() *LeagueController {
	return &LeagueController{}
}

// leagueRow is one line of the table in the API response. Embedding
// services.Standing flattens its fields into the JSON object.
type leagueRow struct {
	Position int    `json:"position"`
	Club     string `json:"club"`
	services.Standing
	GoalDifference int `json:"goal_difference"`
}

// #region table
// Table builds the league table from finished fixtures.
func (r *LeagueController) Table(ctx http.Context) http.Response {
	// facades.DB() is the query builder, like DB::table(). Rows scan into
	// any struct: column home_club_id fills field HomeClubID by name.
	var results []services.Result
	err := facades.DB().Table("fixtures").
		Select("home_club_id", "away_club_id", "home_score", "away_score").
		Where("status", models.Finished).
		Get(&results)
	if err != nil {
		return serverError(ctx, err)
	}

	standings := facades.LeagueTable().Standings(results)

	// Clubs that have not played yet still appear, at the bottom
	var clubs []models.Club
	if err := facades.Orm().Query().Order("name").Get(&clubs); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": tableRows(standings, clubs)})
}

// #endregion table

// tableRows adds positions and club names, and appends clubs with no results.
func tableRows(standings []services.Standing, clubs []models.Club) []leagueRow {
	names := make(map[uint]string, len(clubs))
	for _, club := range clubs {
		names[club.ID] = club.Name
	}

	rows := make([]leagueRow, 0, len(clubs))
	played := make(map[uint]bool, len(standings))
	for _, standing := range standings {
		played[standing.ClubID] = true
		rows = append(rows, leagueRow{Club: names[standing.ClubID], Standing: standing, GoalDifference: standing.GoalDifference()})
	}
	for _, club := range clubs {
		if !played[club.ID] {
			rows = append(rows, leagueRow{Club: club.Name, Standing: services.Standing{ClubID: club.ID}})
		}
	}
	for i := range rows {
		rows[i].Position = i + 1
	}

	return rows
}

// #region goals
// goalsRow receives one row of the raw query. Raw SQL needs explicit db tags
// only when a column name does not match the field name in snake_case.
type goalsRow struct {
	ClubID uint   `json:"club_id"`
	Club   string `json:"club"`
	Goals  int    `json:"goals"`
}

// goalsSQL totals goals per club. A goal can be scored at home or away, so
// the two columns are stacked with UNION ALL and then summed. It sticks to
// SQL that MySQL and SQLite both accept.
const goalsSQL = `
SELECT clubs.id AS club_id, clubs.name AS club, SUM(scored.goals) AS goals
FROM (
	SELECT home_club_id AS club_id, home_score AS goals FROM fixtures WHERE status = ?
	UNION ALL
	SELECT away_club_id AS club_id, away_score AS goals FROM fixtures WHERE status = ?
) AS scored
JOIN clubs ON clubs.id = scored.club_id
GROUP BY clubs.id, clubs.name
ORDER BY goals DESC, clubs.name`

// Goals lists clubs by goals scored, using raw SQL like DB::select().
func (r *LeagueController) Goals(ctx http.Context) http.Response {
	var rows []goalsRow
	if err := facades.DB().Select(&rows, goalsSQL, models.Finished, models.Finished); err != nil {
		return serverError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{"data": rows})
}

// #endregion goals
