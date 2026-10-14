// Package events holds application events.
package events

import (
	"github.com/goravel/framework/contracts/event"
)

// FixtureResultRecorded fires when a fixture gets a final score.
//
// Laravel:
//
//	class FixtureResultRecorded
//	{
//	    use Dispatchable, SerializesModels;
//	    public function __construct(public Fixture $fixture) {}
//	}
//
//	FixtureResultRecorded::dispatch($fixture);
//
// Goravel: the event carries no properties. Its payload is a list of typed
// arguments (here, just the fixture ID), so it can be serialised for queued
// listeners. Listeners are mapped to events in bootstrap/app.go (WithEvents),
// not discovered from type hints.
type FixtureResultRecorded struct{}

// Handle runs before the listeners and can change the arguments they receive.
func (e *FixtureResultRecorded) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

// Args builds the argument list, so callers cannot get the shape wrong.
func (e *FixtureResultRecorded) Args(fixtureID uint) []event.Arg {
	return []event.Arg{{Type: "uint", Value: fixtureID}}
}
