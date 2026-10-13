// Laravel equivalent: routes/console.php - see docs/chapters/17-scheduling.md
package bootstrap

import (
	"time"

	"github.com/goravel/framework/contracts/schedule"

	"goravel/app/facades"
	"goravel/app/models"
)

// #region schedule
// Schedule lists the recurring tasks. The scheduler runner started by
// app.Start() checks them every second; there is no cron entry to add.
func Schedule() []schedule.Event {
	return []schedule.Event{
		// Pull results in every fifteen minutes on match day. If one run is
		// still going when the next is due, skip it rather than overlap.
		facades.Schedule().Command("sync:fixtures").
			EveryFifteenMinutes().
			Saturdays().
			SkipIfStillRunning(),

		// Closures work too, like Schedule::call()
		facades.Schedule().Call(pruneOldNotifications).
			DailyAt("03:00").
			Name("prune-old-notifications"),
	}
}

// #endregion schedule

// pruneOldNotifications deletes result notifications older than 30 days.
func pruneOldNotifications() {
	cutoff := time.Now().UTC().AddDate(0, 0, -30)
	if _, err := facades.Orm().Query().Where("created_at < ?", cutoff).Delete(&models.ResultNotification{}); err != nil {
		facades.Log().Error(err)
	}
}
