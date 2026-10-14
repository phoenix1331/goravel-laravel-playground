package bootstrap

import (
	"time"

	"github.com/goravel/framework/contracts/schedule"

	"goravel/app/facades"
	"goravel/app/models"
)

// Schedule lists the recurring tasks.
//
// Laravel (routes/console.php):
//
//	Schedule::command('sync:fixtures')->everyFifteenMinutes()->saturdays()->withoutOverlapping();
//	Schedule::call(fn () => ...)->dailyAt('03:00')->name('prune-old-notifications');
//
// Goravel: the same fluent API. The scheduler runs inside the app process,
// started by app.Start(), so there is no cron entry. Every running copy of
// the app runs it, so add OnOneServer() before running several.
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

// pruneOldNotifications deletes result notifications older than 30 days.
func pruneOldNotifications() {
	cutoff := time.Now().UTC().AddDate(0, 0, -30)
	if _, err := facades.Orm().Query().Where("created_at < ?", cutoff).Delete(&models.ResultNotification{}); err != nil {
		facades.Log().Error(err)
	}
}
