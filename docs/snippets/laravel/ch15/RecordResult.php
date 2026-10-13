<?php

// app/Actions/RecordResult.php
public function __invoke(Fixture $fixture, int $homeScore, int $awayScore): void
{
    $fixture->update([
        'status' => FixtureStatus::Finished,
        'home_score' => $homeScore,
        'away_score' => $awayScore,
    ]);

    FixtureResultRecorded::dispatch($fixture);
}
