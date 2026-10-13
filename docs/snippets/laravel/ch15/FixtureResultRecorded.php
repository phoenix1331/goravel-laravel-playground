<?php

namespace App\Events;

use App\Models\Fixture;
use Illuminate\Foundation\Events\Dispatchable;
use Illuminate\Queue\SerializesModels;

class FixtureResultRecorded
{
    use Dispatchable, SerializesModels;

    // The event carries the model itself; SerializesModels stores just
    // its ID if a queued listener needs it
    public function __construct(public Fixture $fixture) {}
}
