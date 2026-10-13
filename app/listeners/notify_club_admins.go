// Laravel equivalent: app/Listeners/NotifyClubAdmins.php - see docs/chapters/15-events.md
package listeners

import (
	"fmt"

	"github.com/goravel/framework/contracts/event"
	"github.com/spf13/cast"

	"goravel/app/facades"
	"goravel/app/jobs"
)

// #region listener
// NotifyClubAdmins reacts to a new result by queueing the notifications job.
type NotifyClubAdmins struct{}

func (l *NotifyClubAdmins) Signature() string {
	return "notify_club_admins"
}

// Queue says whether to run the listener itself on a queue. It does not need
// to: it only dispatches a job, which is quick, and the job is what queues.
func (l *NotifyClubAdmins) Queue(args ...any) event.Queue {
	return event.Queue{Enable: false}
}

// Handle receives the event's arguments as ...any, in the order they were
// dispatched. cast.ToUint converts whatever arrived into a uint.
func (l *NotifyClubAdmins) Handle(args ...any) error {
	if len(args) == 0 {
		return fmt.Errorf("notify club admins: missing fixture id")
	}
	fixtureID := cast.ToUint(args[0])

	// #region dispatch
	// Like SendResultNotifications::dispatch($fixture). The job goes on the
	// connection named by QUEUE_CONNECTION: sync in tests, redis in dev.
	job := &jobs.SendResultNotifications{}
	if err := facades.Queue().Job(job, job.Args(fixtureID)).Dispatch(); err != nil {
		return fmt.Errorf("notify club admins: %w", err)
	}
	// #endregion dispatch

	return nil
}

// #endregion listener
