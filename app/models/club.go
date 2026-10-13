// Laravel equivalent: app/Models/Club.php - see docs/chapters/09-models.md
package models

import (
	"github.com/goravel/framework/contracts/database/factory"
	"github.com/goravel/framework/database/orm"

	"goravel/database/factories"
)

// #region model
type Club struct {
	orm.Model
	Name        string  `json:"name"`
	Ground      string  `json:"ground"`
	FoundedYear *int    `json:"founded_year"`
	ShortName   *string `json:"short_name"`

	// #region has-many
	// Has many: a slice of the related model. Empty until loaded with
	// With() or Load(); omitempty leaves it out of JSON when not loaded.
	Players []*Player `json:"players,omitempty"`
	// #endregion has-many
}

// #endregion model

// #region factory
// Factory links the model to its factory, like the HasFactory trait.
func (c *Club) Factory() factory.Factory {
	return &factories.ClubFactory{}
}

// #endregion factory
