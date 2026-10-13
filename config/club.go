// Laravel equivalent: config/club.php - see docs/chapters/02-configuration.md
package config

import (
	"goravel/app/facades"
)

// #region club
// init runs automatically when the config package is imported, before main.
// Each config file registers its own section, so there is no list to maintain.
func init() {
	config := facades.Config()
	config.Add("club", map[string]any{
		// Name shown on the league table
		"league_name": config.Env("LEAGUE_NAME", "Sunday League Division One"),

		// League points awarded per result
		"points": map[string]any{
			"win":  config.Env("LEAGUE_POINTS_WIN", 3),
			"draw": config.Env("LEAGUE_POINTS_DRAW", 1),
			"loss": config.Env("LEAGUE_POINTS_LOSS", 0),
		},

		// Largest squad a club may register
		"max_squad_size": config.Env("CLUB_MAX_SQUAD_SIZE", 25),
	})
}

// #endregion club
