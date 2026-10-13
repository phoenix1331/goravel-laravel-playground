<?php

namespace App\Exceptions;

class AlreadyRecorded extends \RuntimeException {}

// app/Actions/RecordResult.php
if ($fixture->status === FixtureStatus::Finished) {
    throw new AlreadyRecorded("Fixture {$fixture->id} already has a result");
}
