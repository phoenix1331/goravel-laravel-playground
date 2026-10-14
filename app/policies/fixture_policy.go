// Laravel equivalent: app/Policies/FixturePolicy.php
package policies

import (
	"context"

	frameworkaccess "github.com/goravel/framework/auth/access"
	"github.com/goravel/framework/contracts/auth/access"

	"goravel/app/models"
)

type FixturePolicy struct{}

func NewFixturePolicy() *FixturePolicy {
	return &FixturePolicy{}
}

// RecordResult allows only the home club's admin: the home side reports the score.
func (p *FixturePolicy) RecordResult(ctx context.Context, arguments map[string]any) access.Response {
	user, userOK := arguments["user"].(*models.User)
	fixture, fixtureOK := arguments["fixture"].(*models.Fixture)
	if !userOK || !fixtureOK {
		return frameworkaccess.NewDenyResponse("This action is unauthorised.")
	}
	if !isAdminOf(user, fixture.HomeClubID) {
		return frameworkaccess.NewDenyResponse("Only the home club can record the result.")
	}

	return frameworkaccess.NewAllowResponse()
}
