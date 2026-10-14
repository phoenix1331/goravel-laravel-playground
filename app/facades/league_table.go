package facades

import (
	"goravel/app/services"
)

// LeagueTable is a facade: a plain function that asks the container for the
// binding and returns it with a real type. No static proxy, no __callStatic.
// It follows the same shape as the built-in facades in this folder.
func LeagueTable() services.LeagueTable {
	instance, err := App().Make(services.LeagueTableBinding)
	if err != nil {
		Log().Error(err)
		return nil
	}

	return instance.(services.LeagueTable)
}
