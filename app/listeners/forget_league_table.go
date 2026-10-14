package listeners

import (
	"github.com/goravel/framework/contracts/event"

	"goravel/app/facades"
	"goravel/app/services"
)

// ForgetLeagueTable clears the cached league table when a result comes in,
// so the next request rebuilds it. It runs synchronously, so the very next
// request after recording a result already sees the new table.
//
// Laravel: Cache::forget('league_table');
type ForgetLeagueTable struct{}

func (l *ForgetLeagueTable) Signature() string {
	return "forget_league_table"
}

func (l *ForgetLeagueTable) Queue(args ...any) event.Queue {
	return event.Queue{Enable: false}
}

func (l *ForgetLeagueTable) Handle(args ...any) error {
	facades.Cache().Forget(services.LeagueTableCacheKey)

	return nil
}
