// Laravel equivalent: database/factories/UserFactory.php - see docs/chapters/12-factories-seeders.md
package factories

import (
	"github.com/brianvoe/gofakeit/v7"
)

// UserFactory sets Password to an already-hashed value, so creating users in
// tests does not pay for a bcrypt hash each time. Chapter 19 logs in with it.
// The plain-text password is "password", as in Laravel's UserFactory.
type UserFactory struct{}

// PasswordHash is bcrypt("password").
const PasswordHash = "$2a$12$wenAdbQLsf4CQmRLUTQhuuCcM6Hm8T2bQWlWMmARMC.Wd6iBd3526"

func (f *UserFactory) Definition() map[string]any {
	return map[string]any{
		"Name":     gofakeit.Name(),
		"Email":    gofakeit.Email(),
		"Password": PasswordHash,
	}
}
