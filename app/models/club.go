// Laravel equivalent: app/Models/Club.php - see docs/chapters/09-models.md
package models

import (
	"github.com/goravel/framework/database/orm"
)

// #region model
type Club struct {
	orm.Model
	Name        string `json:"name"`
	Ground      string `json:"ground"`
	FoundedYear *int   `json:"founded_year"`
}

// #endregion model
