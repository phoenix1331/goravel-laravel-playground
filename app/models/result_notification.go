// Laravel equivalent: app/Models/ResultNotification.php
package models

import (
	"github.com/goravel/framework/database/orm"
)

// ResultNotification tells a club admin that one of their fixtures has a result.
type ResultNotification struct {
	orm.Model
	UserID    uint   `json:"user_id"`
	FixtureID uint   `json:"fixture_id"`
	Message   string `json:"message"`
}
