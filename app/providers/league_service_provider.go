package providers

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"

	"goravel/app/services"
)

type LeagueServiceProvider struct{}

// Relationship tells Goravel which bindings this provider adds and which it
// needs, so providers are registered in dependency order.
func (r *LeagueServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings:     []string{services.LeagueTableBinding},
		Dependencies: []string{binding.Config},
		ProvideFor:   []string{},
	}
}

// Register binds services into the container, like register() in Laravel.
// The callback runs lazily, the first time something asks for the binding.
func (r *LeagueServiceProvider) Register(app foundation.Application) {
	app.Singleton(services.LeagueTableBinding, func(app foundation.Application) (any, error) {
		config := app.MakeConfig()

		return services.NewPointsLeagueTable(
			config.GetInt("club.points.win"),
			config.GetInt("club.points.draw"),
			config.GetInt("club.points.loss"),
		), nil
	})
}

// Boot runs after every provider is registered, like boot() in Laravel.
// Nothing to do here yet.
func (r *LeagueServiceProvider) Boot(app foundation.Application) {}
