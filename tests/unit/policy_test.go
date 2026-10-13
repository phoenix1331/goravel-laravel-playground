// Laravel equivalent: tests/Unit/ClubPolicyTest.php - see docs/chapters/20-authorisation.md
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"goravel/app/models"
	"goravel/app/policies"
)

// #region test
// Policies are plain methods, so they can be tested without HTTP or a database.
func TestClubPolicyUpdate(t *testing.T) {
	club := &models.Club{}
	club.ID = 7
	own, other := uint(7), uint(8)

	tests := []struct {
		name    string
		user    *models.User
		allowed bool
	}{
		{"the club's admin", &models.User{ClubID: &own}, true},
		{"another club's admin", &models.User{ClubID: &other}, false},
		{"a user with no club", &models.User{}, false},
	}

	policy := policies.NewClubPolicy()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := policy.Update(context.Background(), map[string]any{"user": tt.user, "club": club})
			assert.Equal(t, tt.allowed, response.Allowed())
		})
	}
}

// #endregion test

func TestClubPolicyCreateOnlyForUsersWithoutAClub(t *testing.T) {
	policy := policies.NewClubPolicy()
	clubID := uint(1)

	assert.True(t, policy.Create(context.Background(), map[string]any{"user": &models.User{}}).Allowed())

	response := policy.Create(context.Background(), map[string]any{"user": &models.User{ClubID: &clubID}})
	assert.False(t, response.Allowed())
	assert.Equal(t, "You already run a club.", response.Message())
}

func TestPolicyDeniesWhenArgumentsAreMissing(t *testing.T) {
	assert.False(t, policies.NewClubPolicy().Update(context.Background(), map[string]any{}).Allowed())
	assert.False(t, policies.NewFixturePolicy().RecordResult(context.Background(), map[string]any{}).Allowed())
}
