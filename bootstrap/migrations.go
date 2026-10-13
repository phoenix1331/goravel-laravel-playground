package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20261013182221CreateClubsTable{},
		&migrations.M20261013183449CreatePlayersTable{},
		&migrations.M20261013183455CreateFixturesTable{},
		&migrations.M20261013183501CreateUsersTable{},
		&migrations.M20261013184001AddShortNameToClubsTable{},
		&migrations.M20261013184019AddKickOffIndexToFixturesTable{},
	}
}
