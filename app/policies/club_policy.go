// Package policies holds authorisation rules, registered as gates in
// app/providers/auth_service_provider.go.
package policies

import (
	"context"

	frameworkaccess "github.com/goravel/framework/auth/access"
	"github.com/goravel/framework/contracts/auth/access"

	"goravel/app/models"
)

// ClubPolicy decides what a user may do to a club. Each method has the gate
// callback signature, so it can be passed straight to facades.Gate().Define.
//
// Laravel:
//
//	public function update(User $user, Club $club): Response
//	{
//	    return $user->club_id === $club->id
//	        ? Response::allow()
//	        : Response::deny('You can only manage your own club.');
//	}
//
// Goravel: the user is not passed in for you; it arrives in the arguments
// map, along with the club, and needs a type assertion. No policy discovery:
// each method is registered as a named ability in AuthServiceProvider.
type ClubPolicy struct{}

func NewClubPolicy() *ClubPolicy {
	return &ClubPolicy{}
}

// Create allows users who do not run a club yet; they become its admin.
func (p *ClubPolicy) Create(ctx context.Context, arguments map[string]any) access.Response {
	user, ok := arguments["user"].(*models.User)
	if !ok {
		return frameworkaccess.NewDenyResponse("You must be logged in.")
	}
	if user.ClubID != nil {
		return frameworkaccess.NewDenyResponse("You already run a club.")
	}

	return frameworkaccess.NewAllowResponse()
}

// Update allows only the club's own admin. Delete and ManagePlayers use it too.
func (p *ClubPolicy) Update(ctx context.Context, arguments map[string]any) access.Response {
	user, userOK := arguments["user"].(*models.User)
	club, clubOK := arguments["club"].(*models.Club)
	if !userOK || !clubOK {
		return frameworkaccess.NewDenyResponse("This action is unauthorised.")
	}
	if !isAdminOf(user, club.ID) {
		return frameworkaccess.NewDenyResponse("You can only manage your own club.")
	}

	return frameworkaccess.NewAllowResponse()
}

func (p *ClubPolicy) Delete(ctx context.Context, arguments map[string]any) access.Response {
	return p.Update(ctx, arguments)
}

func (p *ClubPolicy) ManagePlayers(ctx context.Context, arguments map[string]any) access.Response {
	return p.Update(ctx, arguments)
}

// isAdminOf reports whether user administers the club with the given ID.
func isAdminOf(user *models.User, clubID uint) bool {
	return user.ClubID != nil && *user.ClubID == clubID
}
