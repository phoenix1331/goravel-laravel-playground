// Laravel equivalent: app/Models/User.php - see docs/chapters/09-models.md
package models

import (
	"github.com/goravel/framework/database/orm"
)

// #region user
type User struct {
	orm.Model
	// The club this user administers, or nil
	ClubID *uint  `json:"club_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	// json:"-" leaves the field out of every JSON response, like $hidden
	Password string `json:"-"`
}

// #endregion user
