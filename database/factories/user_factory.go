// Laravel equivalent: database/factories/UserFactory.php - see docs/chapters/12-factories-seeders.md
package factories

import (
	"sync"

	"github.com/brianvoe/gofakeit/v7"

	"goravel/app/facades"
)

// The plain-text password for every factory user, as in Laravel's UserFactory.
const Password = "password"

var (
	passwordHash string
	hashOnce     sync.Once
)

// hashedPassword hashes Password once per process and reuses it, like the
// static $password in Laravel's UserFactory. The cost comes from
// BCRYPT_ROUNDS, which .env.testing lowers to 4 so test logins are quick.
func hashedPassword() string {
	hashOnce.Do(func() {
		hash, err := facades.Hash().Make(Password)
		if err != nil {
			panic(err)
		}
		passwordHash = hash
	})

	return passwordHash
}

type UserFactory struct{}

func (f *UserFactory) Definition() map[string]any {
	return map[string]any{
		"Name":     gofakeit.Name(),
		"Email":    gofakeit.Email(),
		"Password": hashedPassword(),
	}
}
