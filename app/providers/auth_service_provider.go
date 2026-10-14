package providers

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"

	"goravel/app/facades"
	"goravel/app/policies"
)

type AuthServiceProvider struct{}

func (r *AuthServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings:     []string{},
		Dependencies: []string{binding.Gate},
		ProvideFor:   []string{},
	}
}

func (r *AuthServiceProvider) Register(app foundation.Application) {}

// Boot names each ability and points it at a policy method. There is no
// automatic Club -> ClubPolicy discovery: every ability is defined here.
func (r *AuthServiceProvider) Boot(app foundation.Application) {
	gate := facades.Gate()

	clubs := policies.NewClubPolicy()
	gate.Define("create-club", clubs.Create)
	gate.Define("update-club", clubs.Update)
	gate.Define("delete-club", clubs.Delete)
	gate.Define("manage-players", clubs.ManagePlayers)

	fixtures := policies.NewFixturePolicy()
	gate.Define("record-result", fixtures.RecordResult)
}
