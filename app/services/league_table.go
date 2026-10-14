package services

import (
	"sort"
)

// LeagueTableBinding is the container key the service is bound under.
const LeagueTableBinding = "club.league_table"

// LeagueTableCacheKey is where the rendered table is cached.
const LeagueTableCacheKey = "league_table"

// Result is one finished fixture.
type Result struct {
	HomeClubID uint
	AwayClubID uint
	HomeScore  int
	AwayScore  int
}

// Standing is one club's row in the league table.
type Standing struct {
	ClubID       uint `json:"club_id"`
	Played       int  `json:"played"`
	Won          int  `json:"won"`
	Drawn        int  `json:"drawn"`
	Lost         int  `json:"lost"`
	GoalsFor     int  `json:"goals_for"`
	GoalsAgainst int  `json:"goals_against"`
	Points       int  `json:"points"`
}

// GoalDifference is a method on the struct, like an accessor on a PHP class.
func (s Standing) GoalDifference() int {
	return s.GoalsFor - s.GoalsAgainst
}

// LeagueTable is what callers depend on. Any type with a matching Standings
// method satisfies it: there is no "implements" keyword in Go.
type LeagueTable interface {
	Standings(results []Result) []Standing
}

// PointsLeagueTable ranks clubs by points, then goal difference, then goals scored.
type PointsLeagueTable struct {
	win, draw, loss int
}

func NewPointsLeagueTable(win, draw, loss int) *PointsLeagueTable {
	return &PointsLeagueTable{win: win, draw: draw, loss: loss}
}

func (t *PointsLeagueTable) Standings(results []Result) []Standing {
	rows := map[uint]*Standing{}
	row := func(clubID uint) *Standing {
		if rows[clubID] == nil {
			rows[clubID] = &Standing{ClubID: clubID}
		}

		return rows[clubID]
	}

	for _, result := range results {
		home, away := row(result.HomeClubID), row(result.AwayClubID)
		t.record(home, result.HomeScore, result.AwayScore)
		t.record(away, result.AwayScore, result.HomeScore)
	}

	table := make([]Standing, 0, len(rows))
	for _, standing := range rows {
		table = append(table, *standing)
	}

	sort.Slice(table, func(i, j int) bool {
		a, b := table[i], table[j]
		if a.Points != b.Points {
			return a.Points > b.Points
		}
		if a.GoalDifference() != b.GoalDifference() {
			return a.GoalDifference() > b.GoalDifference()
		}
		if a.GoalsFor != b.GoalsFor {
			return a.GoalsFor > b.GoalsFor
		}

		return a.ClubID < b.ClubID
	})

	return table
}

func (t *PointsLeagueTable) record(standing *Standing, scored, conceded int) {
	standing.Played++
	standing.GoalsFor += scored
	standing.GoalsAgainst += conceded

	switch {
	case scored > conceded:
		standing.Won++
		standing.Points += t.win
	case scored == conceded:
		standing.Drawn++
		standing.Points += t.draw
	default:
		standing.Lost++
		standing.Points += t.loss
	}
}
