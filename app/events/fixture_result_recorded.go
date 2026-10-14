// Laravel equivalent: app/Events/FixtureResultRecorded.php
package events

import (
	"github.com/goravel/framework/contracts/event"
)

// FixtureResultRecorded fires when a fixture gets a final score.
//
// The event carries its data as a list of typed arguments, not as struct
// fields: listeners may run on a queue, and the arguments are what gets
// serialised. Here the only argument is the fixture ID.
type FixtureResultRecorded struct{}

// Handle runs before the listeners and can change the arguments they receive.
func (e *FixtureResultRecorded) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

// Args builds the argument list, so callers cannot get the shape wrong.
func (e *FixtureResultRecorded) Args(fixtureID uint) []event.Arg {
	return []event.Arg{{Type: "uint", Value: fixtureID}}
}
