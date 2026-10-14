// Package models holds the ORM models.
package models

import (
	"github.com/goravel/framework/contracts/database/factory"
	"github.com/goravel/framework/database/orm"

	"goravel/database/factories"
)

// Club is a football club.
//
// Laravel:
//
//	class Club extends Model
//	{
//	    use HasFactory;
//	    protected $fillable = ['name', 'ground', 'founded_year', 'short_name'];
//	    protected function casts(): array { return ['founded_year' => 'integer']; }
//	    public function players(): HasMany { return $this->hasMany(Player::class); }
//	}
//
// Goravel: a struct with every column declared as a typed field. There is no
// __get, so a misspelt field fails to compile. No $fillable: you build the
// struct yourself. No $casts: the field type is the cast. orm.Model is
// embedded (its ID and timestamps become fields of Club), not inherited.
type Club struct {
	orm.Model
	Name        string  `json:"name"`
	Ground      string  `json:"ground"`
	FoundedYear *int    `json:"founded_year"`
	ShortName   *string `json:"short_name"`

	// Has many, as a field rather than a players() method. Empty until loaded
	// with With("Players") or Load(); there is no lazy loading.
	// omitempty leaves it out of JSON when it was not loaded.
	Players []*Player `json:"players,omitempty"`
}

// Factory links the model to its factory, the equivalent of use HasFactory.
func (c *Club) Factory() factory.Factory {
	return &factories.ClubFactory{}
}
